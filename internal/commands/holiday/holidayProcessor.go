package holidays

import (
	"encoding/json"
	"sort"
	"time"

	"github.com/FGasquez/alum-bot/internal/types"
)

func HolidaysProcessor(jsonData []byte, skipPassed bool) (types.ProcessedHolidays, error) {
	var rawHolidays []types.Holiday
	err := json.Unmarshal(jsonData, &rawHolidays)
	if err != nil {
		return types.ProcessedHolidays{}, err
	}

	var parsed []types.ParsedHolidays
	now := time.Now()
	layout := "2006-01-02"

	for _, h := range rawHolidays {
		date, err := time.Parse(layout, h.Date)
		if err != nil {
			continue // or return error
		}

		// SKIP PASSED HOLIDAYS if needed
		if skipPassed && date.Before(now.Truncate(24*time.Hour)) {
			continue
		}

		daysLeft := int(date.Sub(now).Hours() / 24)

		parsedHoliday := types.ParsedHolidays{
			Date:          h.Date,
			Type:          h.Type,
			Name:          h.Name,
			FormattedDate: date.Format("Monday, 2 January 2006"),
			NamedDate: types.NamedDate{
				Day:   date.Weekday().String(),
				Month: date.Month().String(),
			},
			RawDate: types.RawDate{
				Year:  date.Year(),
				Month: int(date.Month()),
				Day:   date.Day(),
			},
			FullDate:          date.Format(time.RFC3339),
			Count:             0,
			Adjacent:          []types.ParsedHolidays{}, // Filled later
			IsToday:           date.Year() == now.Year() && date.YearDay() == now.YearDay(),
			DaysLeftToHoliday: daysLeft,
		}

		parsed = append(parsed, parsedHoliday)
	}

	// Sort parsed holidays
	sort.Slice(parsed, func(i, j int) bool {
		return parsed[i].Date < parsed[j].Date
	})

	parsed = detectAdjacents(parsed)

	// Find next and previous holidays
	var next types.ParsedHolidays
	var previous types.ParsedHolidays

	for _, p := range parsed {
		date, _ := time.Parse(layout, p.Date)
		if date.After(now) && next.Date == "" {
			next = p
		}
		if date.Before(now) {
			previous = p
		}
	}

	// 🛠️ Filter only real holidays (exclude weekends)
	realHolidays := []types.ParsedHolidays{}
	for _, h := range parsed {
		if h.Type != "weekend" {
			realHolidays = append(realHolidays, h)
		}
	}

	return types.ProcessedHolidays{
		Next:     next,
		Previous: previous,
		All:      realHolidays,
	}, nil
}

func detectAdjacents(holidays []types.ParsedHolidays) []types.ParsedHolidays {
	dateLayout := "2006-01-02"
	holidayMap := make(map[string]types.ParsedHolidays)

	for _, h := range holidays {
		holidayMap[h.Date] = h
	}

	var enhanced []types.ParsedHolidays
	seen := make(map[string]bool)

	for _, holiday := range holidays {
		if seen[holiday.Date] {
			continue
		}

		currentGroup := []types.ParsedHolidays{}
		date, _ := time.Parse(dateLayout, holiday.Date)

		// Look backward
		d := date.AddDate(0, 0, -1)
		for {
			dStr := d.Format(dateLayout)
			if h, ok := holidayMap[dStr]; ok {
				currentGroup = append(currentGroup, h)
				seen[dStr] = true
			} else if isWeekend(d) {
				currentGroup = append(currentGroup, createWeekendHoliday(d))
			} else {
				break // Stop if it's a weekday (Monday-Friday)
			}
			d = d.AddDate(0, 0, -1) // go back one more day
		}

		// Add current holiday
		currentGroup = append(currentGroup, holiday)
		seen[holiday.Date] = true

		// Look forward
		d = date.AddDate(0, 0, 1)
		for {
			dStr := d.Format(dateLayout)
			if h, ok := holidayMap[dStr]; ok {
				currentGroup = append(currentGroup, h)
				seen[dStr] = true
			} else if isWeekend(d) {
				currentGroup = append(currentGroup, createWeekendHoliday(d))
			} else {
				break // Stop if it's a weekday (Monday-Friday)
			}
			d = d.AddDate(0, 0, 1) // go forward one more day
		}

		// Sort group
		sort.Slice(currentGroup, func(i, j int) bool {
			return currentGroup[i].Date < currentGroup[j].Date
		})

		if len(currentGroup) > 1 && !isWeekend(date) {
			for idx := range currentGroup {
				currentGroup[idx].Adjacent = currentGroup
			}
		}

		enhanced = append(enhanced, currentGroup...)
	}

	// Remove duplicates
	final := []types.ParsedHolidays{}
	unique := make(map[string]bool)

	for _, h := range enhanced {
		if !unique[h.Date] {
			final = append(final, h)
			unique[h.Date] = true
		}
	}

	return final
}

func isWeekend(t time.Time) bool {
	return t.Weekday() == time.Saturday || t.Weekday() == time.Sunday
}

func createWeekendHoliday(day time.Time) types.ParsedHolidays {
	return types.ParsedHolidays{
		Date:          day.Format("2006-01-02"),
		Type:          "weekend",
		Name:          day.Weekday().String(), // Saturday or Sunday
		FormattedDate: day.Format("Monday, 2 January 2006"),
		NamedDate: types.NamedDate{
			Day:   day.Weekday().String(),
			Month: day.Month().String(),
		},
		RawDate: types.RawDate{
			Year:  day.Year(),
			Month: int(day.Month()),
			Day:   day.Day(),
		},
		FullDate:          day.Format(time.RFC3339),
		IsToday:           false,
		DaysLeftToHoliday: int(day.Sub(time.Now()).Hours() / 24),
	}
}
