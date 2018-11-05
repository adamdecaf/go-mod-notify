// Copyright 2018 Adam Shannon
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package database

import (
	"time"
)

type User struct {
	ID      string
	Email   string
	Created time.Time
}

// user_emails user_id (ours), email, created_at, blocked_at

type Project struct {
	ID         string
	UserId     string
	ImportPath string
	WorkScore  int32
	Created    time.Time
	Paused     time.Time
}

// projects: project_id, user_id, import_path, work_score, created_at, paused_at

type Scrape struct {
	ID        string
	Error     string
	WorkScore int32
	Started   time.Time
	Finished  time.Time
}

// scrapes: project_id, error_message, work_score, started_at, finished_at
