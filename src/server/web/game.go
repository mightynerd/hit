package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mightynerd/hit/db"
	"github.com/mightynerd/hit/game"
	"github.com/mightynerd/hit/spotify"
)

type CreateGameBody struct {
	PlaylistId string `json:"playlist_id" binding:"required"`
}

func (web *Web) CreateGame(c *gin.Context) {
	var body CreateGameBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	playlist, err := web.db.GetPlaylistById(body.PlaylistId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not get playlist"})
		return
	}

	userInterface, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Missing user"})
		return
	}

	user := userInterface.(*db.User)

	game := &db.Game{
		UserID:     user.ID,
		PlaylistID: playlist.ID,
	}

	gameId, err := web.db.CreateGame(game)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create game"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": gameId})
}

func (web *Web) AdvanceGame(c *gin.Context) {
	gameId := c.Param("game_id")

	userInterface, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Missing user"})
		return
	}

	user := userInterface.(*db.User)

	dbGame, err := web.db.GetGameById(gameId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not get game"})
		return
	}

	if dbGame.UserID != user.ID {
		c.JSON(http.StatusNotFound, gin.H{"error": "Game not found"})
		return
	}

	spotify := spotify.FromUser(user)
	game := game.NewGame(spotify, dbGame, web.db)

	track, err := game.Advance()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not advance game"})
		return
	}

	c.JSON(http.StatusOK, track)
}
