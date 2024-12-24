package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/mux"
	"github.com/mightynerd/hit/db"
	"github.com/mightynerd/hit/discogs"
	"github.com/mightynerd/hit/web"
)

type Server struct {
	ctx    *context.Context
	db     *db.DB
	config *Config
	server *http.Server
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

	router := mux.NewRouter()
	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	server := &Server{
		ctx:    &ctx,
		config: config,
		server: httpServer,
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

	r.GET("/login", web.Login)
	r.GET("/callback", web.Callback)

	authorizedGroup := r.Group("")
	authorizedGroup.Use(web.AuthMiddleware())

	authorizedGroup.POST("/playlists", web.CreatePlaylist)
	authorizedGroup.POST("/games", web.CreateGame)
	authorizedGroup.POST("/games/:game_id/advance", web.AdvanceGame)

	r.Run(":8080")
}
