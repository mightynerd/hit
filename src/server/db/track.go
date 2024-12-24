package db

import (
	"fmt"
	"time"

	"github.com/georgysavva/scany/pgxscan"
)

type Track struct {
	ID         string    `db:"id" json:"id"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
	PlaylistID string    `db:"playlist_id" json:"playlist_id"`
	Title      string    `db:"title" json:"title"`
	Artist     string    `db:"artist" json:"artist"`
	Year       int       `db:"year" json:"year"`
	SpotifyURI string    `db:"spotify_uri" json:"spotify_uri"`
}

func (db *DB) CreateTrack(track *Track) (trackId string, err error) {
	query := `
		INSERT INTO tracks (playlist_id, title, artist, year, spotify_uri)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id;
	`

	err = db.pool.QueryRow(*db.ctx, query,
		track.PlaylistID,
		track.Title,
		track.Artist,
		track.Year,
		track.SpotifyURI,
	).Scan(&trackId)

	if err != nil {
		fmt.Println("failed to insert track", err)
		return "", fmt.Errorf("unable to insert track")
	}

	return trackId, nil
}

/*
Select a track from a certain playlistId such that it does not exist
in game_tracks for a certain gameId
*/
func (db *DB) GetUniqueTrack(playlistId string, gameId string) (*Track, error) {
	fmt.Println("getting unique track from playlist", playlistId, "for game", gameId)
	query := `
		SELECT * FROM tracks
		WHERE playlist_id = $1
		AND id NOT IN (
			SELECT track_id
			FROM game_tracks
			WHERE game_id = $2
		)
		ORDER BY RANDOM()
		LIMIT 1
	`

	var track Track
	err := pgxscan.Get(*db.ctx, db.pool, &track, query, playlistId, gameId)
	if err != nil {
		fmt.Println("failed to get unique track", err)
		return nil, err
	}

	return &track, nil
}

func (db *DB) GetTracks(playlistId string, page int, size int) (*[]Track, error) {
	query := `
		SELECT * from tracks
		WHERE playlist_id = $1
		ORDER BY created_at DESC
		OFFSET $2
		LIMIT $3
	`
	var tracks []Track
	err := pgxscan.Select(*db.ctx, db.pool, &tracks, query, playlistId, page*size, size)

	if err != nil {
		fmt.Println("failed to get tracks", err)
		return nil, err
	}

	return &tracks, nil
}
