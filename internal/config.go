package internal

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

const (
	configFileName     = "config.yaml"
	credentialFileName = "credentials.json"
)

type Config struct {
	Discord struct {
		GuildID string `yaml:"guild_id"`
	}
	Google struct {
		CalendarID  string `yaml:"calendar_id"`
		Credentials []byte `yaml:"credentials,omitempty"`
	}
	Secret struct {
		Token string `yaml:"token"`
	}
}

// NewConfig gets the bot config from the directory of the executable's path
func NewConfig() (*Config, error) {
	config := &Config{}

	var err error
	credentials := os.Getenv("GOOGLE_CREDENTIALS")
	if credentials == "" {
		config.Google.Credentials, err = os.ReadFile(credentialFileName)
		if err != nil {
			return nil, fmt.Errorf("cannot find google credentials: %v", err)
		}
	} else {
		config.Google.Credentials = []byte(credentials)
	}

	calendarID := os.Getenv("GOOGLE_CALENDAR_ID")
	if calendarID != "" {
		config.Google.CalendarID = calendarID
	}

	dig, token := os.Getenv("DISCORD_GUILD_ID"), os.Getenv("DISCORD_TOKEN")
	if dig != "" && token != "" {
		config.Discord.GuildID = dig
		config.Secret.Token = token
		return config, nil
	}

	path, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	file, err := os.Open(filepath.Join(path, configFileName))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	d := yaml.NewDecoder(file)
	if err := d.Decode(&config); err != nil {
		return nil, err
	}
	return config, nil
}
