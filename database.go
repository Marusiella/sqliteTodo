package main

import "database/sql"

func SetupDatabase(db *sql.DB) error {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS todos (id INTEGER PRIMARY KEY, content TEXT);")
	if err != nil {
		return err
	}
	return nil
}
