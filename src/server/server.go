package main

import (
	"context"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/mightynerd/hit/db"
	"github.com/mightynerd/hit/discogs"
	"github.com/mightynerd/hit/web"
)

type Server struct {
	ctx    *context.Context
	db     *db.DB
	config *Config
}

func (s *Server) ConnectToDb() {
	db, err := db.Connect(*s.ctx, s.config.PGConnectionString)
	if err != nil {
		log.Fatal("Failed to connect to db", err)
	}

	s.db = db
}

func main() {
	// Load config
	config := LoadConfig("config.json")

	// Migrate
	migrator := db.NewMigrator(config.PGConnectionString)
	migrator.Migrate()

	ctx := context.Background()

	server := &Server{
		ctx:    &ctx,
		config: config,
	}

	// Connect to DB
	server.ConnectToDb()

	discogs := discogs.NewDiscogsConfig(config.DiscogsAPIKey)

	web := web.NewWeb(
		server.db,
		server.config.ServiceUrl,
		server.config.SpotifyClientId,
		server.config.SpotifyClientSecret,
		discogs,
		config.JWTSecret)

	r := gin.Default()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost", "http://localhost:5173"}
	corsConfig.AllowCredentials = true
	corsConfig.AddAllowHeaders("Authorization")
	r.Use(cors.New(corsConfig))

	r.GET("/login", web.Login)
	r.GET("/callback", web.Callback)

	authorizedGroup := r.Group("")

	authorizedGroup.Use(web.AuthMiddleware())
	authorizedGroup.GET("/playlists", web.GetPlaylists)
	authorizedGroup.POST("/playlists", web.CreatePlaylist)
	authorizedGroup.DELETE("/playlists/:playlist_id", web.DeletePlaylist)

	authorizedGroup.GET("/playlists/:playlist_id/tracks", web.GetTracks)
	authorizedGroup.PATCH("/playlists/:playlist_id/tracks/:track_id", web.UpdateTrack)
	authorizedGroup.DELETE("/playlists/:playlist_id/tracks/:track_id", web.DeleteTrack)

	authorizedGroup.POST("/games", web.CreateGame)
	authorizedGroup.POST("/games/:game_id/advance", web.AdvanceGame)

	r.Run(":8080")
}
