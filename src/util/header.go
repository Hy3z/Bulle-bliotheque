package util

import (
	"errors"
	"github.com/labstack/echo/v4"
)

const (
	headerTemplateKey = "Tmpl"
	BpsType = "bps"
	ResearchType = "rs"
	BookType = "b"
	SerieType = "s"
)

var (
	NoTemplateHeader = errors.New("No template header was given")
	MultipleTemplateHeader = errors.New("Multiple template were given in header")
)
func GetHeaderTemplate(c echo.Context) (string,error) {
	for key,values := range c.Request().Header {
		if key != headerTemplateKey {continue}
		if len(values) == 1 {
			return values[0],nil
		}
		return "",MultipleTemplateHeader
	}
	return "",NoTemplateHeader
}
