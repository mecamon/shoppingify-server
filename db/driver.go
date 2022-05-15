package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/mecamon/shoppingify-server/config"
)


func InitDB(conf *config.App) (*sql.DB, error) {

	connStr := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", conf.DBUser, conf.DBPassword, conf.DBHost, conf.DBName) 

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	return db, nil
}