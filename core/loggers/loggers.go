package loggers

import (
	"log"
	"os"
)

type CustomLoggers struct {
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
}

func Init() *CustomLoggers {
	file, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	infoLogger := log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	warningLogger := log.New(file, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	errorLogger := log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

	loggers := CustomLoggers{
		Info:    infoLogger,
		Warning: warningLogger,
		Error:   errorLogger,
	}
	return &loggers
}
