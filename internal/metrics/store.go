package metrics

import (
	"database/sql"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type AgentRun struct {
	ID          int64
	AgentID     string
	ProjectPath string
	StartedAt   time.Time
	DurationMs  int64
	Success     bool
	ErrorMsg    string
}

type DB struct {
	db *sql.DB
}

func OpenDB(path string) (*DB, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	if err := initSchema(db); err != nil {
		return nil, err
	}
	return &DB{db: db}, nil
}

func (d *DB) Close() error {
	return d.db.Close()
}

func initSchema(db *sql.DB) error {
	schema := `
CREATE TABLE IF NOT EXISTS agent_runs (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    agent_id    TEXT NOT NULL,
    project_path TEXT NOT NULL,
    started_at  TIMESTAMP NOT NULL,
    duration_ms INTEGER NOT NULL,
    success     INTEGER NOT NULL,
    error_msg   TEXT NOT NULL DEFAULT ''
);
CREATE INDEX IF NOT EXISTS idx_agent_runs_agent_id ON agent_runs(agent_id);
CREATE INDEX IF NOT EXISTS idx_agent_runs_project_path ON agent_runs(project_path);
CREATE INDEX IF NOT EXISTS idx_agent_runs_started_at ON agent_runs(started_at);
`
	_, err := db.Exec(schema)
	return err
}

func (d *DB) RecordRun(run *AgentRun) error {
	success := 0
	if run.Success {
		success = 1
	}
	_, err := d.db.Exec(
		`INSERT INTO agent_runs (agent_id, project_path, started_at, duration_ms, success, error_msg)
         VALUES (?, ?, ?, ?, ?, ?)`,
		run.AgentID, run.ProjectPath, run.StartedAt, run.DurationMs, success, run.ErrorMsg,
	)
	return err
}

func (d *DB) QueryRuns(agentID string, limit int, projectPath string) ([]*AgentRun, error) {
	query := `SELECT id, agent_id, project_path, started_at, duration_ms, success, error_msg
              FROM agent_runs WHERE agent_id = ?`
	args := []interface{}{agentID}

	if projectPath != "" {
		query += " AND project_path = ?"
		args = append(args, projectPath)
	}

	query += " ORDER BY started_at DESC"

	if limit > 0 {
		query += " LIMIT ?"
		args = append(args, limit)
	}

	rows, err := d.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var runs []*AgentRun
	for rows.Next() {
		r := &AgentRun{}
		var success int
		if err := rows.Scan(&r.ID, &r.AgentID, &r.ProjectPath, &r.StartedAt, &r.DurationMs, &success, &r.ErrorMsg); err != nil {
			return nil, err
		}
		r.Success = success == 1
		runs = append(runs, r)
	}
	return runs, rows.Err()
}