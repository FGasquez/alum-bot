package holidays

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/FGasquez/alum-bot/internal/types"
	"github.com/sirupsen/logrus"
)

const (
	holidaysURL = "https://api.argentinadatos.com/v1/feriados/%s"
)

type Months int

const (
	January   Months = 1
	February  Months = 2
	March     Months = 3
	April     Months = 4
	May       Months = 5
	June      Months = 6
	July      Months = 7
	August    Months = 8
	September Months = 9
	October   Months = 10
	November  Months = 11
	December  Months = 12
)

// GetHolidays returns the holidays for the given year
func GetHolidays(year int, skipPassed bool) (types.ProcessedHolidays, error) {
	cacheFile := fmt.Sprintf("/tmp/holidays_%d.json", year)
	data := []byte{}
	// Check if the cache file exists and is not older than 24 hours
	if fileInfo, err := os.Stat(cacheFile); err == nil {
		if time.Since(fileInfo.ModTime()) < 24*time.Hour {
			if data, err = os.ReadFile(cacheFile); err != nil {
				logrus.Warnf("Failed to read cache file: %v", err)
			}
		}
	}

	if len(data) < 1 {
		logrus.Info("Cache file not found or is older than 24 hours")
		url := fmt.Sprintf(holidaysURL, strconv.Itoa(year))
		logrus.Infof("Getting holidays for year %d", year)
		logrus.Infof("URL: %s", url)
		resp, err := http.Get(url)

		logrus.Info("Getting next holiday")

		if err != nil {
			return types.ProcessedHolidays{}, err
		}
		defer resp.Body.Close()

		data, err = io.ReadAll(resp.Body)
		if err != nil {
			return types.ProcessedHolidays{}, err
		}

		// Save the response to the cache file
		if err := os.WriteFile(cacheFile, data, 0644); err != nil {
			logrus.Warnf("Failed to write cache file: %v", err)
		}

	}

	processedHolidays, err := HolidaysProcessor(data, skipPassed)
	if err != nil {
		return types.ProcessedHolidays{}, err
	}

	return processedHolidays, nil
}

// IsHoliday returns true if the given date is a holiday
func IsHoliday(date time.Time, holidays []types.Holiday) bool {
	for _, h := range holidays {
		holidayDate, err := time.Parse("2006-01-02", h.Date)
		if err != nil {
			continue
		}
		if date.Equal(holidayDate) {
			return true
		}
	}
	return false
}

func isToday(parsedDate time.Time) bool {
	y1, m1, d1 := time.Now().Date()
	y2, m2, d2 := parsedDate.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}

// NextHoliday returns the next holiday
func NextHoliday(date time.Time, skipWeekends bool, skipToday bool) (*types.ParsedHolidays, bool) {
	holidays, err := GetHolidays(date.Year(), false)
	if err != nil {
		return nil, false
	}

	return &holidays.Next, holidays.Next.IsToday
}

// Calculate how many days are left for the giving holiday
func DaysLeft(skipWeekends bool, skipToday bool) (int, types.ParsedHolidays, bool) {
	holidays, err := GetHolidays(time.Now().Year(), true)
	if err != nil {
		return 0, types.ParsedHolidays{}, false
	}

	return holidays.Next.DaysLeftToHoliday, holidays.Next, holidays.Next.IsToday
}

func GetAllHolidaysOfMonth(month Months, year int) ([]types.ParsedHolidays, error) {
	holidays, err := GetHolidays(year, false)
	if err != nil {
		return nil, err
	}

	var holidaysOfMonth []types.ParsedHolidays
	for _, holiday := range holidays.All {

		if Months(holiday.RawDate.Month) == month {
			holidaysOfMonth = append(holidaysOfMonth, holiday)
		}
	}

	return holidaysOfMonth, nil
}
