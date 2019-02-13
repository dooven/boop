package config

import (
	"encoding/json"
	"fmt"
	"github.com/dooven/boop/utils"
	"io"
	"io/ioutil"
	"os"
)

type RegionOption struct {
	Name   string
	Region string
}

type Config struct {
	Regions []RegionOption
	Users   []string
}

var configFileName = ".boop.json"

func GetOrWriteDefaults() (*Config, error) {
	configPath := fmt.Sprintf("%s/%s", utils.UserHomeDir(), configFileName)

	_, err := os.Lstat(configPath)

	if os.IsNotExist(err) {
		jsonFile, err := os.Create(configPath)
		if err != nil {
			return nil, err
		}

		jsonWriter := io.Writer(jsonFile)
		encoder := json.NewEncoder(jsonWriter)

		err = encoder.Encode(CONFIG_DEFAULTS)

		if err != nil {
			return nil, err
		}

		return &CONFIG_DEFAULTS, nil
	}

	jsonFile, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}

	defer jsonFile.Close()

	jsonData, err := ioutil.ReadAll(jsonFile)

	if err != nil {
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(jsonData, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
