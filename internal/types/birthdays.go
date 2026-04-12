package types

// Birthday structure to parse the YAML file containing friends' birthdays.
type Birthday struct {
	Name     string `yaml:"name"`
	Username string `yaml:"username"` // Discord username ping e.g <@123456>
	Date     string `yaml:"date"`     // Expected format: MM-DD
}
