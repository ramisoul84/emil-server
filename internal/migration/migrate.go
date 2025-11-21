package migration

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
)

type Migrator struct {
	db     *sqlx.DB
	logger *log.Logger
}

func NewMigrator(db *sqlx.DB) *Migrator {
	return &Migrator{
		db:     db,
		logger: log.New(os.Stdout, "🚀 [migrator] ", log.LstdFlags),
	}
}

// Migrate runs all pending migrations from filesystem
func (m *Migrator) Migrate() error {
	m.logger.Println("Starting database migrations...")

	// Get current version
	currentVersion, err := m.getCurrentVersion()
	if err != nil {
		return fmt.Errorf("failed to get current version: %w", err)
	}

	m.logger.Printf("Current database version: %d", currentVersion)

	// Get all migration files from filesystem
	migrations, err := m.loadMigrationsFromFS()
	if err != nil {
		return fmt.Errorf("failed to load migrations: %w", err)
	}

	if len(migrations) == 0 {
		return fmt.Errorf("no migration files found")
	}

	m.logger.Printf("Found %d migration file(s)", len(migrations))

	// Apply pending migrations
	applied := 0
	for _, migration := range migrations {
		if migration.version > currentVersion {
			m.logger.Printf("Applying migration %d: %s", migration.version, migration.name)

			if err := m.applyMigration(migration); err != nil {
				return fmt.Errorf("failed to apply migration %d: %w", migration.version, err)
			}

			m.logger.Printf("✅ Successfully applied migration %d", migration.version)
			applied++
		}
	}

	if applied == 0 {
		m.logger.Println("✅ Database is up to date - no migrations needed")
	} else {
		m.logger.Printf("✅ All migrations completed! Applied %d migration(s)", applied)
	}

	return nil
}

type migration struct {
	version int
	name    string
	content string
}

func (m *Migrator) loadMigrationsFromFS() ([]migration, error) {
	var migrations []migration

	// Look for migrations in multiple possible locations
	possiblePaths := []string{
		"migration/migrations",          // Your current location
		"internal/migration/migrations", // Standard location
		"database/migrations",           // Alternative location
	}

	var migrationsDir string
	for _, path := range possiblePaths {
		if _, err := os.Stat(path); err == nil {
			migrationsDir = path
			break
		}
	}

	if migrationsDir == "" {
		return nil, fmt.Errorf("migrations directory not found. Checked: %v", possiblePaths)
	}

	m.logger.Printf("📁 Found migrations in: %s", migrationsDir)

	files, err := os.ReadDir(migrationsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read migrations directory: %w", err)
	}

	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(strings.ToLower(file.Name()), ".up.sql") { // ONLY .up.sql files!
			continue
		}

		filename := file.Name()
		m.logger.Printf("📄 Processing: %s", filename)

		// Parse version from filename (001_initial_schema.up.sql -> 1)
		parts := strings.Split(filename, "_")
		if len(parts) < 2 {
			return nil, fmt.Errorf("invalid migration filename format: %s (should be like 001_initial.up.sql)", filename)
		}

		version, err := strconv.Atoi(parts[0])
		if err != nil {
			return nil, fmt.Errorf("invalid version number in %s: %w", filename, err)
		}

		content, err := os.ReadFile(filepath.Join(migrationsDir, filename))
		if err != nil {
			return nil, fmt.Errorf("failed to read %s: %w", filename, err)
		}

		migrations = append(migrations, migration{
			version: version,
			name:    filename,
			content: string(content),
		})

		m.logger.Printf("✅ Loaded migration %d: %s", version, filename)
	}

	// Sort by version
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].version < migrations[j].version
	})

	return migrations, nil
}

func (m *Migrator) getCurrentVersion() (int, error) {
	// Check if schema_version table exists
	var tableExists bool
	err := m.db.Get(&tableExists, `
		SELECT EXISTS (
			SELECT FROM information_schema.tables 
			WHERE table_schema = 'public' 
			AND table_name = 'schema_version'
		)
	`)
	if err != nil {
		return 0, fmt.Errorf("failed to check if schema_version table exists: %w", err)
	}

	if !tableExists {
		m.logger.Println("No schema_version table found - starting fresh")
		return 0, nil
	}

	var version int
	err = m.db.Get(&version, "SELECT COALESCE(MAX(version), 0) FROM schema_version")
	if err != nil {
		return 0, fmt.Errorf("failed to get current version: %w", err)
	}

	return version, nil
}

func (m *Migrator) applyMigration(mig migration) error {
	// Start transaction
	tx, err := m.db.Beginx()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	m.logger.Printf("Executing SQL from %s...", mig.name)
	m.logger.Printf("🔍 SQL Content:\n%s", mig.content)

	// Execute migration
	if _, err := tx.Exec(mig.content); err != nil {
		return fmt.Errorf("SQL execution failed: %w", err)
	}

	// Update schema version
	_, err = tx.Exec(
		"INSERT INTO schema_version (version, description) VALUES ($1, $2)",
		mig.version,
		fmt.Sprintf("Applied: %s", mig.name),
	)
	if err != nil {
		return fmt.Errorf("failed to update schema version: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// VerifyTables checks if required tables exist
func (m *Migrator) VerifyTables() error {
	m.logger.Println("🔍 Verifying database tables...")

	requiredTables := []string{"messages", "schema_version"}

	for _, table := range requiredTables {
		var exists bool
		err := m.db.Get(&exists, `
			SELECT EXISTS (
				SELECT FROM information_schema.tables 
				WHERE table_schema = 'public' 
				AND table_name = $1
			)
		`, table)
		if err != nil {
			return fmt.Errorf("failed to check table %s: %w", table, err)
		}

		if !exists {
			return fmt.Errorf("required table '%s' does not exist", table)
		}

		m.logger.Printf("✅ Table '%s' exists", table)
	}

	m.logger.Println("🎉 All tables verified successfully!")
	return nil
}
