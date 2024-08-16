package logger

import (
	"log"
	"os"
	"path"
	"time"
)

//Ce fichier est responsable des logs du serveur HTTP

const (
	//Chemin pour accèder aux logs
	logsPath = "logs/"
)

var (
	//Différents logger en fonction de la sévérité du message
	WarningLogger *log.Logger
	InfoLogger    *log.Logger
	ErrorLogger   *log.Logger
)

// init crée le dossier où mettre les logs, puis initialise les différents logger
func init() {
	filename := time.Now().Format("2006-1-2_15°4'5''.txt")
	_ = os.MkdirAll(path.Dir(logsPath), os.ModePerm)
	file, err := os.OpenFile(logsPath+filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	InfoLogger = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	WarningLogger = log.New(file, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}
