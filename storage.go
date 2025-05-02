// Copyright (c) 2025 Grigoriy Efimov
//
// Licensed under the MIT License. See LICENSE file in the project root for details.

package main

import (
	"database/sql"

	ink "github.com/CatInBeard/inkview"
)

func openDb() (*sql.DB, error) {
	return sql.Open("sqlite3", ink.UserData+"/app-manager.sql")
}

func CreateInitDb() error {
	db, err := openDb()
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS settings (
			id INTEGER PRIMARY KEY,
			key VARCHAR(255),
			val VARCHAR(255)
		);
	`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS installed_apps (
			id INTEGER PRIMARY KEY,
			install_date DATETIME,
			version VARCHAR(255),
			package_name VARCHAR(255)
		);
	`)
	if err != nil {
		return err
	}

	return nil
}
