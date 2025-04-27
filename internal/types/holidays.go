package types

// raw holiday
type Holiday struct {
	Date string `json:"fecha"`
	Type string `json:"tipo"`
	Name string `json:"nombre"`
}

type ParsedHolidays struct {
	Date              string
	Type              string
	Name              string
	FormattedDate     string
	NamedDate         NamedDate
	RawDate           RawDate
	FullDate          string
	Count             int
	Adjacent          []ParsedHolidays
	IsToday           bool
	DaysLeftToHoliday int
}

type ProcessedHolidays struct {
	Next     ParsedHolidays
	Previous ParsedHolidays
	All      []ParsedHolidays
}
