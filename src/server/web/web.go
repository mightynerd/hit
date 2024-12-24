package web

import (
	"crypto/sha512"
	"fmt"

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

func NewWeb(
	db *db.DB,
	serviceURL string,
	spotifyClientId string,
	spotifyClientSecret string,
	discogs *discogs.DiscogsConfig,
	jwtSecret string) *Web {

	jwtSecretHash := sha512.Sum384([]byte(jwtSecret))
	fmt.Println(jwtSecret, jwtSecretHash)

	web := &Web{
		db:                  db,
		serviceURL:          serviceURL,
		spotifyClientId:     spotifyClientId,
		spotifyClientSecret: spotifyClientSecret,
		discogs:             discogs,
		jwtSecret:           jwtSecretHash[:],
	}

	return web
}
