package main

import (
	"database/sql"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

type Migrations struct {
	dbConect string
	log      *logrus.Logger
}

func (mg *Migrations) MigrationInit(version int64) {
	db, err := sql.Open("postgres", mg.dbConect)
	if err != nil {
		mg.log.Fatal(err)
	}
	defer db.Close()

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		mg.log.Fatal(err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://./scripts",
		"postgres", driver)
	if err != nil {
		mg.log.Fatal(err)
	}
	err = m.Migrate(uint(int(version)))
	if err != nil {
		mg.log.Fatal(err)
	}
}

func (mg *Migrations) MigrationUp() {
	db, err := sql.Open("postgres", mg.dbConect)
	if err != nil {
		mg.log.Fatal(err)
	}
	defer db.Close()

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		mg.log.Fatal(err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://./scripts",
		"postgres", driver)
	if err != nil {
		mg.log.Fatal(err)
	}
	err = m.Up()
	if err != nil {
		mg.log.Fatal(err)
	}
}

func (mg *Migrations) MigrationDown() {
	db, err := sql.Open("postgres", mg.dbConect)
	if err != nil {
		mg.log.Fatal(err)
	}
	defer db.Close()

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		mg.log.Fatal(err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://./scripts",
		"postgres", driver)
	if err != nil {
		mg.log.Fatal(err)
	}
	err = m.Down()
	if err != nil {
		mg.log.Fatal(err)
	}
}

func (mg *Migrations) MigrationAddOne() {
	db, err := sql.Open("postgres", mg.dbConect)
	if err != nil {
		mg.log.Fatal(err)
	}
	defer db.Close()

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		mg.log.Fatal(err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://./scripts",
		"postgres", driver)
	if err != nil {
		mg.log.Fatal(err)
	}
	err = m.Steps(1)
	if err != nil {
		mg.log.Fatal(err)
	}
}

func (mg *Migrations) MigrationDelOne() {
	db, err := sql.Open("postgres", mg.dbConect)
	if err != nil {
		mg.log.Fatal(err)
	}
	defer db.Close()

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		mg.log.Fatal(err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://./scripts",
		"postgres", driver)
	if err != nil {
		mg.log.Fatal(err)
	}
	err = m.Steps(-1)
	if err != nil {
		mg.log.Fatal(err)
	}
}