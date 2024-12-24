package web

import (
	"fmt"
	"net/http"

	"github.com/mightynerd/hit/db"
	"github.com/mightynerd/hit/spotify"
)

func (web *Web) getRedirectURL() string {
	return web.serviceURL + "/callback"
}

func (web *Web) Login(w http.ResponseWriter, req *http.Request) {
	url := "https://accounts.spotify.com/authorize?"
	url += "response_type=code&"
	url += "client_id=" + web.spotifyClientId + "&"
	url += "scope=user-modify-playback-state playlist-read-private playlist-read-collaborative&"
	url += "redirect_uri=" + web.getRedirectURL() + "&"
	url += "state=123"

	if req.Method == "GET" {
		http.Redirect(w, req, url, http.StatusSeeOther)
	}
}

func (web *Web) getAccessToken(code string) (string, error) {
	spotifyApp := spotify.NewSpotifyApp(web.spotifyClientId, web.spotifyClientSecret)

	token, err := spotifyApp.GetToken(code, web.getRedirectURL())
	return token, err
}

func (web *Web) Callback(w http.ResponseWriter, req *http.Request) {
	code := req.URL.Query().Get("code")
	state := req.URL.Query().Get("state")
	qerror := req.URL.Query().Get("error")

	if len(qerror) > 0 {
		fmt.Printf("Received error: %s", qerror)
		http.Error(w, qerror, http.StatusInternalServerError)
		return
	}

	fmt.Printf("Code: %s, State: %s", code, state)
	token, err := web.getAccessToken(code)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Could not get access token", http.StatusInternalServerError)
		return
	}

	s := spotify.NewSpotify(token)
	me, err := s.Me()
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Could not get user", http.StatusInternalServerError)
		return
	}

	user := &db.User{
		Name:      me.DisplayName,
		SpotifyId: &me.ID,
		Token:     &token,
	}

	id, err := web.db.PutUser(user)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Could not create user", http.StatusInternalServerError)
		return
	}

	jwt, err := web.createSignedUserJWT(id)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Could not sign JWT", http.StatusInternalServerError)
	}

	fmt.Fprintf(w, jwt)
}
