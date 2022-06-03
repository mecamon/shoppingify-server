package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	//_ "github.com/lib/pq"
	"github.com/mecamon/shoppingify-server/config"
)

const maxOpenDbConn = 10
const maxIdleDbConn = 5
const maxDbLifeTime = 5 * time.Minute

func InitDB(conf *config.App) (*sql.DB, error) {
	connStr := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", conf.DBUser, conf.DBPassword, conf.DBHost, conf.DBName)

	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(maxOpenDbConn)
	db.SetMaxIdleConns(maxIdleDbConn)
	db.SetConnMaxLifetime(maxDbLifeTime)

	err = pingDB(db)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func pingDB(conn *sql.DB) error {
	err := conn.Ping()
	if err != nil {
		return err
	}
	log.Println("DB pinged!!!")
	return nil
}
