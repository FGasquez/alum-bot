package helpers

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	holydaysURL = "https://api.argentinadatos.com/v1/feriados/%s"
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

// Holyday represents a holyday
type Holyday struct {
	Date string `json:"fecha"`
	Type string `json:"tipo"`
	Name string `json:"nombre"`
}

// GetHolydays returns the holydays for the given year
func GetHolydays(year int) ([]Holyday, error) {
	cacheFile := fmt.Sprintf("/tmp/holydays_%d.json", year)

	// Check if the cache file exists and is not older than 24 hours
	if fileInfo, err := os.Stat(cacheFile); err == nil {
		if time.Since(fileInfo.ModTime()) < 24*time.Hour {
			if data, err := os.ReadFile(cacheFile); err == nil {
				var cachedHolydays []Holyday
				if err := json.Unmarshal(data, &cachedHolydays); err == nil {
					logrus.Infof("Loaded holydays for year %d from cache", year)
					return cachedHolydays, nil
				}
			}
		}
	}

	url := fmt.Sprintf(holydaysURL, strconv.Itoa(year))
	logrus.Infof("Getting holydays for year %d", year)
	logrus.Infof("URL: %s", url)
	resp, err := http.Get(url)

	logrus.Info("Getting next holyday")

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var holydays []Holyday
	if err := json.Unmarshal(body, &holydays); err != nil {
		return nil, err
	}

	// Save the response to the cache file
	if err := os.WriteFile(cacheFile, body, 0644); err != nil {
		logrus.Warnf("Failed to write cache file: %v", err)
	}

	return holydays, nil
}

// IsHolyday returns true if the given date is a holyday
func IsHolyday(date time.Time, holydays []Holyday) bool {
	for _, h := range holydays {
		holydayDate, err := time.Parse("2006-01-02", h.Date)
		if err != nil {
			continue
		}
		if date.Equal(holydayDate) {
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

// NextHolyday returns the next holyday
func NextHolyday(date time.Time, skipWeekends bool, skipToday bool) (*Holyday, bool) {
	holydays, err := GetHolydays(date.Year())
	if err != nil {
		return nil, false
	}

	for _, h := range holydays {
		holydayDate, err := time.Parse("2006-01-02", h.Date)
		if err != nil {
			continue
		}

		if holydayDate.After(date) || (!skipToday && isToday(holydayDate)) {
			if skipWeekends && (holydayDate.Weekday() == time.Saturday || holydayDate.Weekday() == time.Sunday) {
				continue
			}
			return &h, isToday(holydayDate)
		}
	}
	return nil, false
}

// Calculate how many days are left for the giving holyday
func DaysLeft(skipWeekends bool, skipToday bool) int {
	date := time.Now()
	nextHoliday, isToday := NextHolyday(date, skipWeekends, skipToday)
	if isToday {
		return 0
	}

	parsedDate, err := time.Parse("2006-01-02", nextHoliday.Date)
	if err != nil {
		return -1
	}

	return int(math.Ceil((parsedDate.Sub(date).Hours() + 3) / 24))
}

func GetAllHolidaysOfMonth(month Months, year int) ([]Holyday, error) {
	holydays, err := GetHolydays(year)
	if err != nil {
		return nil, err
	}

	var holydaysOfMonth []Holyday
	for _, h := range holydays {
		holydayDate, err := time.Parse("2006-01-02", h.Date)
		if err != nil {
			continue
		}

		if holydayDate.Month() == time.Month(month) {
			holydaysOfMonth = append(holydaysOfMonth, h)
		}
	}

	return holydaysOfMonth, nil
}
