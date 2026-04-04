package main

import (
	"context"
	"encoding/json"
	"os"
)

var configfile = "config.json"

// EXPLINATION: parsing json file into struct
type Config struct {
	Host    string `json:"host"`
	Port    int    `json:"port"`
	TimeOut int    `json:"timeout"`
}

// WARNING:  type cfgKey is a private type
// to avoid key collision, preserves typesaftey
type cfgKey string

const cfgInterfaceKey cfgKey = "cfg_interface"

func loadConfig(ctx context.Context) (context.Context, error) {
	var data Config

	jsonData, err := os.ReadFile(configfile)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(jsonData, &data)
	if err != nil {
		return nil, err
	}

	ctx = context.WithValue(ctx, cfgInterfaceKey, data)

	return ctx, nil
}
