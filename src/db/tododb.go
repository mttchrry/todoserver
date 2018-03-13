package db

import (
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // Required so that sqlx has access to the postgres database driver.
)

// TODO - the Postgres client
var Todo *sqlx.DB = sqlx.MustConnect(
	"postgres",
	os.Getenv("DATABASE_URL"))
