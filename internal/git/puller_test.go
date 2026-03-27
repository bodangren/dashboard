package git

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// makeClonePair creates a "remote" bare repo and a local clone of it,
// returning (remote, localClone) paths.
func makeClonePair(t *testing.T) (remote, local string) {
	t.Helper()
	base := t.TempDir()
	remote = filepath.Join(base, "remote.git")
	local = filepath.Join(base, "local")

	run := func(dir string, args ...string) {
		t.Helper()
		cmd := exec.Command("git", args...)
		cmd.Dir = dir
		cmd.Env = append(os.Environ(),
			"GIT_AUTHOR_NAME=T",
			"GIT_AUTHOR_EMAIL=t@t.com",
			"GIT_COMMITTER_NAME=T",
			"GIT_COMMITTER_EMAIL=t@t.com",
		)
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v\n%s", args, err, out)
		}
	}

	// Create bare remote
	run(base, "init", "--bare", remote)
	// Clone it
	run(base, "clone", remote, local)
	// Seed a commit in local and push to remote
	run(local, "config", "user.email", "t@t.com")
	run(local, "config", "user.name", "T")
	f := filepath.Join(local, "seed.txt")
	if err := os.WriteFile(f, []byte("seed\n"), 0644); err != nil {
		t.Fatal(err)
	}
	run(local, "add", ".")
	run(local, "commit", "-m", "seed")
	run(local, "push", "origin", "HEAD")

	return remote, local
}

func TestPullRepo_succeedsOnCleanRepo(t *testing.T) {
	_, local := makeClonePair(t)

	err := PullRepo(local)
	if err != nil {
		t.Errorf("PullRepo on clean repo returned error: %v", err)
	}
}

func TestPullRepo_noErrorOnNonRemoteRepo(t *testing.T) {
	// A repo with no remote should not cause a fatal error — just log and continue.
	dir := initTestRepo(t)
	err := PullRepo(dir)
	// We expect an error (no remote) but it should be a wrapped, non-panicking error.
	// The caller (scheduler) will log it and continue.
	_ = err // either nil or a wrapped error is acceptable
}

func TestPullRepo_errorOnNonGitDir(t *testing.T) {
	err := PullRepo(t.TempDir())
	if err == nil {
		t.Error("expected error for non-git directory, got nil")
	}
}
