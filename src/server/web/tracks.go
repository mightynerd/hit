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

func (web *Web) UpdateTrack(c *gin.Context) {
	trackId := c.Param("track_id")

	userInterface, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Missing user"})
		return
	}
	user := userInterface.(*db.User)

	track, err := web.db.GetTrackById(trackId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not get track"})
		return
	}

	playlist, err := web.db.GetPlaylistById(track.PlaylistID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not get playlist"})
		return
	}

	if user.ID != playlist.UserID {
		c.JSON(http.StatusNotFound, gin.H{"error": "Could not find track"})
		return
	}

	var body db.Track
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
	}

	if len(body.Title) > 0 {
		track.Title = body.Title
	}
	if len(body.Artist) > 0 {
		track.Artist = body.Artist
	}
	if body.Year > 0 {
		track.Year = body.Year
	}

	updated, err := web.db.UpdateTrack(track)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update track"})
		return
	}

	c.JSON(http.StatusOK, updated)
}

func (web *Web) DeleteTrack(c *gin.Context) {
	trackId := c.Param("track_id")

	userInterface, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Missing user"})
		return
	}
	user := userInterface.(*db.User)

	track, err := web.db.GetTrackById(trackId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not get track"})
		return
	}

	playlist, err := web.db.GetPlaylistById(track.PlaylistID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not get playlist"})
		return
	}

	if user.ID != playlist.UserID {
		c.JSON(http.StatusNotFound, gin.H{"error": "Could not find track"})
		return
	}

	err = web.db.DeleteTrack(trackId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": "Failed to delete track"})
		return
	}

	c.Status(http.StatusOK)
}
