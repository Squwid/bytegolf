package sqldb

import (
	"crypto/tls"
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

var (
	host     = os.Getenv("BGPG_HOST")
	port     = 5432
	user     = "postgres"
	password = os.Getenv("BGPG_PASSWORD")
	dbname   = os.Getenv("BGPG_DBNAME")
)

var DB *bun.DB

func Open() error {
	sqldb := sql.OpenDB(pgdriver.NewConnector(
		pgdriver.WithAddr(fmt.Sprintf("%s:%v", host, port)),
		pgdriver.WithTLSConfig(&tls.Config{InsecureSkipVerify: true}),
		pgdriver.WithInsecure(true),
		pgdriver.WithDatabase(dbname),
		pgdriver.WithUser(user),
		pgdriver.WithPassword(password),
		pgdriver.WithTimeout(5*time.Second),
		pgdriver.WithDialTimeout(5*time.Second),
		pgdriver.WithReadTimeout(5*time.Second),
		pgdriver.WithWriteTimeout(5*time.Second),
	))
	if err := sqldb.Ping(); err != nil {
		return err
	}

	DB = bun.NewDB(sqldb, pgdialect.New())

	return DB.Ping()
}

func Close() error { return DB.Close() }
