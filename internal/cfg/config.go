package cfg

import (
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/BurntSushi/toml"
	"github.com/adrg/xdg"
)

type API struct {
	BaseURL string `toml:"base-url"`
	Key     string `toml:"key"`
}

type Config struct {
	API API `toml:"api"`
}

var configFile = "jrnl/config.toml"
var devConfigFile = "tmp/config.toml"

var (
	IsDev      = os.Getenv("DEV") == "true"
	EnableLogs = os.Getenv("JRNL_ENABLE_LOGS") == "true"
	_, b, _, _ = runtime.Caller(0)
	basepath   = filepath.Dir(b)
)

func GetEnv() string {
	if os.Getenv("DEV") == "true" {
		return "DEV"
	}
	if os.Getenv("TEST") == "true" {
		return "TEST"
	}

	return "PROD"
}

func GetProjectRoot() string {
	return filepath.Join(basepath, "../..")
}

func GetConfigPath() string {
	env := GetEnv()
	if env == "DEV" || env == "TEST" {
		return filepath.Join(GetProjectRoot(), "/tmp/config.toml")
	}

	configFilePath, err := xdg.ConfigFile(configFile)
	if err != nil {
		log.Fatal(err)
	}

	return configFilePath
}

var configData []byte

func GetConfig() Config {
	var f []byte
	var err error

	if configData != nil {
		f = configData
	} else {
		f, err = os.ReadFile(GetConfigPath())
		if err != nil {
			log.Fatal(err)
		}
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
