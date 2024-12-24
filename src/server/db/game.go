package db

import (
	"fmt"
	"time"

	"github.com/georgysavva/scany/pgxscan"
)

type Game struct {
	ID         string    `db:"id"`
	CreatedAt  time.Time `db:"created_at"`
	UserID     string    `db:"user_id"`
	PlaylistID string    `db:"playlist_id"`
}

func (db *DB) CreateGame(game *Game) (gameID string, err error) {
	query := `
		INSERT INTO games (user_id, playlist_id)
		VALUES ($1, $2)
		RETURNING id;
	`

	err = db.pool.QueryRow(*db.ctx, query,
		game.UserID,
		game.PlaylistID,
	).Scan(&gameID)

	if err != nil {
		fmt.Println(err)
		return "", fmt.Errorf("unable to insert game")
	}

	return gameID, nil
}

func (db *DB) GetGameById(gameId string) (*Game, error) {
	fmt.Println("getting game by id", gameId)
	query := `
		SELECT * FROM games
		WHERE id = $1
	`

	var game Game
	err := pgxscan.Get(*db.ctx, db.pool, &game, query, gameId)
	if err != nil {
		fmt.Println("unable to get game", err)
		return nil, fmt.Errorf("unable to get game")
	}

	return &game, nil
}
