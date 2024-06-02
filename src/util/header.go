package util

import (
	"errors"
	"github.com/labstack/echo/v4"
	"strconv"
)

const (
	headerTemplateKey  = "Tmpl"
	headerSerieModeKey = "SerieMode"
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
	NoTemplateHeader       = errors.New("No template header was given")
	MultipleTemplateHeader = errors.New("Multiple template were given in header")
	MultipleSerieMode      = errors.New("Multiple serie modes were given in header")
)

func GetHeaderTemplate(c echo.Context) (string, error) {
	for key, values := range c.Request().Header {
		if key != headerTemplateKey {
			continue
		}
		if len(values) == 1 {
			return values[0], nil
		}
		return "", MultipleTemplateHeader
	}
	return "", NoTemplateHeader
}

func IsSerieMode(c echo.Context) (bool, error) {
	for key, values := range c.Request().Header {
		if key != headerSerieModeKey {
			continue
		}

		if len(values) > 1 {
			return false, MultipleSerieMode
		}

		serieMode, err := strconv.ParseBool(values[0])
		if err != nil {
			return false, err
		}

		return serieMode, nil
	}
	return false, nil
}
