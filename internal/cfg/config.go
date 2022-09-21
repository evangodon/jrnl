package cfg

import (
	"log"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/adrg/xdg"
)

type API struct {
	BaseURL string `toml:"base-url"`
}

type Config struct {
	API API `toml:"api"`
}

var configFile = "jrnl/config.toml"

func GetConfigPath() string {
	configFilePath, err := xdg.ConfigFile(configFile)
	if err != nil {
		log.Fatal(err)
	}

	return configFilePath
}

func GetConfig() Config {
	f, err := os.ReadFile(GetConfigPath())

	if err != nil {
		log.Fatal(err)
	}

	var config Config

	if _, err := toml.Decode(string(f[:]), &config); err != nil {
		log.Fatal(err)
	}

	return config
}

func CreateConfigFile() error {
	configFilePath, err := xdg.ConfigFile(configFile)
	if err != nil {
		return err
	}
	f, err := os.OpenFile(configFilePath, os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		return err
	}

	configTemplate := Config{}

	if err := toml.NewEncoder(f).Encode(configTemplate); err != nil {
		return err
	}

	if err := f.Close(); err != nil {
		return err
	}

	return nil
}
