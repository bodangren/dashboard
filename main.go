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
				out[i] = api.Commit{
					Hash:      c.Hash,
					Message:   c.Message,
					Body:      c.Body,
					Notes:     c.Notes,
					Author:    c.Author,
					Timestamp: c.Timestamp,
				}
			}
			return out, nil
		},
		GetDiffFunc: func(repoPath, hash string) (string, error) {
			return gitpkg.GetDiff(repoPath, hash)
		},
		PullFunc: gitpkg.PullRepo,
	})

	staticFS, err := fs.Sub(staticFiles, "static")
	if err != nil {
		log.Fatalf("failed to sub static fs: %v", err)
	}
	mux.Handle("/", http.FileServer(http.FS(staticFS)))

	agentHandler := api.NewAgentHandler(agents.ReadCrontab, api.WithOpenCodeBinary(""))
	agentHandler.SetRepos(repos)
	mux.HandleFunc("/api/agents", agentHandler.HandleAgents)
	mux.HandleFunc("/api/agents/", agentHandler.HandleAgentAction)
	mux.HandleFunc("/api/models", agentHandler.HandleModels)

	addr := ":8080"
	log.Printf("Git Dashboard → http://localhost%s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
