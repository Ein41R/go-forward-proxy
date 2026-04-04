package main

import (
	"context"
	"encoding/json"
	"os"
)

var configfile = "config.json"

type Config struct {
	Host    string `json:"host"`
	Port    int    `json:"port"`
	TimeOut int    `json:"timeout"`
}

// private type to avoid key collisions in context
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
	// cfg := (*ctx).Value("cfg_interface").(Config).Host

	// log.Printf("Config loaded: %+v\n", reflect.TypeOf(cfg))

	return ctx, nil
}
