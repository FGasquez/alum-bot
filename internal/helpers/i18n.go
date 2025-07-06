package helpers

import (
	"fmt"
	"strconv"
	"time"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var spanishWeekdays = map[time.Weekday]string{
	time.Monday:    "Lunes",
	time.Tuesday:   "Martes",
	time.Wednesday: "Miércoles",
	time.Thursday:  "Jueves",
	time.Friday:    "Viernes",
	time.Saturday:  "Sábado",
	time.Sunday:    "Domingo",
}

var spanishMonths = map[time.Month]string{
	time.January:   "Enero",
	time.February:  "Febrero",
	time.March:     "Marzo",
	time.April:     "Abril",
	time.May:       "Mayo",
	time.June:      "Junio",
	time.July:      "Julio",
	time.August:    "Agosto",
	time.September: "Septiembre",
	time.October:   "Octubre",
	time.November:  "Noviembre",
	time.December:  "Diciembre",
}

func FormatDateToSpanishUnparsed(date string) (string, string, string, string) {
	parsedDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		return "", "", "", ""
	}
	return FormatDateToSpanish(parsedDate)

}
func FormatDateToSpanish(parsedDate time.Time) (string, string, string, string) {
	// Days and months in Spanish
	days := []string{"domingo", "lunes", "martes", "miércoles", "jueves", "viernes", "sábado"}
	months := []string{"enero", "febrero", "marzo", "abril", "mayo", "junio", "julio", "agosto", "septiembre", "octubre", "noviembre", "diciembre"}

	day := days[parsedDate.Weekday()]
	month := months[parsedDate.Month()-1]
	year := parsedDate.Format("2006")
	caser := cases.Title(language.Spanish)
	dayNumber := parsedDate.Day()
	dayFormatted := caser.String(day) + ", " + formatDayNumber(dayNumber) + " de " + month

	return dayFormatted, day, month, year
}

func formatDayNumber(day int) string {
	if day == 1 {
		return "1ro"
	}
	return strconv.Itoa(day)
}

func MonthsToSpanish(month int64) string {
	months := []string{"enero", "febrero", "marzo", "abril", "mayo", "junio", "julio", "agosto", "septiembre", "octubre", "noviembre", "diciembre"}
	return months[month-1]
}

func FormatDate(dateStr string) string {
	layout := "2006-01-02"
	t, err := time.Parse(layout, dateStr)
	if err != nil {
		return dateStr // fallback
	}

	day := t.Day()
	weekday := spanishWeekdays[t.Weekday()]
	month := spanishMonths[t.Month()]

	return fmt.Sprintf("%s %d de %s", weekday, day, month)
}
