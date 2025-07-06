package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

func LoadConfig() (map[string]interface{}, error) {
	configFile, err := os.ReadFile("evoke.yaml")
	if err != nil {
		// If the file doesn't exist, return an empty config
		if os.IsNotExist(err) {
			return make(map[string]interface{}), nil
		}
		return nil, err
	}

	var config map[string]interface{}
	err = yaml.Unmarshal(configFile, &config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
