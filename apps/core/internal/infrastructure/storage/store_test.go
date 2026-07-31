package storage

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"testing"
	"time"

	_ "modernc.org/sqlite"
)

func TestOpenMigratesAndReopensIdempotently(t *testing.T) {
	ctx := context.Background()
	databasePath := filepath.Join(t.TempDir(), "nested", "corvus.db")

	first, err := Open(ctx, databasePath)
	if err != nil {
		t.Fatalf("first open: %v", err)
	}
	if first.SchemaVersion() != 2 {
		t.Fatalf("first schema version = %d, want 2", first.SchemaVersion())
	}
	if first.LastBackupPath() != "" {
		t.Fatalf("fresh database backup = %q, want none", first.LastBackupPath())
	}
	assertPragmas(t, first.db)
	if err := first.Close(); err != nil {
		t.Fatalf("close first store: %v", err)
	}

	second, err := Open(ctx, databasePath)
	if err != nil {
		t.Fatalf("second open: %v", err)
	}
	t.Cleanup(func() { _ = second.Close() })
	if second.SchemaVersion() != 2 {
		t.Fatalf("second schema version = %d, want 2", second.SchemaVersion())
	}
	if second.LastBackupPath() != "" {
		t.Fatalf("idempotent reopen backup = %q, want none", second.LastBackupPath())
	}
}

func TestOpenBacksUpExistingDatabaseBeforeMigration(t *testing.T) {
	ctx := context.Background()
	databasePath := filepath.Join(t.TempDir(), "corvus.db")
	seedDatabase(t, databasePath)

	fixedTime := time.Date(2026, time.July, 29, 10, 0, 0, 123, time.UTC)
	store, err := open(ctx, databasePath, func() time.Time { return fixedTime })
	if err != nil {
		t.Fatalf("open existing database: %v", err)
	}
	t.Cleanup(func() { _ = store.Close() })

	wantBackup := filepath.Join(
		filepath.Dir(databasePath),
		"backups",
		"corvus-pre-migration-v0-to-v2-20260729T100000.000000123Z.db",
	)
	if store.LastBackupPath() != wantBackup {
		t.Fatalf("backup path = %q, want %q", store.LastBackupPath(), wantBackup)
	}
	if _, err := os.Stat(wantBackup); err != nil {
		t.Fatalf("stat backup: %v", err)
	}

	backup, err := sql.Open(driverName, dataSourceName(wantBackup))
	if err != nil {
		t.Fatalf("open backup: %v", err)
	}
	t.Cleanup(func() { _ = backup.Close() })
	var value string
	if err := backup.QueryRowContext(ctx, "SELECT value FROM seed_data WHERE id = 1").Scan(&value); err != nil {
		t.Fatalf("read seed row from backup: %v", err)
	}
	if value != "preserve me" {
		t.Fatalf("backup seed value = %q", value)
	}
}

func TestDataSourceNameEscapesPathAndSetsPragmas(t *testing.T) {
	path := filepath.Join(t.TempDir(), "path with spaces", "corvus.db")
	dsn := dataSourceName(path)
	if dsn == path {
		t.Fatalf("DSN was not converted to URI: %q", dsn)
	}

	store, err := Open(context.Background(), path)
	if err != nil {
		t.Fatalf("open path with spaces: %v", err)
	}
	t.Cleanup(func() { _ = store.Close() })
	assertPragmas(t, store.db)
}

func seedDatabase(t *testing.T, path string) {
	t.Helper()
	db, err := sql.Open(driverName, dataSourceName(path))
	if err != nil {
		t.Fatalf("open seed database: %v", err)
	}
	if _, err := db.Exec(`
CREATE TABLE seed_data (id INTEGER PRIMARY KEY, value TEXT NOT NULL) STRICT;
INSERT INTO seed_data (id, value) VALUES (1, 'preserve me');
`); err != nil {
		_ = db.Close()
		t.Fatalf("seed database: %v", err)
	}
	if err := db.Close(); err != nil {
		t.Fatalf("close seed database: %v", err)
	}
}

func assertPragmas(t *testing.T, db *sql.DB) {
	t.Helper()
	var foreignKeys int
	if err := db.QueryRow("PRAGMA foreign_keys").Scan(&foreignKeys); err != nil {
		t.Fatalf("read foreign_keys pragma: %v", err)
	}
	if foreignKeys != 1 {
		t.Fatalf("foreign_keys = %d, want 1", foreignKeys)
	}
	var journalMode string
	if err := db.QueryRow("PRAGMA journal_mode").Scan(&journalMode); err != nil {
		t.Fatalf("read journal_mode pragma: %v", err)
	}
	if journalMode != "wal" {
		t.Fatalf("journal_mode = %q, want wal", journalMode)
	}
}
