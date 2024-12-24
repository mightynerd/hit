package spotify

import (
	"encoding/base64"
	"fmt"

	"github.com/go-resty/resty/v2"
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
	unEncodedAuth := fmt.Sprintf("%s:%s", s.clientId, s.clientSecret)
	encodedAuth := base64.StdEncoding.EncodeToString([]byte(unEncodedAuth))

	var response GetTokenResp
	_, err := resty.New().
		R().
		SetFormData(map[string]string{
			"code":         code,
			"redirect_url": redirectUrl,
			"grant_type":   "authorization_code",
		}).
		SetResult(&response).
		SetHeader("Authorization", "Basic "+encodedAuth).
		Post("https://accounts.spotify.com/api/token")

	if err != nil {
		return "", err
	}

	return response.AccessToken, nil
}
