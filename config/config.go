package config

import (
	"flag"
	"github.com/mecamon/shoppingify-server/core/loggers"
)

type App struct {
	Loggers    *loggers.CustomLoggers
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
}

var conf = App{}

func Set() {
	conf.Loggers = loggers.Init()
	flag.StringVar(&conf.DBHost, "db-host", "test-database", "DB host")
	flag.StringVar(&conf.DBPort, "db-port", "5432", "Port on which runs the db")
	flag.StringVar(&conf.DBUser, "db-user", "developer", "Database's username")
	flag.StringVar(&conf.DBPassword, "db-password", "123456789", "Database's user password")
	flag.StringVar(&conf.DBName, "db-name", "shoppingify-test", "Database's name")
	flag.Parse()
}

func Get() *App {
	return &conf
}
