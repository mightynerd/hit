package web

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mightynerd/hit/db"
	"github.com/mightynerd/hit/library"
	"github.com/mightynerd/hit/spotify"
)

type CreatePlaylistBody struct {
	Name string `json:"name"`
	From struct {
		Source string `json:"source"`
		Id     string `json:"id"`
	} `json:"from"`
}

func (web *Web) handleImport(user *db.User, body *CreatePlaylistBody, playlistId string) error {
	if body.From.Source == "spotify_playlist" {
		spotify := spotify.FromUser(user)
		lib := library.NewLibrary(web.db, spotify, web.discogs)
		err := lib.ImportSpotifyPlaylist(playlistId, body.From.Id)
		if err != nil {
			fmt.Println(err)
			web.db.UpdatePlaylistStatus(playlistId, db.PlaylistStatusFailed)
			return err
		}
	}

	web.db.UpdatePlaylistStatus(playlistId, db.PlaylistStatusActive)
	return nil
}

func (web *Web) CreatePlaylist(c *gin.Context) {
	var body CreatePlaylistBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
	}

	userInterface, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Missing user"})
		return
	}

	user := userInterface.(*db.User)

	playlist := &db.Playlist{
		UserID: user.ID,
		Name:   body.Name,
		Status: db.PlaylistStatusImporting,
	}

	playlistId, err := web.db.CreatePlaylist(playlist)

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create playlist"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": playlistId})

	go web.handleImport(user, &body, playlistId)
}

func (web *Web) GetPlaylists(c *gin.Context) {
	qPage, _ := c.GetQuery("page")
	qSize, _ := c.GetQuery("size")
	page, size := parsePagination(qPage, qSize, 0, 20)

	userInterface, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Missing user"})
		return
	}

	user := userInterface.(*db.User)

	playlists, err := web.db.GetPlaylists(user.ID, page, size)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not get playlists"})
		return
	}

	c.JSON(http.StatusOK, playlists)
}

func (web *Web) DeletePlaylist(c *gin.Context) {
	playlistId := c.Param("playlist_id")

	userInterface, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Missing user"})
		return
	}
	user := userInterface.(*db.User)

	playlist, err := web.db.GetPlaylistById(playlistId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not get playlist"})
		return
	}

	if user.ID != playlist.UserID {
		c.JSON(http.StatusNotFound, gin.H{"error": "Could not find playlist"})
		return
	}

	err = web.db.DeletePlaylist(playlistId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": "Failed to delete playlist"})
		return
	}

	c.Status(http.StatusOK)
}
