package spotify

import (
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/mightynerd/hit/db"
)

type Spotify struct {
	token string
}

const APIURL = "https://api.spotify.com/v1"

func NewSpotify(token string) *Spotify {
	spotify := &Spotify{
		token: token,
	}

	return spotify
}

func FromUser(user *db.User) *Spotify {
	return NewSpotify(*user.Token)
}

func (s *Spotify) newRequest() *resty.Request {
	client := resty.New()
	client.SetHeader("Authorization", fmt.Sprintf("Bearer %s", s.token))
	client.BaseURL = APIURL

	return client.R()
}
