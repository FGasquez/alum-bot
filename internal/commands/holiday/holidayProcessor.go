package holidays

import (
	"encoding/json"
	"math"
	"sort"
	"time"

	"github.com/FGasquez/alum-bot/internal/types"
)

const dateLayout = "2006-01-02"

// HolidaysProcessor processes raw holiday data, applying filters and identifying relationships.
func HolidaysProcessor(jsonData []byte, skipPassed, adjacents, skipWeekends, skipToday bool) (types.ProcessedHolidays, error) {
	var rawHolidays []types.Holiday
	if err := json.Unmarshal(jsonData, &rawHolidays); err != nil {
		return types.ProcessedHolidays{}, err
	}

	now := time.Now()
	// Use the start of today for consistent date comparisons.
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	// Pre-allocate slice capacity to prevent reallocations, improving performance.
	parsedHolidays := make([]types.ParsedHolidays, 0, len(rawHolidays))

	// Iterate over raw holidays
	for _, h := range rawHolidays {
		date, err := time.Parse(dateLayout, h.Date)
		if err != nil {
			continue // Skip records with invalid dates.
		}

		isHolidayToday := date.Year() == today.Year() && date.YearDay() == today.YearDay()

		// Apply filters
		if skipToday && isHolidayToday {
			continue
		}
		if skipPassed && date.Before(today) && !isToday(date) {
			continue
		}
		if skipWeekends && isWeekend(date) {
			continue
		}

		// Add parsed holiday to slice
		parsedHolidays = append(parsedHolidays, types.ParsedHolidays{
			Date:              h.Date,
			Type:              h.Type,
			Name:              h.Name,
			FormattedDate:     date.Format("Monday, 2 January 2006"),
			NamedDate:         types.NamedDate{Day: date.Weekday().String(), Month: date.Month().String()},
			RawDate:           types.RawDate{Year: date.Year(), Month: int(date.Month()), Day: date.Day()},
			FullDate:          date.Format(time.RFC3339),
			IsToday:           isHolidayToday,
			DaysLeftToHoliday: int(math.Ceil(date.Sub(today).Hours() / 24)),
		})
	}

	sort.Slice(parsedHolidays, func(i, j int) bool {
		return parsedHolidays[i].Date < parsedHolidays[j].Date
	})

	// Group adjacent holidays
	if adjacents {
		parsedHolidays = groupAdjacentHolidays(parsedHolidays)
	}

	next, previous := findNextAndPrevious(parsedHolidays, today)

	return types.ProcessedHolidays{
		Next:     next,
		Previous: previous,
		All:      parsedHolidays,
	}, nil
}

func groupAdjacentHolidays(sortedHolidays []types.ParsedHolidays) []types.ParsedHolidays {
	if len(sortedHolidays) == 0 {
		return sortedHolidays
	}

	var allGroupedHolidays []types.ParsedHolidays
	var holidayGroups [][]*types.ParsedHolidays

	// Find groups of holidays that are directly adjacent (e.g., Mon, Tue).
	for i := 0; i < len(sortedHolidays); {
		group := []*types.ParsedHolidays{&sortedHolidays[i]}
		j := i
		for j+1 < len(sortedHolidays) {
			prevDate, _ := time.ParseInLocation(dateLayout, sortedHolidays[j].Date, time.Local)
			nextDate, _ := time.ParseInLocation(dateLayout, sortedHolidays[j+1].Date, time.Local)

			// Check if the next holiday is exactly one day after the current one.
			if !prevDate.AddDate(0, 0, 1).Equal(nextDate) {
				break // Not directly adjacent.
			}
			group = append(group, &sortedHolidays[j+1])
			j++
		}
		holidayGroups = append(holidayGroups, group)
		i = j + 1
	}

	// For each group, expand it with all adjacent weekends.
	for _, group := range holidayGroups {
		var finalGroup []types.ParsedHolidays

		// Get the start and end dates of the core holiday block.
		firstHolidayDate, _ := time.ParseInLocation(dateLayout, group[0].Date, time.Local)
		lastHolidayDate, _ := time.ParseInLocation(dateLayout, group[len(group)-1].Date, time.Local)

		// Find adjacent previous weekends.
		finalGroup = append(finalGroup, findPrecedingWeekends(firstHolidayDate)...)

		for i, holiday := range group {
			finalGroup = append(finalGroup, *holiday)
			// If there's another holiday in the group, check for a weekend gap.
			if i+1 < len(group) {
				currentDate, _ := time.ParseInLocation(dateLayout, holiday.Date, time.Local)
				nextDate, _ := time.ParseInLocation(dateLayout, group[i+1].Date, time.Local)
				finalGroup = append(finalGroup, createWeekendHolidaysBetween(currentDate, nextDate)...)
			}
		}

		finalGroup = append(finalGroup, findSucceedingWeekends(lastHolidayDate)...)

		// Link all items in the final group together.
		if len(finalGroup) > 1 {
			for i := range finalGroup {
				finalGroup[i].Adjacent = finalGroup
			}
		}
		allGroupedHolidays = append(allGroupedHolidays, finalGroup...)
	}

	return allGroupedHolidays
}

// findPrecedingWeekends finds all weekend days immediately before a given date.
func findPrecedingWeekends(startDate time.Time) []types.ParsedHolidays {
	var weekends []types.ParsedHolidays
	d := startDate.AddDate(0, 0, -1) // Start checking the day before.
	for isWeekend(d) {
		weekends = append([]types.ParsedHolidays{createWeekendHoliday(d)}, weekends...)
		d = d.AddDate(0, 0, -1)
	}
	return weekends
}

// findSucceedingWeekends finds all weekend days immediately after a given date.
func findSucceedingWeekends(endDate time.Time) []types.ParsedHolidays {
	var weekends []types.ParsedHolidays
	d := endDate.AddDate(0, 0, 1) // Start checking the day after.
	for isWeekend(d) {
		weekends = append(weekends, createWeekendHoliday(d))
		d = d.AddDate(0, 0, 1)
	}
	return weekends
}

// createWeekendHolidaysBetween finds weekend days that fall between two dates.
func createWeekendHolidaysBetween(start, end time.Time) []types.ParsedHolidays {
	var weekends []types.ParsedHolidays
	d := start.AddDate(0, 0, 1)
	for d.Before(end) {
		if isWeekend(d) {
			weekends = append(weekends, createWeekendHoliday(d))
		}
		d = d.AddDate(0, 0, 1)
	}
	return weekends
}

// findNextAndPrevious finds the next and previous holidays relative to today
func findNextAndPrevious(holidays []types.ParsedHolidays, today time.Time) (next, previous types.ParsedHolidays) {
	index := sort.Search(len(holidays), func(i int) bool {
		date, _ := time.ParseInLocation(dateLayout, holidays[i].Date, time.Local)
		return !date.Before(today)
	})

	// If a holiday is found at or after today
	if index < len(holidays) {
		date, _ := time.ParseInLocation(dateLayout, holidays[index].Date, time.Local)

		if date.Equal(today) {
			next = holidays[index]
			if index > 0 {
				previous = holidays[index-1]
			}
		} else {
			next = holidays[index]
			if index > 0 {
				previous = holidays[index-1]
			}
		}
	} else if len(holidays) > 0 {
		previous = holidays[len(holidays)-1]
	}

	return next, previous
}

func isWeekend(t time.Time) bool {
	weekday := t.Weekday()
	return weekday == time.Saturday || weekday == time.Sunday
}

// createWeekendHoliday is a factory function to build a weekend holiday object.
func createWeekendHoliday(day time.Time) types.ParsedHolidays {
	return types.ParsedHolidays{
		Date:          day.Format(dateLayout),
		Type:          types.Weekend,
		Name:          day.Weekday().String(),
		FormattedDate: day.Format("Monday, 2 January 2006"),
		NamedDate:     types.NamedDate{Day: day.Weekday().String(), Month: day.Month().String()},
		RawDate:       types.RawDate{Year: day.Year(), Month: int(day.Month()), Day: day.Day()},
		FullDate:      day.Format(time.RFC3339),
	}
}
