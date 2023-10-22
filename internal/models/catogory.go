package models

type Category struct {
	ID        int    `db:"id" json:"id"`
	GroupID   int    `db:"group_id" json:"group_id"`
	Reg       string `db:"reg" json:"reg"`
	GroupName string `db:"group_name" json:"group_name"`
}

type CategoryGroup struct {
	ID       int    `db:"id" json:"id"`
	Name     string `db:"name" json:"name"`
	Severity int    `db:"severity" json:"severity"`
}
