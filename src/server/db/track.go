package db

import (
	"errors"
	"fmt"
	"time"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
)

type Track struct {
	ID         string    `db:"id" json:"id"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
	PlaylistID string    `db:"playlist_id" json:"playlist_id"`
	Title      string    `db:"title" json:"title"`
	Artist     string    `db:"artist" json:"artist"`
	Year       int       `db:"year" json:"year"`
	SpotifyURI string    `db:"spotify_uri" json:"spotify_uri"`
	Enhanced   bool      `db:"enhanced"`
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

func (db *DB) CreateTracks(tracks []Track) (int64, error) {
	rows := [][]any{}
	for _, t := range tracks {
		rows = append(rows, []any{
			t.PlaylistID,
			t.Title,
			t.Artist,
			t.Year,
			t.SpotifyURI,
			false,
		})
	}

	count, err := db.pool.CopyFrom(
		*db.ctx,
		pgx.Identifier{"tracks"},
		[]string{"playlist_id", "title", "artist", "year", "spotify_uri", "enhanced"},
		pgx.CopyFromRows(rows),
	)

	return count, err
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

func (db *DB) GetTrackToEnhance() (*Track, error) {
	const query = `
		SELECT * FROM tracks
		WHERE enhanced = false
		ORDER BY created_at ASC
		LIMIT 1
	`

	var track Track
	err := pgxscan.Get(*db.ctx, db.pool, &track, query)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &track, nil
}

func (db *DB) SetTrackEnhanced(trackId string) error {
	const query = `
		UPDATE tracks
		SET enhanced = true
		WHERE id = $1
	`

	_, err := db.pool.Query(*db.ctx, query, trackId)
	return err
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

func (db *DB) GetTrackById(trackId string) (*Track, error) {
	query := `
		SELECT * FROM tracks
		WHERE id = $1
	`

	var track Track
	err := pgxscan.Get(*db.ctx, db.pool, &track, query, trackId)
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("unable to get track")
	}

	return &track, nil
}

func (db *DB) UpdateTrack(track *Track) (*Track, error) {
	query := `
		UPDATE tracks
		SET title = $2, artist = $3, year = $4, enhanced = $5
		WHERE id = $1
		RETURNING *
	`

	var updated Track
	err := pgxscan.Get(*db.ctx, db.pool, &updated, query,
		track.ID,
		track.Title,
		track.Artist,
		track.Year,
		track.Enhanced,
	)
	if err != nil {
		fmt.Println("failed to update track", err)
		return nil, err
	}

	return &updated, nil
}

func (db *DB) CountUnenhancedTracks(playlistId string) (int, error) {
	const query = `
		SELECT COUNT(*)
		FROM tracks
		WHERE playlist_id = $1
		AND enhanced = false
	`

	var count int
	err := pgxscan.Get(*db.ctx, db.pool, &count, query, playlistId)

	return count, err
}

func (db *DB) DeleteTrack(trackId string) error {
	query := `
		DELETE FROM tracks
		WHERE id = $1
	`

	_, err := db.pool.Query(*db.ctx, query, trackId)
	if err != nil {
		fmt.Println("failed to delete track", err)
		return err
	}

	return nil
}
