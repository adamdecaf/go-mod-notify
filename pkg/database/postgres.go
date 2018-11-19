// Copyright 2018 Adam Shannon
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/lib/pq"
)

var (
	minimumRescrapeTime = -12 * time.Hour

	timeZero = time.Unix(0, 0)
)

func PostgresFromEnv() (*sql.DB, error) {
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

type PostgresRepository struct {
	DB *sql.DB
}

// - one transaction
//   - grab N projects from `projects` table (by nonce) without 'started_at is not null and finished_at is null'
//   - insert new rows into scrapes
//
// projects: project_id, user_id, import_path, work_score, created_at, paused_at
// scrapes: project_id, error_message, work_score, started_at, finished_at

func (r *PostgresRepository) ScrapeableProjects(workScore uint32, limit int) ([]*Project, error) {
	tx, err := r.DB.Begin()
	if err != nil {
		return nil, fmt.Errorf("postgres: ScrapeableProjects: problem starting transaction: %v", err)
	}

	// TODO(adam): look at p.paused_at
	// TODO(adam): this needs to be a left outer join (we want to exclude those in scrapes and include those not in scrapes)
	// but then also find stale projects
	stmt, err := tx.Prepare(`select p.project_id, p.user_id, p.import_path, p.work_score, p.created_at, p.paused_at
from projects as p inner join scrapes as s on p.project_id = s.project_id
where s.started_at is null and s.finished_at is null -- can't be started by another worker already
and work_score > ? or work_score <= ?                -- close to our workscore
and s.finished_at < ?                                -- last scrape must be 'N hours ago'
order by s.finished_at asc                           -- grab the most stale first
limit ?`)
	if err != nil {
		return nil, fmt.Errorf("postgres: ScrapeableProjects: problem with prepare: %v", err)
	}
	defer stmt.Close()

	oldEnough := time.Now().Add(minimumRescrapeTime)
	rows, err := tx.Query(fmt.Sprintf("%d", workScore), workScore, oldEnough, limit)
	if err != nil {
		return nil, fmt.Errorf("postgres: ScrapeableProjects: problem with query: %v", err)
	}

	var projects []*Project
	for rows.Next() {
		p := &Project{}
		err = rows.Scan(&p.ID, &p.UserId, &p.ImportPath, &p.WorkScore, &p.Created, &p.Paused)
		if err != nil {
			// TODO(adam): log here
			continue
		}
		projects = append(projects, p)
	}

	// Insert our scrapability into scrapes table // TODO(adam): pull this into a helper that takes a *sql.Tx
	// scrapes is table name, others are column names
	stmt, err = tx.Prepare(pq.CopyIn("scrapes", "project_id", "error_message", "work_score", "started_at", "finished_at"))
	if err != nil {
		log.Fatal(err)
	}
	for i := range projects {
		p := projects[i]
		if _, err = stmt.Exec(p.ID, "", workScore, time.Now(), timeZero); err != nil {
			return nil, fmt.Errorf("postgres: ScrapeableProjects: problem inserting scrapes: %v (rollback error=%v)", err, tx.Rollback())
		}
	}
	if _, err := stmt.Exec(); err != nil {
		return nil, fmt.Errorf("postgres: ScrapeableProjects: problem with copy: %v", err)
	}
	if err := stmt.Close(); err != nil {
		return nil, fmt.Errorf("postgres: ScrapeableProjects: problem closing statement: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("postgres: ScrapeableProjects: problem with commit: %v", err)
	}
	return projects, nil
}
