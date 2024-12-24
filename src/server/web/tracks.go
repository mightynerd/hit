package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mightynerd/hit/db"
)

func (web *Web) GetTracks(c *gin.Context) {
	qPage, _ := c.GetQuery("page")
	qSize, _ := c.GetQuery("size")
	page, size := parsePagination(qPage, qSize, 0, 20)

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

	if playlist.UserID != user.ID {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
		return
	}

	tracks, err := web.db.GetTracks(playlist.ID, page, size)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not get tracks"})
		return
	}

	c.JSON(http.StatusOK, tracks)
}
