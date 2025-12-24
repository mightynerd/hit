package web

import (
	"crypto/sha512"

	"github.com/danielgtaylor/huma/v2"
	"github.com/mightynerd/hit/db"
	"github.com/mightynerd/hit/discogs"
)

type Web struct {
	API                 huma.API
	DB                  *db.DB
	ServiceURL          string
	SpotifyClientId     string
	SpotifyClientSecret string
	JWTSecret           []byte
	Discogs             *discogs.DiscogsConfig
}

func NewWeb(
	api huma.API,
	db *db.DB,
	serviceURL string,
	spotifyClientId string,
	spotifyClientSecret string,
	discogs *discogs.DiscogsConfig,
	jwtSecret string,
) *Web {
	jwtSecretHash := sha512.Sum384([]byte(jwtSecret))

	web := &Web{
		API:                 api,
		DB:                  db,
		ServiceURL:          serviceURL,
		SpotifyClientId:     spotifyClientId,
		SpotifyClientSecret: spotifyClientSecret,
		Discogs:             discogs,
		JWTSecret:           jwtSecretHash[:],
	}

	return web
}
