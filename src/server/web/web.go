package web

import (
	"github.com/mightynerd/hit/db"
	"github.com/mightynerd/hit/discogs"
)

type Web struct {
	db                  *db.DB
	serviceURL          string
	spotifyClientId     string
	spotifyClientSecret string
	jwtSecret           []byte
	discogs             *discogs.DiscogsConfig
}

func NewWeb(db *db.DB, serviceURL string, spotifyClientId string, spotifyClientSecret string, discogs *discogs.DiscogsConfig) *Web {
	web := &Web{
		db:                  db,
		serviceURL:          serviceURL,
		spotifyClientId:     spotifyClientId,
		spotifyClientSecret: spotifyClientSecret,
		discogs:             discogs,
	}

	return web
}
