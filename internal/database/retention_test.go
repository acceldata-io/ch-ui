package database

import (
	"testing"
	"time"
)

func retentionTimestamp(daysAgo int) string {
	return time.Now().UTC().AddDate(0, 0, -daysAgo).Format(time.RFC3339)
}

func mustExec(t *testing.T, db *DB, query string, args ...interface{}) {
	t.Helper()
	if _, err := db.conn.Exec(query, args...); err != nil {
		t.Fatalf("exec %q: %v", query, err)
	}
}

func countRows(t *testing.T, db *DB, table string) int {
	t.Helper()
	var n int
	if err := db.conn.QueryRow("SELECT COUNT(*) FROM " + table).Scan(&n); err != nil {
		t.Fatalf("count %s: %v", table, err)
	}
	return n
}

func TestRunRetention_PrunesOldRowsKeepsRecent(t *testing.T) {
	db := openTestDB(t)

	// audit_logs (default retention 90 days, no FKs)
	mustExec(t, db,
		`INSERT INTO audit_logs (id, action, username, created_at) VALUES (?, ?, ?, ?)`,
		"audit-old", "user.login", "alice", retentionTimestamp(120))
	mustExec(t, db,
		`INSERT INTO audit_logs (id, action, username, created_at) VALUES (?, ?, ?, ?)`,
		"audit-new", "user.login", "bob", retentionTimestamp(1))

	// alert_events (60 days) + alert_dispatch_jobs (30 days, FK CASCADE to events)
	mustExec(t, db,
		`INSERT INTO alert_events (id, event_type, severity, title, message, status, created_at) VALUES (?, ?, ?, ?, ?, 'new', ?)`,
		"event-old", "schedule.failed", "warning", "old", "old", retentionTimestamp(90))
	mustExec(t, db,
		`INSERT INTO alert_events (id, event_type, severity, title, message, status, created_at) VALUES (?, ?, ?, ?, ?, 'new', ?)`,
		"event-new", "schedule.failed", "warning", "new", "new", retentionTimestamp(1))

	// github_sync_logs (30 days, no FK constraint on connection_id)
	mustExec(t, db,
		`INSERT INTO github_sync_logs (id, connection_id, status, started_at, created_at) VALUES (?, ?, 'success', ?, ?)`,
		"gh-old", "conn-1", retentionTimestamp(45), retentionTimestamp(45))
	mustExec(t, db,
		`INSERT INTO github_sync_logs (id, connection_id, status, started_at, created_at) VALUES (?, ?, 'success', ?, ?)`,
		"gh-new", "conn-1", retentionTimestamp(2), retentionTimestamp(2))

	// pipeline_runs (90 days) + pipeline_run_logs (30 days, FK CASCADE to runs)
	mustExec(t, db,
		`INSERT INTO connections (id, name, tunnel_token) VALUES ('conn-1', 'test', 'tok-1')`)
	mustExec(t, db,
		`INSERT INTO pipelines (id, name, connection_id) VALUES ('pipe-1', 'p', 'conn-1')`)
	mustExec(t, db,
		`INSERT INTO pipeline_runs (id, pipeline_id, status, started_at, created_at) VALUES (?, 'pipe-1', 'success', ?, ?)`,
		"run-old", retentionTimestamp(120), retentionTimestamp(120))
	mustExec(t, db,
		`INSERT INTO pipeline_runs (id, pipeline_id, status, started_at, created_at) VALUES (?, 'pipe-1', 'success', ?, ?)`,
		"run-new", retentionTimestamp(1), retentionTimestamp(1))
	mustExec(t, db,
		`INSERT INTO pipeline_run_logs (id, run_id, level, message, created_at) VALUES (?, 'run-old', 'info', 'old', ?)`,
		"log-old", retentionTimestamp(120))
	mustExec(t, db,
		`INSERT INTO pipeline_run_logs (id, run_id, level, message, created_at) VALUES (?, 'run-new', 'info', 'new', ?)`,
		"log-new", retentionTimestamp(1))

	// gov_schema_changes (180 days, FK to connections)
	mustExec(t, db,
		`INSERT INTO gov_schema_changes (id, connection_id, change_type, database_name, detected_at, created_at) VALUES (?, 'conn-1', 'table_added', 'default', ?, ?)`,
		"sc-old", retentionTimestamp(200), retentionTimestamp(200))
	mustExec(t, db,
		`INSERT INTO gov_schema_changes (id, connection_id, change_type, database_name, detected_at, created_at) VALUES (?, 'conn-1', 'table_added', 'default', ?, ?)`,
		"sc-new", retentionTimestamp(2), retentionTimestamp(2))

	stats := db.RunRetention()

	if stats.LastRunAt == "" {
		t.Fatal("expected LastRunAt to be set")
	}
	if stats.LastError != "" {
		t.Fatalf("unexpected retention error: %s", stats.LastError)
	}

	checks := []struct {
		table string
		want  int
	}{
		{"audit_logs", 1},
		{"alert_events", 1},
		{"github_sync_logs", 1},
		{"pipeline_runs", 1},
		{"pipeline_run_logs", 1},
		{"gov_schema_changes", 1},
	}
	for _, c := range checks {
		if got := countRows(t, db, c.table); got != c.want {
			t.Errorf("%s: expected %d row(s) after retention, got %d", c.table, c.want, got)
		}
	}

	deleted := map[string]int64{
		"audit_logs":         1,
		"alert_events":       1,
		"github_sync_logs":   1,
		"pipeline_runs":      1,
		"pipeline_run_logs":  1,
		"gov_schema_changes": 1,
	}
	for table, want := range deleted {
		if got := stats.RowsDeleted[table]; got != want {
			t.Errorf("stats for %s: expected %d deleted, got %d", table, want, got)
		}
	}
	if stats.TotalDeleted != 6 {
		t.Errorf("expected 6 total deleted, got %d", stats.TotalDeleted)
	}

	// Stats snapshot should be exposed via RetentionStats.
	if snap := db.RetentionStats(); snap.LastRunAt != stats.LastRunAt || snap.TotalDeleted != stats.TotalDeleted {
		t.Errorf("RetentionStats snapshot mismatch: %+v vs %+v", snap, stats)
	}
}

func TestRunRetention_ZeroDisablesPruning(t *testing.T) {
	db := openTestDB(t)

	cfg := DefaultRetentionConfig()
	cfg.AuditLogs = 0
	if err := db.SetRetentionConfig(cfg); err != nil {
		t.Fatalf("SetRetentionConfig: %v", err)
	}

	mustExec(t, db,
		`INSERT INTO audit_logs (id, action, username, created_at) VALUES (?, ?, ?, ?)`,
		"audit-ancient", "user.login", "alice", retentionTimestamp(3000))

	stats := db.RunRetention()

	if got := countRows(t, db, "audit_logs"); got != 1 {
		t.Errorf("expected ancient audit log kept with retention 0, got %d rows", got)
	}
	if _, pruned := stats.RowsDeleted["audit_logs"]; pruned {
		t.Errorf("expected audit_logs to be skipped, got stats entry %v", stats.RowsDeleted)
	}
}

func TestRetentionConfig_RoundTripAndValidation(t *testing.T) {
	db := openTestDB(t)

	// Unset -> defaults.
	cfg, err := db.GetRetentionConfig()
	if err != nil {
		t.Fatalf("GetRetentionConfig: %v", err)
	}
	if cfg != DefaultRetentionConfig() {
		t.Fatalf("expected defaults when unset, got %+v", cfg)
	}

	cfg.AuditLogs = 30
	cfg.GovSchemaChanges = 0
	if err := db.SetRetentionConfig(cfg); err != nil {
		t.Fatalf("SetRetentionConfig: %v", err)
	}
	loaded, err := db.GetRetentionConfig()
	if err != nil {
		t.Fatalf("GetRetentionConfig after set: %v", err)
	}
	if loaded != cfg {
		t.Fatalf("round trip mismatch: %+v vs %+v", loaded, cfg)
	}

	bad := DefaultRetentionConfig()
	bad.AlertEvents = -1
	if err := db.SetRetentionConfig(bad); err == nil {
		t.Fatal("expected error for negative retention")
	}
}

func TestOpen_EnablesIncrementalAutoVacuumOnNewDatabases(t *testing.T) {
	db := openTestDB(t)

	var mode int
	if err := db.conn.QueryRow("PRAGMA auto_vacuum").Scan(&mode); err != nil {
		t.Fatalf("read auto_vacuum: %v", err)
	}
	if mode != 2 { // 2 = INCREMENTAL
		t.Fatalf("expected auto_vacuum=2 (INCREMENTAL) on a new database, got %d", mode)
	}
}
