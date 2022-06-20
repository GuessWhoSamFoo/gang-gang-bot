package internal

import (
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

type Config struct {
	Discord struct {
		GuildID string `yaml:"guild_id"`
	}
	Secret struct {
		Token string `yaml:"token"`
	}
}

// NewConfig gets the bot config from the directory of the executable's path
func NewConfig() (*Config, error) {
	config := &Config{}

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
	file, err := os.Open(filepath.Join(path, "config.yaml"))
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
