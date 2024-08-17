package util

import (
	"errors"
	"strconv"

	"github.com/labstack/echo/v4"
)

const (
	//Clé du paramètre HTTP signifiant la nature de la réponse HTML voulue (élément ou page entière)
	headerTemplateKey = "Tmpl"
	//Clé du paramètre HTTP indiquant si le mode série est activé
	headerSerieModeKey = "Smode"
	//Valeur du paramètre HTTP headerTemplateKey pour un élement
	MainContentType = "main"
	//Valeur du paramètre HTTP headerTemplateKey pour un ensemble de prévisuels
	PreviewSetContentType = "previewSet"
)

var (
	NoTemplateHeader = errors.New("no template header was given")
)

// GetHeaderTemplate lit le header HTTP et renvoit la valeur du paramètre headerTemplateKey
func GetHeaderTemplate(c echo.Context) (string, error) {
	header := c.Request().Header.Get(headerTemplateKey)
	if header == "" {
		return "", NoTemplateHeader
	}
	return header, nil
}

// IsSerieMode renvoit true si le paramètre HTTP headerSerieModeKey est présent et vrai
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
