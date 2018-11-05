// Copyright 2018 Adam Shannon
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package database

import (
	"database/sql"
	"fmt"
	"os"
)

func createPostgresConnection() (*sql.DB, error) {
	user, pass := os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASS")
	host, dbname := os.Getenv("POSTGRES_HOSTNAME"), os.Getenv("POSTGRES_DATABASE")
	if user == "" || pass == "" || host == "" || dbname == "" {
		return nil, fmt.Errorf("missing postgres creds (user=%s) (pass:%v) (host=%s) (database:%s)", user, pass != "", host, dbname)
	}
	conn := fmt.Sprintf("user=%s password=%s host=%s dbname=%s sslmode=verify-full", user, pass, host, dbname)
	db, err := sql.Open("postgres", conn)
	if err != nil {
		return nil, fmt.Errorf("problem creating postgres connection to %s (user=%s) (pass:%v) (database:%s)", host, user, pass != "", dbname)
	}
	return db, nil
}
