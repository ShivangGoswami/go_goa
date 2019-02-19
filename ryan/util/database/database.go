package database

import (
	"database/sql"

	_ "github.com/go-goracle/goracle"
)

func Connect() (*sql.DB, error) {
	return sql.Open("goracle", "system/1001289@localhost:1521/xe")
}
