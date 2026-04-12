package birthday

import (
	"fmt"
	"io"
	"math"
	"os"
	"sort"
	"time"

	"github.com/FGasquez/alum-bot/internal/config"
	"github.com/FGasquez/alum-bot/internal/types"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

// GetBirthdays loads and parses the birthdays YAML file.
func GetBirthdays() ([]types.Birthday, error) {
	path := config.GetBirthdaysPath()
	if path == "" {
		return nil, fmt.Errorf("birthdays-file configuration is empty")
	}

	yamlFile, err := os.Open(path)
	if err != nil {
		logrus.Errorf("Error opening birthdays file: %s", err)
		return nil, err
	}
	defer yamlFile.Close()

	byteValue, err := io.ReadAll(yamlFile)
	if err != nil {
		return nil, err
	}

	var rawBirthdays []types.Birthday
	if err := yaml.Unmarshal(byteValue, &rawBirthdays); err != nil {
		logrus.Errorf("Error unmarshalling birthdays file: %s", err)
		return nil, err
	}

	return rawBirthdays, nil
}

type ParsedBirthday struct {
	types.Birthday
	NextOccurrence time.Time
	DaysLeft       int
}

// GetNextBirthday determines which birthday is coming up next relative to today.
func GetNextBirthday(birthdays []types.Birthday) (*ParsedBirthday, error) {
	if len(birthdays) == 0 {
		return nil, fmt.Errorf("no birthdays provided")
	}

	now := time.Now().In(config.GetTimezone())
	// Normalize today to start-of-day for consistent duration calculation
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	var parsed []ParsedBirthday

	for _, b := range birthdays {
		// Expect format MM-DD, we prefix it with a leap-year to allow clean parsing.
		t, err := time.Parse("01-02", b.Date)
		if err != nil {
			logrus.Warnf("Invalid date format for %s: %s. Expected MM-DD", b.Name, b.Date)
			continue
		}

		occurrence := time.Date(today.Year(), t.Month(), t.Day(), 0, 0, 0, 0, today.Location())

		// If the birthday has already passed this year, the next one is next year.
		// If it's today, occurrence.Before(today) is false, so it triggers today.
		if occurrence.Before(today) {
			occurrence = occurrence.AddDate(1, 0, 0)
		}

		daysLeft := int(math.Round(occurrence.Sub(today).Hours() / 24))

		parsed = append(parsed, ParsedBirthday{
			Birthday:       b,
			NextOccurrence: occurrence,
			DaysLeft:       daysLeft,
		})
	}

	if len(parsed) == 0 {
		return nil, fmt.Errorf("no valid birthdays could be parsed")
	}

	// Sort by whoever's birthday is coming first.
	sort.Slice(parsed, func(i, j int) bool {
		return parsed[i].NextOccurrence.Before(parsed[j].NextOccurrence)
	})

	return &parsed[0], nil
}

// CalculateDaysLeft calculates the remaining days for a given birthday object.
func CalculateDaysLeft(targetBirthday *types.Birthday) (*ParsedBirthday, error) {
	now := time.Now().In(config.GetTimezone())
	// Normalize today to start-of-day for consistent duration calculation
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	t, err := time.Parse("01-02", targetBirthday.Date)
	if err != nil {
		logrus.Warnf("Invalid date format for %s: %s. Expected MM-DD", targetBirthday.Name, targetBirthday.Date)
		return nil, fmt.Errorf("invalid date format for user")
	}

	occurrence := time.Date(today.Year(), t.Month(), t.Day(), 0, 0, 0, 0, today.Location())

	if occurrence.Before(today) {
		occurrence = occurrence.AddDate(1, 0, 0)
	}

	daysLeft := int(math.Round(occurrence.Sub(today).Hours() / 24))

	return &ParsedBirthday{
		Birthday:       *targetBirthday,
		NextOccurrence: occurrence,
		DaysLeft:       daysLeft,
	}, nil
}

// GetDaysLeft determines how many days are left until a specific user's birthday based on username.
func GetDaysLeft(username string, birthdays []types.Birthday) (*ParsedBirthday, error) {
	if len(birthdays) == 0 {
		return nil, fmt.Errorf("no birthdays provided")
	}

	var targetBirthday *types.Birthday
	for _, b := range birthdays {
		// Comparing against the provided username
		if b.Username == username {
			bCopy := b
			targetBirthday = &bCopy
			break
		}
	}

	if targetBirthday == nil {
		return nil, fmt.Errorf("user %s not found in birthdays list", username)
	}

	return CalculateDaysLeft(targetBirthday)
}
