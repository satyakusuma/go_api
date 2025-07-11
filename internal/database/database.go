package database

import (
	"database/sql"
	"fmt"
	"log"
	"github.com/satyakusuma/go-rest-api/internal/config"
	"github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	migrate_mysql "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func Connect() (*sql.DB, error) {
	cfg := mysql.Config{
		User:   config.GetEnv("DB_USER", "root"),
		Passwd: config.GetEnv("DB_PASSWORD", ""),
		Net:    "tcp",
		Addr:   fmt.Sprintf("%s:%s", config.GetEnv("DB_HOST", "localhost"), config.GetEnv("DB_PORT", "3306")),
		DBName: config.GetEnv("DB_NAME", "go_rest_api"),
		AllowNativePasswords: true, 
		ParseTime:            true, 
	}

	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func RunMigrations() error {
	db, err := Connect()
	if err != nil {
		return err
	}
	defer db.Close()

	driver, err := migrate_mysql.WithInstance(db, &migrate_mysql.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://internal/migrations",
		"mysql",
		driver,
	)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	log.Println("Database migrations applied successfully")
	return nil
}