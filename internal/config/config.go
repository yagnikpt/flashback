package config

import (
	"os"

	"github.com/BurntSushi/toml"
)

type Config struct {
	ShowHelp bool   `toml:"show_help"`
	APIKey   string `toml:"api_key"`
}

func LoadConfig(filePath string) (Config, error) {
	var cfg Config

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		cfg = Config{
			ShowHelp: true,
			APIKey:   "",
		}
		SaveConfig(filePath, cfg)
		return cfg, nil
	}

	_, err := toml.DecodeFile(filePath, &cfg)
	if err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func SaveConfig(filePath string, cfg Config) error {
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()
	encoder := toml.NewEncoder(f)
	return encoder.Encode(cfg)
}
