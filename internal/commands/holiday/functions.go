package holidays

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"sort"
	"strconv"
	"time"

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

// Holiday represents a holiday
type Holiday struct {
	Date string `json:"fecha"`
	Type string `json:"tipo"`
	Name string `json:"nombre"`
}

// GetHolidays returns the holidays for the given year
func GetHolidays(year int) ([]Holiday, error) {
	cacheFile := fmt.Sprintf("/tmp/holidays_%d.json", year)

	// Check if the cache file exists and is not older than 24 hours
	if fileInfo, err := os.Stat(cacheFile); err == nil {
		if time.Since(fileInfo.ModTime()) < 24*time.Hour {
			if data, err := os.ReadFile(cacheFile); err == nil {
				var cachedHolidays []Holiday
				if err := json.Unmarshal(data, &cachedHolidays); err == nil {
					logrus.Infof("Loaded holidays for year %d from cache", year)
					return cachedHolidays, nil
				}
			}
		}
	}

	url := fmt.Sprintf(holidaysURL, strconv.Itoa(year))
	logrus.Infof("Getting holidays for year %d", year)
	logrus.Infof("URL: %s", url)
	resp, err := http.Get(url)

	logrus.Info("Getting next holiday")

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var holidays []Holiday
	if err := json.Unmarshal(body, &holidays); err != nil {
		return nil, err
	}

	// Save the response to the cache file
	if err := os.WriteFile(cacheFile, body, 0644); err != nil {
		logrus.Warnf("Failed to write cache file: %v", err)
	}

	return holidays, nil
}

// IsHoliday returns true if the given date is a holiday
func IsHoliday(date time.Time, holidays []Holiday) bool {
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
func NextHoliday(date time.Time, skipWeekends bool, skipToday bool) (*Holiday, bool) {
	holidays, err := GetHolidays(date.Year())
	if err != nil {
		return nil, false
	}

	for _, h := range holidays {
		holidayDate, err := time.Parse("2006-01-02", h.Date)
		if err != nil {
			continue
		}

		if holidayDate.After(date) || (!skipToday && isToday(holidayDate)) {
			if skipWeekends && (holidayDate.Weekday() == time.Saturday || holidayDate.Weekday() == time.Sunday) {
				continue
			}
			return &h, isToday(holidayDate)
		}
	}
	return nil, false
}

// Calculate how many days are left for the giving holiday
func DaysLeft(skipWeekends bool, skipToday bool) int {
	date := time.Now()
	nextHoliday, isToday := NextHoliday(date, skipWeekends, skipToday)
	if isToday {
		return 0
	}

	parsedDate, err := time.Parse("2006-01-02", nextHoliday.Date)
	if err != nil {
		return -1
	}

	return int(math.Ceil((parsedDate.Sub(date).Hours() + 3) / 24))
}

func GetAllHolidaysOfMonth(month Months, year int) ([]Holiday, [][]Holiday, error) {
	holidays, err := GetHolidays(year)
	if err != nil {
		return nil, nil, err
	}

	var holidaysOfMonth []Holiday
	for _, h := range holidays {
		holidayDate, err := time.Parse("2006-01-02", h.Date)
		if err != nil {
			continue
		}

		if holidayDate.Month() == time.Month(month) {
			holidaysOfMonth = append(holidaysOfMonth, h)
		}
	}

	adjacents := largeHolidays(holidaysOfMonth)

	return holidaysOfMonth, adjacents, nil
}

/*
*	Large holidays are two or more consecutive holidays
*	and/or adjacents to weekends
 */
func largeHolidays(holidays []Holiday) [][]Holiday {
	var adjacents [][]Holiday

	sort.Slice(holidays, func(i, j int) bool {
		return holidays[i].Date < holidays[j].Date
	})

	holidayMap := make(map[string]Holiday)
	for _, h := range holidays {
		holidayMap[h.Date] = h
	}

	visited := make(map[string]bool)
	for _, holiday := range holidays {
		if visited[holiday.Date] {
			continue
		}

		group := []Holiday{holiday}
		visited[holiday.Date] = true

		currentDate, _ := time.Parse("2006-01-02", holiday.Date)

		// Check previous days
		for d := currentDate.AddDate(0, 0, -1); ; d = d.AddDate(0, 0, -1) {
			weekday := d.Weekday()
			dateStr := d.Format("2006-01-02")

			if h, exists := holidayMap[dateStr]; exists && !visited[dateStr] {
				group = append([]Holiday{h}, group...)
				visited[dateStr] = true
			} else if weekday == time.Sunday || weekday == time.Saturday {
				group = append([]Holiday{{Date: dateStr, Type: "Weekend", Name: weekday.String()}}, group...)
				break
			} else if weekday == time.Monday {
				break
			}
		}

		// Check next days
		for d := currentDate.AddDate(0, 0, 1); ; d = d.AddDate(0, 0, 1) {
			weekday := d.Weekday()
			dateStr := d.Format("2006-01-02")

			if h, exists := holidayMap[dateStr]; exists && !visited[dateStr] {
				group = append(group, h)
				visited[dateStr] = true
			} else if weekday == time.Saturday || weekday == time.Sunday {
				group = append(group, Holiday{Date: dateStr, Type: "Weekend", Name: weekday.String()})
			} else if weekday == time.Friday {
				break
			}
		}

		// Discard groups with only a single holiday
		hasAdjacents := len(group) > 1
		if hasAdjacents {
			adjacents = append(adjacents, group)
		}
	}

	return adjacents
}
