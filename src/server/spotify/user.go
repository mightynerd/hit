package spotify

import "fmt"

type SpotifyUser struct {
	Country         string `json:"country"`
	DisplayName     string `json:"display_name"`
	Email           string `json:"email"`
	ExplicitContent struct {
		FilterEnabled bool `json:"filter_enabled"`
		FilterLocked  bool `json:"filter_locked"`
	} `json:"explicit_content"`
	ExternalUrls struct {
		Spotify string `json:"spotify"`
	} `json:"external_urls"`
	Followers struct {
		Href  string `json:"href"`
		Total int    `json:"total"`
	} `json:"followers"`
	Href   string `json:"href"`
	ID     string `json:"id"`
	Images []struct {
		URL    string `json:"url"`
		Height int    `json:"height"`
		Width  int    `json:"width"`
	} `json:"images"`
	Product string `json:"product"`
	Type    string `json:"type"`
	URI     string `json:"uri"`
}

func (s *Spotify) Me() (*SpotifyUser, error) {
	var me SpotifyUser
	resp, err := s.newRequest().SetResult(&me).Get("/me")

	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("failed to get me %s", string(resp.Body()))
	}

	return &me, nil
}
