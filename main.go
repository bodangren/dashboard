package main

import (
	"embed"
	"encoding/json"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"dashboard/internal/agents"
	"dashboard/internal/ai"
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
		GetHealthFunc: func(repoPath string) (api.RepoHealth, error) {
			h, err := gitpkg.GetRepoHealth(repoPath)
			if err != nil {
				return api.RepoHealth{}, err
			}
			return api.RepoHealth{
				Dirty: api.DirtyStatus{
					Modified:  h.Dirty.Modified,
					Staged:    h.Dirty.Staged,
					Untracked: h.Dirty.Untracked,
					Total:     h.Dirty.Total,
				},
				Divergence: api.BranchDivergence{
					Ahead:  h.Divergence.Ahead,
					Behind: h.Divergence.Behind,
				},
				StaleBranches: api.StaleBranchInfo{
					Count:    h.StaleBranches.Count,
					Branches: h.StaleBranches.Branches,
				},
			}, nil
		},
		GetBranchesFunc:    gitpkg.GetBranches,
		CreateBranchFunc:   gitpkg.CreateBranch,
		DeleteBranchFunc:   gitpkg.DeleteBranch,
		CheckoutBranchFunc: gitpkg.CheckoutBranch,
		GetStashFunc: func(repoPath string) ([]api.StashEntry, error) {
			entries, err := gitpkg.GetStashList(repoPath)
			if err != nil {
				return nil, err
			}
			out := make([]api.StashEntry, len(entries))
			for i, e := range entries {
				out[i] = api.StashEntry{
					Index:     e.Index,
					Message:   e.Message,
					Author:    e.Author,
					Timestamp: e.Timestamp,
				}
			}
			return out, nil
		},
		ApplyStashFunc: gitpkg.ApplyStash,
		DropStashFunc:  gitpkg.DropStash,
	})

	staticFS, err := fs.Sub(staticFiles, "static")
	if err != nil {
		log.Fatalf("failed to sub static fs: %v", err)
	}
	mux.Handle("/", http.FileServer(http.FS(staticFS)))

	logHub := ws.NewHub()
	logHub.Start()
	mux.Handle("/ws/logs", logHub)

	activityHub := ws.NewActivityHub()
	activityHub.Start()
	mux.Handle("/ws/activity", activityHub)

	summarizer, err := ai.DefaultSummarizer()
	if err != nil {
		log.Printf("warning: failed to create summarizer: %v", err)
	}
	var activityEnhancer *ai.ActivityEnhancer
	if summarizer != nil {
		activityEnhancer = ai.NewActivityEnhancer(summarizer)
	}

	activityHandler := api.NewActivityHandler(
		api.WithActivityRepos(repos),
		api.WithActivityGetCommits(func(repoPath string, n int) ([]api.Commit, error) {
			gitCommits, err := gitpkg.GetCommits(repoPath, n)
			if err != nil {
				return nil, err
			}
			out := make([]api.Commit, len(gitCommits))
			for i, c := range gitCommits {
				out[i] = c.ToAPICommit()
			}
			return out, nil
		}),
		api.WithActivityReadCrontab(agents.ReadCrontab),
		api.WithActivityHub(activityHub),
		api.WithActivityEnhancer(activityEnhancer),
	)
	mux.HandleFunc("/api/activity", activityHandler.HandleActivity)

	watcherManager := ws.NewWatcherManager(logHub)
	agentHandler := api.NewAgentHandler(agents.ReadCrontab,
		api.WithOpenCodeBinary(""),
		api.WithWatcherManager(watcherManager),
		api.WithAgentCompleteCallback(func(agentID string, exitCode int, lastError string) {
			var status string
			if exitCode == 0 {
				status = "completed"
			} else {
				status = "failed"
			}
			agentName := agentID
			if idx := strings.LastIndex(agentID, ":"); idx >= 0 {
				if modelIdx := strings.LastIndex(agentID[:idx], ":"); modelIdx >= 0 {
					agentName = agentID[modelIdx+1 : idx]
				}
			}
			meta, err := json.Marshal(api.AgentEventMetadata{
				AgentID:   agentID,
				AgentName: agentName,
				Status:    status,
				ExitCode:  exitCode,
				Error:     lastError,
			})
			if err != nil {
				return
			}
			activityHandler.RecordAgentEvent(api.ActivityEvent{
				ID:        "agent-" + agentID,
				Type:      api.EventTypeAgent,
				Repo:      "",
				Message:   "Agent " + status,
				Timestamp: time.Now(),
				Metadata:  meta,
			})
		}),
	)
	agentHandler.SetRepos(repos)
	mux.HandleFunc("/api/agents", agentHandler.HandleAgents)
	mux.HandleFunc("/api/agents/", agentHandler.HandleAgentAction)
	mux.HandleFunc("/api/models", agentHandler.HandleModels)

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
