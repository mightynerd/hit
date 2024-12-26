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
	AllowOrigin         string `json:"allow_origin"`
}

func loadConfigFromFile(file string) *Config {
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

func getRequiredEnv(name string) string {
	env := os.Getenv(name)
	if len(env) < 1 {
		log.Fatal("Missing env variable", name)
	}
	return env
}

func loadConfigFromEnv() *Config {
	config := &Config{
		PGConnectionString:  getRequiredEnv("PG_CONNECTION_STRING"),
		SpotifyClientId:     getRequiredEnv("SPOTIFY_CLIENT_ID"),
		SpotifyClientSecret: getRequiredEnv("SPOTIFY_CLIENT_SECRET"),
		ServiceUrl:          getRequiredEnv("SERVICE_URL"),
		DiscogsAPIKey:       getRequiredEnv("DISCOGS_API_KEY"),
		JWTSecret:           getRequiredEnv("JWT_SECRET"),
		AllowOrigin:         getRequiredEnv("ALLOW_ORIGIN"),
	}

	return config
}

func LoadConfig(file string) *Config {
	if _, err := os.Stat(file); err == nil {
		return loadConfigFromFile(file)
	} else {
		return loadConfigFromEnv()
	}
}
