package util

import (
	"errors"
	"github.com/labstack/echo/v4"
	"strconv"
)

const (
	headerTemplateKey  = "Tmpl"
	headerSerieModeKey = "Smode"
	//headerOriginKey    = "Origin"
	//PageType          = "page"
	MainContentType = "main"
	//ResearchContentType   = "research"
	PreviewSetContentType = "previewSet"

	/*
		BpsType           = "bps"
		//PreviewSetType    = "ps"
		ResearchType = "rs"
		BookType     = "b"
		SerieType    = "s"
		ContactType  = "c"
		AuthType     = "a"*/
)

var (
	NoTemplateHeader = errors.New("No template header was given")
)

func GetHeaderTemplate(c echo.Context) (string, error) {
	header := c.Request().Header.Get(headerTemplateKey)
	if header == "" {
		return "", NoTemplateHeader
	}
	return header, nil
}

func IsSerieMode(c echo.Context) bool {
	for key, values := range c.Request().Header {
		if key != headerSerieModeKey {
			continue
		}

		serieMode, err := strconv.ParseBool(values[0])
		if err != nil {
			return false
		}

		return serieMode
	}
	return false
}

/* GetOrigin(c echo.Context) (string, error) {
	header := c.Request().Header.Get(headerOriginKey)
	if header == "" {
		return "", NoTemplateHeader
	}
	return header, nil
}*/
