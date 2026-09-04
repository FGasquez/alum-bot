package holidays

import (
	"fmt"
	"testing"
	"time"

	"github.com/FGasquez/alum-bot/internal/types"
)

func TestHolidaysProcessorOptimized(t *testing.T) {
	// Mock JSON Data with two holidays
	mockData := `[
		{"fecha":"2026-05-01","tipo":"inamovible","nombre":"Día del Trabajador"},
		{"fecha":"2026-05-25","tipo":"inamovible","nombre":"Día de la Revolución de Mayo"}
	]`

	processed, err := HolidaysProcessor([]byte(mockData), false, false, false, false)

	if err != nil {
		t.Fatalf("Failed to process holidays: %v", err)
	}

	if len(processed.All) != 2 {
		t.Fatalf("Expected 2 holidays, got %d", len(processed.All))
	}

	// Verify ParsedDate is correctly populated
	expectedDate1, _ := time.ParseInLocation("2006-01-02", "2026-05-01", time.Local)
	if !processed.All[0].ParsedDate.Equal(expectedDate1) {
		t.Errorf("Expected ParsedDate to be %v, got %v", expectedDate1, processed.All[0].ParsedDate)
	}

	expectedDate2, _ := time.ParseInLocation("2006-01-02", "2026-05-25", time.Local)
	if !processed.All[1].ParsedDate.Equal(expectedDate2) {
		t.Errorf("Expected ParsedDate to be %v, got %v", expectedDate2, processed.All[1].ParsedDate)
	}
}

func TestIsHolidayOptimized(t *testing.T) {
	mockHolidays := []types.Holiday{
		{Date: "2026-05-01", Type: "inamovible", Name: "Día del Trabajador"},
	}

	testDate, _ := time.Parse("2006-01-02", "2026-05-01")
	if !IsHoliday(testDate, mockHolidays) {
		t.Errorf("Expected %v to be recognized as a holiday", testDate)
	}

	testDateNotHoliday, _ := time.Parse("2006-01-02", "2026-05-02")
	if IsHoliday(testDateNotHoliday, mockHolidays) {
		t.Errorf("Expected %v to NOT be recognized as a holiday", testDateNotHoliday)
	}
}

func TestProcessedHolidaysStructure(t *testing.T) {
	now := time.Now()
	// Create dates relative to today to deterministically test Next and Previous
	yesterday := now.AddDate(0, 0, -1).Format("2006-01-02")
	today := now.Format("2006-01-02")
	tomorrow := now.AddDate(0, 0, 1).Format("2006-01-02")
	nextWeek := now.AddDate(0, 0, 7).Format("2006-01-02")

	mockData := fmt.Sprintf(`[
		{"fecha":"%s","tipo":"inamovible","nombre":"Yesterday Holiday"},
		{"fecha":"%s","tipo":"inamovible","nombre":"Today Holiday"},
		{"fecha":"%s","tipo":"inamovible","nombre":"Tomorrow Holiday"},
		{"fecha":"%s","tipo":"inamovible","nombre":"Next Week Holiday"}
	]`, yesterday, today, tomorrow, nextWeek)

	// Since we are not skipping passed holidays or today, we expect all 4 in All.
	processed, err := HolidaysProcessor([]byte(mockData), false, false, false, false)

	if err != nil {
		t.Fatalf("Failed to process holidays: %v", err)
	}

	if len(processed.All) != 4 {
		t.Fatalf("Expected 4 holidays, got %d", len(processed.All))
	}

	// Verify Previous is Yesterday
	if processed.Previous.Date != yesterday {
		t.Errorf("Expected Previous date to be %s, got %s", yesterday, processed.Previous.Date)
	}

	// Verify Next is Today (since today is a holiday and we didn't skip it, findNextAndPrevious sets Next to today)
	if processed.Next.Date != today {
		t.Errorf("Expected Next date to be %s, got %s", today, processed.Next.Date)
	}
}

func TestProcessedHolidaysSkipToday(t *testing.T) {
	now := time.Now()
	yesterday := now.AddDate(0, 0, -1).Format("2006-01-02")
	today := now.Format("2006-01-02")
	tomorrow := now.AddDate(0, 0, 1).Format("2006-01-02")

	mockData := fmt.Sprintf(`[
		{"fecha":"%s","tipo":"inamovible","nombre":"Yesterday Holiday"},
		{"fecha":"%s","tipo":"inamovible","nombre":"Today Holiday"},
		{"fecha":"%s","tipo":"inamovible","nombre":"Tomorrow Holiday"}
	]`, yesterday, today, tomorrow)

	// Skip today = true
	processed, err := HolidaysProcessor([]byte(mockData), false, false, false, true)

	if err != nil {
		t.Fatalf("Failed to process holidays: %v", err)
	}

	if len(processed.All) != 2 {
		t.Fatalf("Expected 2 holidays (yesterday and tomorrow), got %d", len(processed.All))
	}

	if processed.Previous.Date != yesterday {
		t.Errorf("Expected Previous date to be %s, got %s", yesterday, processed.Previous.Date)
	}

	if processed.Next.Date != tomorrow {
		t.Errorf("Expected Next date to be %s, got %s", tomorrow, processed.Next.Date)
	}
}

func TestProcessedHolidaysSkipWeekend(t *testing.T) {
	// A Wednesday, Saturday, and Sunday in May 2026. Data should simulate real dates.
	mockData := `[
		{"fecha":"2026-05-20","tipo":"inamovible","nombre":"Wednesday Holiday"},
		{"fecha":"2026-05-23","tipo":"inamovible","nombre":"Saturday Holiday"},
		{"fecha":"2026-05-24","tipo":"inamovible","nombre":"Sunday Holiday"}
	]`

	processed, err := HolidaysProcessor([]byte(mockData), false, false, true, false) // skipWeekends = true

	if err != nil {
		t.Fatalf("Failed to process holidays: %v", err)
	}

	if len(processed.All) != 1 {
		t.Fatalf("Expected 1 holiday (Wednesday), got %d", len(processed.All))
	}

	if processed.All[0].Date != "2026-05-20" {
		t.Errorf("Expected remaining holiday to be 2026-05-20, got %s", processed.All[0].Date)
	}
}

func TestProcessedHolidaysFlagsCombined(t *testing.T) {
	now := time.Now()
	// Let's force today to be a known weekday, but we can't easily fake time.Now() without a mock.
	// To avoid time.Now() flakiness, we'll just test that skipToday and skipWeekend logic
	// applies correctly using absolute dates where weekend behavior is known, and relative dates for today.

	today := now.Format("2006-01-02")
	// Let's use known weekend dates in 2026 for the weekend exclusion.
	// 2026-05-23 is a Saturday.
	// 2026-05-24 is a Sunday.

	mockData := fmt.Sprintf(`[
		{"fecha":"%s","tipo":"inamovible","nombre":"Today Holiday"},
		{"fecha":"2026-05-20","tipo":"inamovible","nombre":"Wednesday Holiday"},
		{"fecha":"2026-05-23","tipo":"inamovible","nombre":"Saturday Holiday"},
		{"fecha":"2026-05-24","tipo":"inamovible","nombre":"Sunday Holiday"}
	]`, today)

	// skipWeekends = true, skipToday = true
	processed, err := HolidaysProcessor([]byte(mockData), false, false, true, true)

	if err != nil {
		t.Fatalf("Failed to process holidays: %v", err)
	}

	// We expect Today to be skipped (if it equals today in Processor)
	// We expect 2026-05-23 and 2026-05-24 to be skipped because they are weekends.
	// The only remaining one should be 2026-05-20.

	var remainingHolidays []types.ParsedHolidays
	for _, h := range processed.All {
		if h.Date == today {
			// In case today = "2026-05-20", then the Wednesday holiday was the same as today Holiday, and both would be dropped.
			// But since we can't control 'today', we just enforce what should definitely NOT be there.
			t.Errorf("Today (%s) should have been skipped", today)
		}
		if h.Date == "2026-05-23" || h.Date == "2026-05-24" {
			t.Errorf("Weekend (%s) should have been skipped", h.Date)
		}
		remainingHolidays = append(remainingHolidays, h)
	}

	// We can't strictly assert len == 1 because today might be 2026-05-20,
	// but we can assert weekends and today were filtered.
}

func TestGroupAdjacentHolidays(t *testing.T) {
	// Let's create an artificial list of parsed holidays spanning adjacent days
	d1, _ := time.Parse("2006-01-02", "2026-05-01") // Friday
	d2, _ := time.Parse("2006-01-02", "2026-05-02") // Saturday
	d3, _ := time.Parse("2006-01-02", "2026-05-03") // Sunday
	d4, _ := time.Parse("2006-01-02", "2026-05-10") // Next week (not adjacent)

	holidays := []types.ParsedHolidays{
		{ParsedDate: d1, Date: "2026-05-01", Name: "H1"},
		{ParsedDate: d2, Date: "2026-05-02", Name: "H2"},
		{ParsedDate: d3, Date: "2026-05-03", Name: "H3"},
		{ParsedDate: d4, Date: "2026-05-10", Name: "H4"},
	}

	grouped := groupAdjacentHolidays(holidays)
	if len(grouped) != 5 {
		t.Fatalf("Expected 5 total flattened items (includes preceding Saturday 09), got %d", len(grouped))
	}
	if len(grouped[0].Adjacent) != 3 {
		t.Errorf("Expected 3 holidays in first group adjacency, got %d", len(grouped[0].Adjacent))
	}
	if len(grouped[4].Adjacent) != 2 {
		t.Errorf("Expected 2 adjacency for isolated sunday (which groups with preceding saturday), got %d", len(grouped[4].Adjacent))
	}
}

func TestCreateWeekendHolidaysBetween(t *testing.T) {
	// Friday 2026-05-01 to Monday 2026-05-04 -> should generate Saturday 02 and Sunday 03
	start, _ := time.Parse("2006-01-02", "2026-05-01")
	end, _ := time.Parse("2006-01-02", "2026-05-04")

	weekends := createWeekendHolidaysBetween(start, end)
	if len(weekends) != 2 {
		t.Fatalf("Expected 2 weekends, got %d", len(weekends))
	}
	if weekends[0].Date != "2026-05-02" || weekends[1].Date != "2026-05-03" {
		t.Errorf("Expected generated weekends for 02 and 03, got %s and %s", weekends[0].Date, weekends[1].Date)
	}

	// Between adjacent days (no weekends between)
	none := createWeekendHolidaysBetween(start, start.AddDate(0, 0, 1))
	if len(none) != 0 {
		t.Errorf("Expected 0 weekends, got %d", len(none))
	}
}

func TestFindPrecedingSucceedingWeekends(t *testing.T) {
	// A Monday holiday
	monday, _ := time.Parse("2006-01-02", "2026-05-04") // May 4th is Monday
	preceding := findPrecedingWeekends(monday)
	if len(preceding) != 2 {
		t.Fatalf("Expected 2 preceding weekend days, got %d", len(preceding))
	}
	if preceding[0].Date != "2026-05-02" || preceding[1].Date != "2026-05-03" {
		t.Errorf("Expected matching preceding dates")
	}

	succeeding := findSucceedingWeekends(monday)
	if len(succeeding) != 0 {
		t.Errorf("Expected 0 succeeding from a Monday, got %d", len(succeeding))
	}

	// A Friday holiday
	friday, _ := time.Parse("2006-01-02", "2026-05-01") // May 1st is Friday
	succeedingFri := findSucceedingWeekends(friday)
	if len(succeedingFri) != 2 {
		t.Fatalf("Expected 2 succeeding weekend days for Friday, got %d", len(succeedingFri))
	}
	if succeedingFri[0].Date != "2026-05-02" || succeedingFri[1].Date != "2026-05-03" {
		t.Errorf("Expected matching succeeding dates")
	}
}

func TestCreateWeekendHolidayDaysLeft(t *testing.T) {
	today, _ := time.Parse("2006-01-02", "2026-09-04")
	saturday, _ := time.Parse("2006-01-02", "2026-10-10")

	wh := createWeekendHoliday(saturday, today)
	if wh.DaysLeftToHoliday != 36 {
		t.Errorf("Expected 36 days left, got %d", wh.DaysLeftToHoliday)
	}
	if wh.IsToday {
		t.Errorf("Expected IsToday to be false")
	}

	todayWh := createWeekendHoliday(today, today)
	if todayWh.DaysLeftToHoliday != 0 {
		t.Errorf("Expected 0 days left for today, got %d", todayWh.DaysLeftToHoliday)
	}
	if !todayWh.IsToday {
		t.Errorf("Expected IsToday to be true for today")
	}
}

func TestGetNextLargeHoliday(t *testing.T) {
	satDate, _ := time.Parse("2006-01-02", "2026-10-10")
	sunDate, _ := time.Parse("2006-01-02", "2026-10-11")
	monDate, _ := time.Parse("2006-01-02", "2026-10-12")

	sat := types.ParsedHolidays{Date: "2026-10-10", ParsedDate: satDate, DaysLeftToHoliday: 36, Name: "Saturday", Type: types.Weekend}
	sun := types.ParsedHolidays{Date: "2026-10-11", ParsedDate: sunDate, DaysLeftToHoliday: 37, Name: "Sunday", Type: types.Weekend}
	mon := types.ParsedHolidays{Date: "2026-10-12", ParsedDate: monDate, DaysLeftToHoliday: 38, Name: "Día del Respeto a la Diversidad Cultural"}

	group := []types.ParsedHolidays{sat, sun, mon}
	sat.Adjacent = group
	sun.Adjacent = group
	mon.Adjacent = group

	processed := types.ProcessedHolidays{
		All: []types.ParsedHolidays{sat, sun, mon},
	}

	nextLarge := GetNextLargeHoliday(processed)
	if nextLarge == nil {
		t.Fatalf("Expected large holiday, got nil")
	}
	if nextLarge.Date != "2026-10-10" {
		t.Errorf("Expected date to be 2026-10-10, got %s", nextLarge.Date)
	}
	if nextLarge.DaysLeftToHoliday != 36 {
		t.Errorf("Expected DaysLeftToHoliday to be 36, got %d", nextLarge.DaysLeftToHoliday)
	}
}

func TestGetNextLargeHolidaySkipsPast(t *testing.T) {
	pastSat, _ := time.Parse("2006-01-02", "2026-08-15")
	futSat, _ := time.Parse("2006-01-02", "2026-10-10")

	past := types.ParsedHolidays{
		Date:              "2026-08-15",
		ParsedDate:        pastSat,
		DaysLeftToHoliday: -20,
		Adjacent:          []types.ParsedHolidays{{Date: "2026-08-15"}, {Date: "2026-08-16"}},
	}
	fut := types.ParsedHolidays{
		Date:              "2026-10-10",
		ParsedDate:        futSat,
		DaysLeftToHoliday: 36,
		Adjacent:          []types.ParsedHolidays{{Date: "2026-10-10"}, {Date: "2026-10-11"}},
	}

	processed := types.ProcessedHolidays{
		All: []types.ParsedHolidays{past, fut},
	}

	nextLarge := GetNextLargeHoliday(processed)
	if nextLarge == nil {
		t.Fatalf("Expected upcoming large holiday, got nil")
	}
	if nextLarge.Date != "2026-10-10" {
		t.Errorf("Expected date to be 2026-10-10, got %s", nextLarge.Date)
	}
}

func TestGroupAdjacentHolidaysBridgeWeekend(t *testing.T) {
	// Friday May 1st and Monday May 4th should be connected across the weekend
	d1, _ := time.Parse("2006-01-02", "2026-05-01") // Friday
	d2, _ := time.Parse("2006-01-02", "2026-05-04") // Monday

	holidays := []types.ParsedHolidays{
		{ParsedDate: d1, Date: "2026-05-01", Name: "Workers Day"},
		{ParsedDate: d2, Date: "2026-05-04", Name: "Monday Holiday"},
	}

	grouped := groupAdjacentHolidays(holidays)
	// Friday + Saturday + Sunday + Monday = 4 days
	if len(grouped) != 4 {
		t.Fatalf("Expected 4 items in grouped bridge holiday, got %d", len(grouped))
	}
	for _, h := range grouped {
		if len(h.Adjacent) != 4 {
			t.Errorf("Expected 4 adjacent items in group, got %d for %s", len(h.Adjacent), h.Date)
		}
	}
}
