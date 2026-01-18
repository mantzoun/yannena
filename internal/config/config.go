package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	SocketAddress string
	SaveFile      string
	LogLevel      int
}

func ReadConfig(filename string) (Config, error) {
	var c Config
	data, err := os.ReadFile(filename)
	if err != nil {
		return c, err
	}
	err = json.Unmarshal(data, &c)
	return c, err
}

func WriteConfig(c Config, filename string) error {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}
