package spotify

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type SpotifyApp struct {
	clientId     string
	clientSecret string
}

type GetTokenResp struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

func NewSpotifyApp(clientId string, clientSecret string) *SpotifyApp {
	app := &SpotifyApp{
		clientId:     clientId,
		clientSecret: clientSecret,
	}

	return app
}

func (s *SpotifyApp) GetToken(code string, redirectUrl string) (string, error) {
	form := url.Values{}
	form.Add("code", code)
	form.Add("redirect_uri", redirectUrl)
	form.Add("grant_type", "authorization_code")

	client := &http.Client{}
	request, err := http.NewRequest("POST", "https://accounts.spotify.com/api/token", strings.NewReader(form.Encode()))
	if err != nil {
		return "", err
	}

	unEncodedAuth := fmt.Sprintf("%s:%s", s.clientId, s.clientSecret)
	encodedAuth := base64.StdEncoding.EncodeToString([]byte(unEncodedAuth))
	request.Header.Add("Authorization", fmt.Sprintf("Basic %s", encodedAuth))
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(request)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var response GetTokenResp
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", err
	}

	return response.AccessToken, nil
}
