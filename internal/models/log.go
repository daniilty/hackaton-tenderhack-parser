package models

import "time"

type Log struct {
	ID         string     `db:"id" json:"id"`
	TS         time.Time  `db:"ts" json:"ts"`
	CategoryID int        `db:"category_id" json:"category_id"`
	Data       string     `db:"data" json:"data"`
	ResolvedAt *time.Time `db:"resolved_at" json:"resolved_at"`
}
