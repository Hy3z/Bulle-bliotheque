package dbconvert

import (
	"bb/logger"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

//This package is responsible for database convertion from the old .cvs files to neo4j

const (
	PATH             = "D:/Code/Bulle-bliotheque/src/"
	RESPATH          = "dbconvert/resources/"
	OUTPATH          = "out/"
	BD               = "oldinv_BDs.csv"
	COMICS           = "oldinv_Comics.csv"
	MANGAS           = "oldinv_Mangas.csv"
	ISBN_FOLDER_PATH = PATH + "data/book"
	UTF8_ZERO        = 48
	LOGPATH          = PATH + RESPATH
)

func CreateBDs() {
	err := os.MkdirAll(OUTPATH, os.ModePerm)
	if err != nil {
		logger.ErrorLogger.Panicf("Error creating output directory: %s\n", err)
	}

	inpF, err := os.Open(RESPATH + BD)
	if err != nil {
		logger.ErrorLogger.Panicf("Error opening BDs: %s\n", err)
	}

	outF, err := os.Create(OUTPATH + "BD.cypher")
	if err != nil {
		logger.ErrorLogger.Panicf("Error creating output file: %s\n", err)
	}

	reader := csv.NewReader(inpF)
	reader.FieldsPerRecord = -1
	data, err := reader.ReadAll()
	if err != nil {
		logger.ErrorLogger.Panicf("Error reading data: %s\n", err)
	}

	for j, row := range data {
		isbn13 := row[0]
		if len(isbn13) == 10 {
			isbn13 = ISBN10to13(isbn13)
		}
		link := row[1]
		serie := row[2]
		num := row[3]
		title := row[4]

		cote := row[5]
		cover := row[17]

		if link == "" || isbn13 == "" || title == "" {
			logger.InfoLogger.Printf("Row %d was dismissed\n", j)
			continue
		}

		if serie != "" && num != "" {
			title = serie + " T." + num + ": " + title
		}

		var isbn10, description, publisher, publishedDate string
		var pageCount float64
		res, err := http.Get("https://www.googleapis.com/books/v1/volumes?q=isbn:" + isbn13 + "&key=AIzaSyDrmo1yjzUx7-LsnI42itmHi76ElpIgVps")
		if err != nil {
			logger.ErrorLogger.Printf("Couldn't open url for %s\n", isbn13)
		} else {
			var m map[string]any
			err := json.NewDecoder(res.Body).Decode(&m)
			if err != nil {
				logger.ErrorLogger.Printf("Couldn't decode json: %s\n", err)
			} else {

				totalItems, ok := m["totalItems"]
				if !ok {
					logger.ErrorLogger.Printf("No totalItems field on %s\n", isbn13)
				} else {
					if totalItems.(float64) == 1 {
						items := m["items"].([]any)
						item := items[0].(map[string]any)
						infos := item["volumeInfo"].(map[string]any)

						idents := infos["industryIdentifiers"].([]any)
						for _, indet := range idents {
							ind := indet.(map[string]any)
							t := ind["type"].(string)
							if t == "ISBN_10" {
								isbn10 = ind["identifier"].(string)
								break
							}
						}

						des, ok := infos["description"]
						if ok {
							description = des.(string)
						}
						pub, ok := infos["publisher"]
						if ok {
							publisher = pub.(string)
						}

						pubd, ok := infos["publishedDate"]
						if ok {
							publishedDate = pubd.(string)
						}

						pagc, ok := infos["pageCount"]
						if ok {
							pageCount = pagc.(float64)
						}

					} else {
						logger.WarningLogger.Printf("Found 0 or multiple books with book: %s\n", isbn13)
					}
				}

			}

			res.Body.Close()
		}

		query := fmt.Sprintf("CREATE (:Book {ISBN_13: %q, title: %q, cote: %q,", isbn13, title, cote)
		if isbn10 != "" {
			query += fmt.Sprintf(" ISBN_10: %q,", isbn10)
		}
		if description != "" {
			query += fmt.Sprintf(" description: %q,", description)
		}
		if publisher != "" {
			query += fmt.Sprintf(" publisher: %q,", publisher)
		}
		if publishedDate != "" {
			query += fmt.Sprintf(" publishedDate: %q,", publishedDate)
		}
		if pageCount != 0 {
			query += fmt.Sprintf(" pageCount: %g,", pageCount)
		}
		query = query[:len(query)-1]
		query += "})\n"

		outF.WriteString(query)

		if cover == "" {
			logger.WarningLogger.Printf("Row %d had no cover\n", j)
			continue
		}

		res, err = http.Get(cover)
		if err != nil {
			logger.ErrorLogger.Printf("Error trying get image for row %d: %s\n", j, err)
			continue
		}

		err = os.MkdirAll(OUTPATH+isbn13, os.ModePerm)
		if err != nil {
			logger.ErrorLogger.Printf("Error creating cover folder for row %d: %s\n", j, err)
			continue
		}

		outI, err := os.Create(OUTPATH + isbn13 + "/cover.jpg")
		if err != nil {
			logger.ErrorLogger.Printf("Error creating cover file for row %d: %s\n", j, err)
			continue
		}

		_, err = io.Copy(outI, res.Body)
		if err != nil {
			logger.ErrorLogger.Printf("Error copying cover file for row %d: %s\n", j, err)
			continue
		}

		res.Body.Close()
	}

	inpF.Close()
	outF.Close()
	logger.InfoLogger.Println("Job completed")
}

func RenameIsbnFolders() {
	entries, err := os.ReadDir(ISBN_FOLDER_PATH)
	if err != nil {
		logger.ErrorLogger.Panicf("Could read directory: %s\n", err)
	}
	for _, e := range entries {
		if e.IsDir() && len(e.Name()) == 10 {
			isbn13 := ISBN10to13(e.Name())
			err := os.Rename(ISBN_FOLDER_PATH+"/"+e.Name(), ISBN_FOLDER_PATH+"/"+isbn13)
			if err != nil {
				logger.ErrorLogger.Printf("Couldn't rename %s: %s \n", e.Name(), err)
			}
		}
	}
	logger.InfoLogger.Println("Job completed")
}

func ISBN10to13(isbn10 string) string {
	isbn10 = isbn10[:len(isbn10)-1]
	isbn13 := "978" + isbn10

	var sum uint8 = 0
	for i, _ := range isbn13 {
		c := isbn13[i]
		n := c - UTF8_ZERO
		if i%2 == 1 {
			sum += 3 * n
		} else {
			sum += n
		}
	}
	check := sum % 10

	if check == 0 {
		return isbn13 + "0"
	}

	return isbn13 + string(UTF8_ZERO+(10-check))
}
