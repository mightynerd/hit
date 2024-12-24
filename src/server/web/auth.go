package web

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mightynerd/hit/db"
	"github.com/mightynerd/hit/spotify"
)

func (web *Web) getRedirectURL() string {
	return web.serviceURL + "/callback"
}

func (web *Web) getLoginURL() string {
	url := "https://accounts.spotify.com/authorize?"
	url += "response_type=code&"
	url += "client_id=" + web.spotifyClientId + "&"
	url += "scope=user-modify-playback-state playlist-read-private playlist-read-collaborative&"
	url += "redirect_uri=" + web.getRedirectURL() + "&"
	url += "state=123"
	return url
}

func (web *Web) Login(c *gin.Context) {
	c.Redirect(http.StatusSeeOther, web.getLoginURL())
}

func (web *Web) getAccessToken(code string) (string, error) {
	spotifyApp := spotify.NewSpotifyApp(web.spotifyClientId, web.spotifyClientSecret)

	token, err := spotifyApp.GetToken(code, web.getRedirectURL())
	return token, err
}

func (web *Web) Callback(c *gin.Context) {
	code := c.Query("code")
	//state := c.Query("state")
	qerror := c.Query("error")

	if len(qerror) > 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Callback failed " + qerror})
		return
	}

	token, err := web.getAccessToken(code)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not get access token"})
		return
	}

	s := spotify.NewSpotify(token)
	me, err := s.Me()
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not get user"})
		return
	}
	fmt.Println(me)
	user := &db.User{
		Name:      me.DisplayName,
		SpotifyId: &me.ID,
		Token:     &token,
	}

	id, err := web.db.PutUser(user)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create user"})
		return
	}

	jwt, err := web.createSignedUserJWT(id)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Faild to sign JWT"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": jwt})
}
