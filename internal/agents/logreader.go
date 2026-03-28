package agents

import (
	"bufio"
	"fmt"
	"os"
)

func ReadLogFile(path string, n int) (*LogInfo, error) {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &LogInfo{Exists: false}, nil
		}
		return nil, fmt.Errorf("stat log file %s: %w", path, err)
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open log file %s: %w", path, err)
	}
	defer f.Close()

	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read log file %s: %w", path, err)
	}

	total := len(lines)
	truncated := total > n
	start := 0
	if truncated {
		start = total - n
	}

	return &LogInfo{
		Exists:    true,
		LastRun:   info.ModTime(),
		Lines:     lines[start:],
		Truncated: truncated,
	}, nil
}
