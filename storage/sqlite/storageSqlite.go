package sqlite

import (
	"database/sql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
	"log"

)

type Storage struct {
	Sql *sql.DB
}

func New() *Storage  {

	db,err := sql.Open("sqlite3", ":memory:")
	if err != nil{
		log.Fatalf("Cannot connect to DB: %v\n", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("DB did not respond to ping: %v\n", err)
	}

	log.Println("Connect to DB successful")

	doMigrate(db)

	return &Storage{
		Sql: db,
	}
}

func doMigrate(db *sql.DB) {
	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		log.Fatalf("Instance error: %v\n", err)
	}
	fileMigration, err := (&file.File{}).Open("file://migrations")
	if err != nil {
		log.Fatalf("Openning file error: %v\n", err)
	}

	m,err := migrate.NewWithInstance("file", fileMigration, "myDB", driver)
	if err != nil {
		log.Fatalf("Migrate error: %v\n", err)
	}

	if err = m.Up(); err != nil {
		log.Fatalf("Migrate UP error: %v\n", err)
	}

	log.Println("Migrate UP done with success")
}