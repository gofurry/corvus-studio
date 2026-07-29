package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gofurry/corvus-studio/apps/core/internal/infrastructure/storage/sqlc"
	"github.com/gofurry/corvus-studio/apps/core/migrations"
	"github.com/pressly/goose/v3"
	_ "modernc.org/sqlite"
)

const (
	driverName    = "sqlite"
	busyTimeoutMS = 5000
)

type Store struct {
	db             *sql.DB
	queries        *sqlc.Queries
	schemaVersion  int64
	lastBackupPath string
}

func Open(ctx context.Context, path string) (*Store, error) {
	return open(ctx, path, time.Now)
}

func open(ctx context.Context, path string, now func() time.Time) (*Store, error) {
	absolutePath, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("resolve database path: %w", err)
	}
	absolutePath = filepath.Clean(absolutePath)

	existed, err := nonEmptyFileExists(absolutePath)
	if err != nil {
		return nil, err
	}
	if err := os.MkdirAll(filepath.Dir(absolutePath), 0o700); err != nil {
		return nil, fmt.Errorf("create database directory: %w", err)
	}

	db, err := sql.Open(driverName, dataSourceName(absolutePath))
	if err != nil {
		return nil, fmt.Errorf("open SQLite database: %w", err)
	}
	db.SetMaxOpenConns(4)
	db.SetMaxIdleConns(4)

	closeWithError := func(cause error) (*Store, error) {
		if closeErr := db.Close(); closeErr != nil {
			return nil, errors.Join(cause, fmt.Errorf("close SQLite database: %w", closeErr))
		}
		return nil, cause
	}

	if err := db.PingContext(ctx); err != nil {
		return closeWithError(fmt.Errorf("ping SQLite database: %w", err))
	}

	provider, err := goose.NewProvider(
		goose.DialectSQLite3,
		db,
		migrations.Files(),
		goose.WithDisableGlobalRegistry(true),
		goose.WithLogger(goose.NopLogger()),
	)
	if err != nil {
		return closeWithError(fmt.Errorf("create migration provider: %w", err))
	}

	currentVersion, targetVersion, err := provider.GetVersions(ctx)
	if err != nil {
		return closeWithError(fmt.Errorf("inspect migration versions: %w", err))
	}

	var backupPath string
	if existed && currentVersion < targetVersion {
		backupPath, err = backupBeforeMigration(ctx, db, absolutePath, currentVersion, targetVersion, now().UTC())
		if err != nil {
			return closeWithError(err)
		}
	}

	if _, err := provider.Up(ctx); err != nil {
		return closeWithError(fmt.Errorf("apply database migrations: %w", err))
	}
	schemaVersion, err := provider.GetDBVersion(ctx)
	if err != nil {
		return closeWithError(fmt.Errorf("read database schema version: %w", err))
	}

	store := &Store{
		db:             db,
		queries:        sqlc.New(db),
		schemaVersion:  schemaVersion,
		lastBackupPath: backupPath,
	}
	if err := store.Ping(ctx); err != nil {
		return closeWithError(err)
	}

	return store, nil
}

func (s *Store) Ping(ctx context.Context) error {
	if s == nil || s.db == nil || s.queries == nil {
		return errors.New("storage is not open")
	}
	if err := s.db.PingContext(ctx); err != nil {
		return fmt.Errorf("ping SQLite database: %w", err)
	}
	value, err := s.queries.Ping(ctx)
	if err != nil {
		return fmt.Errorf("run sqlc health query: %w", err)
	}
	if value != 1 {
		return fmt.Errorf("sqlc health query returned %d, want 1", value)
	}
	return nil
}

func (s *Store) SchemaVersion() int64 {
	if s == nil {
		return 0
	}
	return s.schemaVersion
}

func (s *Store) LastBackupPath() string {
	if s == nil {
		return ""
	}
	return s.lastBackupPath
}

func (s *Store) Close() error {
	if s == nil || s.db == nil {
		return nil
	}
	if err := s.db.Close(); err != nil {
		return fmt.Errorf("close SQLite database: %w", err)
	}
	s.db = nil
	s.queries = nil
	return nil
}

func nonEmptyFileExists(path string) (bool, error) {
	info, err := os.Stat(path)
	switch {
	case err == nil:
		if info.IsDir() {
			return false, fmt.Errorf("database path %q is a directory", path)
		}
		return info.Size() > 0, nil
	case errors.Is(err, os.ErrNotExist):
		return false, nil
	default:
		return false, fmt.Errorf("inspect database path %q: %w", path, err)
	}
}

func dataSourceName(path string) string {
	normalizedPath := filepath.ToSlash(path)
	if !strings.HasPrefix(normalizedPath, "/") {
		normalizedPath = "/" + normalizedPath
	}
	uri := url.URL{Scheme: "file", Path: normalizedPath}
	query := url.Values{}
	query.Set("_busy_timeout", strconv.Itoa(busyTimeoutMS))
	query.Set("_foreign_keys", "1")
	query.Set("_journal_mode", "WAL")
	query.Set("_synchronous", "NORMAL")
	uri.RawQuery = query.Encode()
	return uri.String()
}

func backupBeforeMigration(
	ctx context.Context,
	db *sql.DB,
	databasePath string,
	currentVersion int64,
	targetVersion int64,
	now time.Time,
) (string, error) {
	backupDir := filepath.Join(filepath.Dir(databasePath), "backups")
	if err := os.MkdirAll(backupDir, 0o700); err != nil {
		return "", fmt.Errorf("create migration backup directory: %w", err)
	}
	backupName := fmt.Sprintf(
		"corvus-pre-migration-v%d-to-v%d-%s.db",
		currentVersion,
		targetVersion,
		now.Format("20060102T150405.000000000Z"),
	)
	backupPath := filepath.Join(backupDir, backupName)
	if _, err := os.Stat(backupPath); err == nil {
		return "", fmt.Errorf("migration backup path already exists: %q", backupPath)
	} else if !errors.Is(err, os.ErrNotExist) {
		return "", fmt.Errorf("inspect migration backup path: %w", err)
	}

	if _, err := db.ExecContext(ctx, "VACUUM INTO ?", backupPath); err != nil {
		return "", fmt.Errorf("create pre-migration SQLite backup: %w", err)
	}
	return backupPath, nil
}
