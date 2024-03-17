package logger

import (
	"log"
	"os"
	"path"
	"time"
)

const (
	logsPath = "logs/"
)

var (
	WarningLogger *log.Logger
	InfoLogger    *log.Logger
	ErrorLogger   *log.Logger
)

func init() {
	filename := time.Now().Format("2006-1-2_15°4'5''.txt")
	os.MkdirAll(path.Dir(logsPath), os.ModePerm)
	file, err := os.OpenFile(logsPath+filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	InfoLogger = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	WarningLogger = log.New(file, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}
