package database

import (
	"database/sql"
	"os"
	"path/filepath"
)

func RunMigrations(db *sql.DB) error {
	migrationsDir := "./migrations"

	entries, err := os.ReadDir(migrationsDir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if filepath.Ext(entry.Name()) == ".sql" {
			content, err := os.ReadFile(filepath.Join(migrationsDir, entry.Name()))
			if err != nil {
				return err
			}

			if _, err := db.Exec(string(content)); err != nil {
				return err
			}
		}
	}

	return nil
}
