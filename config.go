package main

import (
	"encoding/json"
	"fmt"
	"os"
)

var configfile = "config.json"

type Config struct {
}

func loadConfig() (map[string]interface{}, error) {
	var data Config

	jsonData, err := os.ReadFile(configfile)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(jsonData, &data)
	if err != nil {
		return nil, err
	}

	fmt.Print(data) //TODO implement later
	return nil, nil
}
