package main

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	PGConnectionString  string `json:"pg_connection_string"`
	SpotifyClientId     string `json:"spotify_client_id"`
	SpotifyClientSecret string `json:"spotify_client_secret"`
	ServiceUrl          string `json:"service_url"`
	DiscogsAPIKey       string `json:"discogs_api_key"`
	JWTSecret           string `json:"jwt_secret"`
}

func LoadConfig(file string) *Config {
	data, err := os.ReadFile(file)
	if err != nil {
		log.Fatal("Failed to read config file", err)
	}

	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		log.Fatal("Failed to parse config file", err)
	}

	return &config
}
