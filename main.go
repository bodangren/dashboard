package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"dashboard/internal/agents"
	"dashboard/internal/api"
	gitpkg "dashboard/internal/git"
	"dashboard/internal/scheduler"
	"dashboard/internal/ws"
)

//go:embed static
var staticFiles embed.FS

func main() {
	// Discover repos under ~/Desktop
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("could not determine home dir: %v", err)
	}
	desktop := filepath.Join(home, "Desktop")

	repos, err := gitpkg.ScanRepos(desktop)
	if err != nil {
		log.Printf("warning: repo scan failed: %v", err)
	}
	log.Printf("found %d repos under %s", len(repos), desktop)

	// Start scheduler: pull all repos every 4 hours, 2s between each
	sched := scheduler.New(repos, 4*time.Hour, 2*time.Second, gitpkg.PullRepo)
	sched.Start()

	// HTTP server — HandlerConfig provides git functions directly
	mux := http.NewServeMux()
	api.RegisterRoutes(mux, api.HandlerConfig{
		Repos: repos,
		GetCommitsFunc: func(repoPath string, n int) ([]api.Commit, error) {
			gitCommits, err := gitpkg.GetCommits(repoPath, n)
			if err != nil {
				return nil, err
			}
			out := make([]api.Commit, len(gitCommits))
			for i, c := range gitCommits {
				out[i] = c.ToAPICommit()
			}
			return out, nil
		},
		GetDiffFunc: func(repoPath, hash string) (string, error) {
			return gitpkg.GetDiff(repoPath, hash)
		},
		GetCommitInfoFunc: func(repoPath, hash string) (string, string, time.Time, error) {
			info, err := gitpkg.GetCommitInfo(repoPath, hash)
			if err != nil {
				return "", "", time.Time{}, err
			}
			return info.Message, info.Author, info.Timestamp, nil
		},
		PullFunc: gitpkg.PullRepo,
	})

	staticFS, err := fs.Sub(staticFiles, "static")
	if err != nil {
		log.Fatalf("failed to sub static fs: %v", err)
	}
	mux.Handle("/", http.FileServer(http.FS(staticFS)))

	logHub := ws.NewHub()
	logHub.Start()
	mux.Handle("/ws/logs", logHub)

	watcherManager := ws.NewWatcherManager(logHub)
	agentHandler := api.NewAgentHandler(agents.ReadCrontab, api.WithOpenCodeBinary(""), api.WithWatcherManager(watcherManager))
	agentHandler.SetRepos(repos)
	mux.HandleFunc("/api/agents", agentHandler.HandleAgents)
	mux.HandleFunc("/api/agents/", agentHandler.HandleAgentAction)
	mux.HandleFunc("/api/models", agentHandler.HandleModels)

	go watchAllAgentLogs(watcherManager, agents.ReadCrontab)

	go watchAllAgentLogs(watcherManager, agents.ReadCrontab)

	addr := ":8080"
	log.Printf("Git Dashboard → http://localhost%s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func watchAllAgentLogs(wm *ws.WatcherManager, readCrontab agents.ReadFunc) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		raw, err := readCrontab()
		if err != nil {
			continue
		}
		ct := agents.ParseCrontab(raw)
		for _, a := range ct.Agents() {
			if a.LogPath != "" {
				wm.StartWatching(a.AgentID(), a.LogPath)
			}
		}
	}
}
