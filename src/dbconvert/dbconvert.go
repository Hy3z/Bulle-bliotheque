package dbconvert

import (
	"bb/logger"
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"os"
)

//This package is responsible for database convertion from the old .cvs files to neo4j

const(
	PATH = "D:/Code/Bulle-bliotheque/src/"
	RESPATH = "dbconvert/resources/"
	OUTPATH = "out/"
	BD = "oldinv_BDs.csv"
	COMICS = "oldinv_Comics.csv"
	MANGAS = "oldinv_Mangas.csv"
)


func CreateBDs() {
	err := os.MkdirAll(OUTPATH, os.ModePerm)
	if err != nil {
		logger.ErrorLogger.Panicf("Error creating output directory: %s \n", err)
	}

	inpF, err := os.Open(RESPATH+BD)
	if err != nil {
		logger.ErrorLogger.Panicf("Error opening BDs: %s \n", err)
	}

	outF, err := os.Create(OUTPATH+"BD.cypher")
	if err != nil {
		logger.ErrorLogger.Panicf("Error creating output file: %s \n", err)
	}

	reader := csv.NewReader(inpF)
	reader.FieldsPerRecord = -1
	data, err := reader.ReadAll()
	if err != nil {
		logger.ErrorLogger.Panicf("Error reading data: %s \n", err)
	}

	for j, row := range data {
		isbn := row[0]
		link := row[1]
		serie := row[2]
		num := row[3]
		title := row[4]

		cote := row[5]
		cover := row[17]

		if link == "" || isbn == "" || title == "" {
			logger.InfoLogger.Printf("Row %d was dismissed \n", j)
			continue
		}

		if serie != "" && num != "" {
			title = serie+" T."+num+": "+title
		}

		query := fmt.Sprintf("CREATE (b:Book {ISBN: %q, link: %q, title: %q, cote: %q, cover: %q})\n", isbn, link, title, cote, cover)
		outF.WriteString(query)

		if cover == ""{
			logger.WarningLogger.Printf("Row %d had no cover \n", j)
			continue
		}

		res, err := http.Get(cover)
		if err != nil {
			logger.ErrorLogger.Printf("Error trying get image for row %d: %s \n", j, err)
			continue
		}

		err = os.MkdirAll(OUTPATH+isbn, os.ModePerm)
		if err != nil {
			logger.ErrorLogger.Printf("Error creating cover folder for row %d: %s \n", j, err)
			continue
		}

		outI, err := os.Create(OUTPATH+isbn+"/cover.jpg")
		if err != nil {
			logger.ErrorLogger.Printf("Error creating cover file for row %d: %s \n", j, err)
			continue
		}

		_,err = io.Copy(outI, res.Body)
		if err != nil {
			logger.ErrorLogger.Printf("Error copying cover file for row %d: %s \n", j, err)
			continue
		}

		res.Body.Close()
	}

	inpF.Close()
	outF.Close()
}

