package spotify

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

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

func (s *Spotify) newAPIRequest(method string, path string, body io.Reader) (*http.Request, error) {
	request, err := http.NewRequest(method, APIURL+path, body)
	if err != nil {
		return nil, err
	}
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", s.token))

	return request, nil

}

func (s *Spotify) get(path string, respBody any) error {
	request, err := s.newAPIRequest("GET", "/me", nil)
	if err != nil {
		return err
	}

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return err
	}

	if err := json.Unmarshal(body, respBody); err != nil {
		return err
	}

	return nil
}
