package types

type NamedDate struct {
	Day   string
	Month string
}

type RawDate struct {
	Day   int
	Month int
	Year  int
}

type TemplateValues struct {
	HolidayName   string
	HolidayList   []ParsedHolidays
	DaysLeft      int
	FormattedDate string
	NamedDate     NamedDate
	RawDate       RawDate
	FullDate      string
	IsToday       bool
	Length        int
	Adjacents     []ParsedHolidays
}

type MonthTemplateValues struct {
	Month        string
	Count        int
	HolidaysList []ParsedHolidays
	Adjacents    [][]ParsedHolidays
}
