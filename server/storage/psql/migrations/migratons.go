package migrations

import (
	"database/sql"
	"errors"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v4/stdlib"
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"
)

// RunMigrations apply migrations to database
func RunMigrations(dbDSN string, dbName string) error {
	db, err := sql.Open("pgx", dbDSN)
	if err != nil {
		return err
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Println(err.Error())
		}
	}()
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}
	defer func() {
		if err := driver.Close(); err != nil {
			log.Println(err.Error())
		}
	}()
	ss := getMigrationsFolder()
	log.Println(ss)
	dir := getMigrationsRelPath()
	cc := os.DirFS(dir)
	log.Println(cc)
	m, err := migrate.NewWithDatabaseInstance(dir, dbName, driver)
	if err != nil {
		return err
	}
	if err := m.Up(); err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			return err
		}
	}
	return nil
}

// getFixturesDir returns current file directory.
func getMigrationsFolder() string {
	_, filePath, _, ok := runtime.Caller(1)
	if !ok {
		return ""
	}

	return path.Dir(filePath)
}

func getMigrationsRelPath() string {
	p, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	dir := getMigrationsFolder()
	rel, err := filepath.Rel(p, dir)
	if err != nil {
	}
	rel = "file://" + filepath.ToSlash(rel)
	return rel
}
