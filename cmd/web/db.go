package main

import (
	"database/sql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func createDatabasePool(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn+"?parseTime=true")
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	migrateTables(db)

	return db, nil
}

func migrateTables(db *sql.DB) {
	driver, _ := mysql.WithInstance(db, &mysql.Config{})
	m, _ := migrate.NewWithDatabaseInstance(
		"file://../../pkg/models/mysql/schemas",
		"mysql",
		driver,
	)

	m.Steps(2)
	//tableUsers := `
	//	CREATE TABLE IF NOT EXISTS user (
	//		id SMALLINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
	//		username VARCHAR(80) UNIQUE,
	//		hashed_password VARCHAR(60) NOT NULL
	//	);`
	//_, err := db.Exec(tableUsers)
	//
	//if err != nil {
	//	errLogger.Fatalln(err)
	//}
	//
	//tableTodo := `
	//	CREATE TABLE IF NOT EXISTS todo (
	//		id SMALLINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
	//		title VARCHAR(255) NOT NULL,
	//		done BOOlEAN NOT NULL DEFAULT FALSE,
	//		created_at DATETIME NOT NULL,
	//		expires DATETIME NOT NULL
	//	);`
	//_, err2 := db.Exec(tableTodo)
	//
	//if err2 != nil {
	//	errLogger.Fatalln(err)
	//}
}
