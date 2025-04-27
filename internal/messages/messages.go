package messages

import (
	"bytes"
	"io"
	"os"
	"text/template"

	"github.com/FGasquez/alum-bot/internal/config"
	"github.com/FGasquez/alum-bot/internal/helpers"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type MessageKeysStruct struct {
	NextHoliday      string
	DaysLeft         string
	HolidaysOfMonth  string
	NextLargeHoliday string
}

var MessageKeys = MessageKeysStruct{
	NextHoliday:      "nextHoliday",
	DaysLeft:         "daysLeft",
	HolidaysOfMonth:  "holidaysOfMonth",
	NextLargeHoliday: "nextLargeHoliday",
}

var Messages map[string]string

var defaultMessages = map[string]string{
	MessageKeys.NextHoliday:      "The next holiday is **{{ .HolidayName }}**",
	MessageKeys.DaysLeft:         "There are **{{ .Days }}** days left for **{{ .HolidayName }}**",
	MessageKeys.HolidaysOfMonth:  "There are **{{ .Count }}** holidays in **{{ .Month }}**: {{ range .HolidaysList }}**{{ .Name }}**, {{ end }}",
	MessageKeys.NextLargeHoliday: "The next large holiday is **{{ .HolidayName }}**",
}

func ParseMessagesFromFile(filename string) map[string]string {
	var messages map[string]string
	yamlFile, err := os.Open(filename)
	if err != nil {
		logrus.Infof("Error opening file %s: %s", filename, err)
		return defaultMessages
	}
	defer yamlFile.Close()

	byteValue, _ := io.ReadAll(yamlFile)
	err = yaml.Unmarshal(byteValue, &messages)
	if err != nil {
		return defaultMessages
	}

	return messages
}

func GetMessage(key string) string {

	logrus.Infof("Loading message %s", key)
	var fileMessages = ParseMessagesFromFile(config.GetMessagesPath())
	for k, v := range fileMessages {
		logrus.Infof("%s: %s", k, v)
	}

	if fileMessages[key] != "" {
		return fileMessages[key]
	}

	return defaultMessages[key]
}

func TemplateMessage(message string, data interface{}) string {
	logrus.Infof("Message: %s", message)

	funcMap := template.FuncMap{
		"sub":        func(a, b int) int { return a - b },
		"formatDate": helpers.FormatDate,
	}

	tmpl, err := template.New("message").Funcs(funcMap).Parse(message)
	if err != nil {
		logrus.Errorf("Failed to parse message: %v", err)
		return ""
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		logrus.Errorf("Failed to execute template: %v", err)
		return ""
	}

	return buf.String()
}
