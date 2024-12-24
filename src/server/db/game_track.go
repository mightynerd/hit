package db

import (
	"fmt"
)

type GameTrack struct {
	GameID  string
	TrackID string
}

func (db *DB) CreateGameTrack(gameID string, trackID string) error {
	query := `
		INSERT INTO game_tracks (game_id, track_id)
		VALUES ($1, $2)
		RETURNING *;
	`

	err := db.pool.QueryRow(*db.ctx, query, gameID, trackID).Scan(&gameID, &trackID)

	if err != nil {
		fmt.Println("unable to insert game track", err)
		return fmt.Errorf("unable to insert game track")
	}

	return nil
}
