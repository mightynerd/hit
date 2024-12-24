package db

import (
	"fmt"
	"time"

	"github.com/georgysavva/scany/pgxscan"
)

type Playlist struct {
	ID        string    `db:"id" json:"id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	Name      string    `db:"name" json:"name"`
	UserID    string    `db:"user_id" json:"user_id"`
}

func (db *DB) CreatePlaylist(playlist *Playlist) (playlistId string, err error) {
	query := `
		INSERT INTO playlists (name, user_id)
		VALUES ($1, $2)
		RETURNING id;
	`

	err = db.pool.QueryRow(*db.ctx, query, playlist.Name, playlist.UserID).Scan(&playlistId)
	if err != nil {
		fmt.Println(err)
		return "", fmt.Errorf("unable to insert playlist")
	}

	return playlistId, nil
}

func (db *DB) GetPlaylistById(playlistId string) (playlist Playlist, err error) {
	query := `
		SELECT * FROM playlists
		WHERE id = $1
	`

	err = pgxscan.Get(*db.ctx, db.pool, &playlist, query, playlistId)
	if err != nil {
		fmt.Println(err)
		return playlist, fmt.Errorf("unable to get playlist")
	}

	return playlist, nil
}

func (db *DB) GetPlaylists(userId string, page int, size int) (*[]Playlist, error) {
	query := `
		SELECT * from playlists
		WHERE user_id = $1
		OFFSET $2
		LIMIT $3
		ORDER BY created_at DESC
	`
	var playlists []Playlist
	err := pgxscan.Select(*db.ctx, db.pool, &playlists, query, userId, page, size)

	if err != nil {
		fmt.Println("failed to get playlists", err)
		return nil, err
	}

	return &playlists, nil
}
