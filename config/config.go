package config

import "flag"

type App struct {
	DBHost string
	DBport string
	DBUser string
	DBPassword string
	DBName string
}

var conf = App{}

func Set() {
	flag.StringVar(&conf.DBHost, "db-host", "postgres_db", "DB host")
	flag.StringVar(&conf.DBport, "db-port", "5432", "Port on which runs the db")
	flag.StringVar(&conf.DBUser, "db-user", "", "Database's username")
	flag.StringVar(&conf.DBPassword, "db-password", "", "Database's user password")
	flag.StringVar(&conf.DBName, "db-name", "", "Database's name")
	flag.Parse()
}

func Get() *App {
	return &conf
}