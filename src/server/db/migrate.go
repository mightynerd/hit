package db

import (
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type Migrator struct {
	dbUrl string
}

func NewMigrator(dbUrl string) *Migrator {
	migrator := &Migrator{
		dbUrl: dbUrl,
	}

	return migrator
}

func (m *Migrator) Migrate() {
	migration, err := migrate.New("file://db/migrations", m.dbUrl)
	if err != nil {
		log.Fatal("migration failed", err)
	}
	err = migration.Up()
	if err != nil && err != migrate.ErrNoChange {
		log.Fatal("migration failed", err)
	}
}
