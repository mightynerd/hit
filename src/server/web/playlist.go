package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mightynerd/hit/db"
)

func (web *Web) DeletePlaylist(c *gin.Context) {
	playlistId := c.Param("playlist_id")

	userInterface, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Missing user"})
		return
	}
	user := userInterface.(*db.User)

	playlist, err := web.DB.GetPlaylistById(playlistId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not get playlist"})
		return
	}

	if user.ID != playlist.UserID {
		c.JSON(http.StatusNotFound, gin.H{"error": "Could not find playlist"})
		return
	}

	err = web.DB.DeletePlaylist(playlistId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": "Failed to delete playlist"})
		return
	}

	c.Status(http.StatusOK)
}
