package game

import (
	"github.com/mightynerd/hit/db"
	"github.com/mightynerd/hit/spotify"
)

type Game struct {
	spotify *spotify.Spotify
	game    *db.Game
	db      *db.DB
}

func NewGame(spotify *spotify.Spotify, dbGame *db.Game, db *db.DB) *Game {
	game := &Game{
		spotify: spotify,
		game:    dbGame,
		db:      db,
	}

	return game
}

func (game *Game) Advance() (*db.Track, error) {
	next, err := game.db.GetUniqueTrack(game.game.PlaylistID, game.game.ID)

	if err != nil {
		return nil, err
	}

	err = game.db.CreateGameTrack(game.game.ID, next.ID)

	if err != nil {
		return nil, err
	}

	err = game.spotify.Play(next.SpotifyURI)

	if err != nil {
		return nil, err
	}

	return next, nil
}
