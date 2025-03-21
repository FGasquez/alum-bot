package helpers

import (
	"strconv"
	"time"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func FormatDateToSpanish(parsedDate time.Time) (string, string, string) {
	// Days and months in Spanish
	days := []string{"domingo", "lunes", "martes", "miércoles", "jueves", "viernes", "sábado"}
	months := []string{"enero", "febrero", "marzo", "abril", "mayo", "junio", "julio", "agosto", "septiembre", "octubre", "noviembre", "diciembre"}

	day := days[parsedDate.Weekday()]
	month := months[parsedDate.Month()-1]
	year := parsedDate.Format("2006")
	caser := cases.Title(language.Spanish)
	dayNumber := parsedDate.Day()
	dayFormatted := caser.String(day) + ", " + formatDayNumber(dayNumber) + " de " + month

	return dayFormatted, month, year
}

func formatDayNumber(day int) string {
	if day == 1 {
		return "1ro"
	}
	return strconv.Itoa(day)
}
