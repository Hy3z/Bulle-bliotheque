package dbconvert

import (
	"bb/database"
	"bb/logger"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
)

const (
	API_URL    = "http://bookcover.longitood.com/bookcover/"
	ISBN_PATH  = "D:/Code/Bulle-bliotheque/src/data/book"
	SERIE_PATH = "D:/Code/Bulle-bliotheque/src/data/serie"
	URL_KEY    = "url"
)

func DownloadCovers() {
	res, _ := database.Query(context.Background(), "MATCH (b:Book) return b.ISBN_13", map[string]any{})
	for _, rec := range res.Records {
		isbn13, _ := rec.Values[0].(string)
		_, err := os.Stat(ISBN_PATH + "/" + isbn13)
		_, err2 := os.Stat(ISBN_PATH + "/" + isbn13 + "/cover.jpg")
		if os.IsExist(err) && os.IsExist(err2) {
			continue
		}

		err = os.MkdirAll(ISBN_PATH+"/"+isbn13, os.ModePerm)
		if err != nil {
			logger.ErrorLogger.Printf("Error creating directory: %s\n", err)
			continue
		}

		res, err := http.Get(API_URL + isbn13)
		if err != nil {
			logger.ErrorLogger.Printf("HTTP error: %s\n", err)
			continue
		}
		defer res.Body.Close()

		var m map[string]any
		if json.NewDecoder(res.Body).Decode(&m) != nil {
			logger.ErrorLogger.Printf("JSON error: %s\n", err)
			continue
		}

		url, ok := m[URL_KEY].(string)
		if !ok {
			logger.InfoLogger.Printf("Couldn't find url for %s\n", isbn13)
			continue
		}

		res, err = http.Get(url)
		if err != nil {
			logger.ErrorLogger.Printf("HTTP error: %s\n", err)
			continue
		}
		defer res.Body.Close()

		outI, err := os.Create(ISBN_PATH + "/" + isbn13 + "/cover.jpg")
		if err != nil {
			logger.ErrorLogger.Printf("Error creating cover file for %s: %s\n", isbn13, err)
			continue
		}

		io.Copy(outI, res.Body)

		logger.InfoLogger.Printf("Created cover for %s\n", isbn13)
	}
}

func GetMissingCovers() (isbns []string) {
	res, _ := database.Query(context.Background(), "MATCH (b:Book) return b.ISBN_13", map[string]any{})
	for _, rec := range res.Records {
		isbn13, _ := rec.Values[0].(string)
		_, err1 := os.Stat(ISBN_PATH + "/" + isbn13)
		_, err2 := os.Stat(ISBN_PATH + "/" + isbn13 + "/cover.jpg")
		if errors.Is(err1, os.ErrNotExist) || errors.Is(err2, os.ErrNotExist) {
			isbns = append(isbns, isbn13)
		}
	}
	return
}

func PrintMissingCovers() {
	logger.InfoLogger.Println("Begin")
	for _, isbn13 := range GetMissingCovers() {
		logger.InfoLogger.Printf("%s is missing it's cover\n", isbn13)
	}
	logger.InfoLogger.Println("End")
}

func SerieCoverFromBook() {
	logger.InfoLogger.Println("Begin")
	res, err := database.Query(context.Background(), ""+
		"match(s:Serie)<-[r:PART_OF]-(b:Book) "+
		"with s,min(r.opus) as minopus "+
		"match (s:Serie)<-[rp:PART_OF{opus:minopus}]-(bp:Book) "+
		"return s.UUID,bp.ISBN_13", map[string]any{})
	if err != nil {
		logger.ErrorLogger.Printf("Couldn't query database: %s\n", err)
		return
	}

	for _, rec := range res.Records {
		uuid, _ := rec.Values[0].(string)
		isbn13, _ := rec.Values[1].(string)
		from, err := os.Open(ISBN_PATH + "/" + isbn13 + "/cover.jpg")
		if err != nil {
			logger.ErrorLogger.Printf("Couldn't open from: %s\n", err)
			continue
		}
		os.Mkdir(SERIE_PATH+"/"+uuid, os.ModePerm)
		to, err := os.Create(SERIE_PATH + "/" + uuid + "/cover.jpg")
		if err != nil {
			logger.ErrorLogger.Printf("Couldn't create to: %s\n", err)
			continue
		}
		io.Copy(to, from)
		logger.InfoLogger.Printf("Created cover for %s\n", uuid)
	}
	logger.InfoLogger.Println("End")
}
