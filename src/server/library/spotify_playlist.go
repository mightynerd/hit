package library

import (
	"fmt"

	"github.com/mightynerd/hit/db"
	"github.com/mightynerd/hit/discogs"
	"github.com/mightynerd/hit/spotify"
)

type Library struct {
	db      *db.DB
	spotify *spotify.Spotify
	discogs *discogs.DiscogsConfig
}

func NewLibrary(db *db.DB, spotify *spotify.Spotify, discogs *discogs.DiscogsConfig) *Library {
	library := &Library{
		db:      db,
		spotify: spotify,
		discogs: discogs,
	}

	return library
}

func (lib *Library) ImportSpotifyPlaylist(playlistId string, spotifyPlaylistId string) error {
	fmt.Println("Importing spotify playlist", spotifyPlaylistId, "to", playlistId)
	tracks, err := lib.spotify.GetPlaylistItems(spotifyPlaylistId)
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Println("Got", len(*tracks), "tracks, enhancing with discogs")

	fmt.Println("Inserting into db")
	for _, track := range *tracks {
		lib.discogs.EnhanceYear(&track)
		track.PlaylistID = playlistId
		_, err := lib.db.CreateTrack(&track)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}

	return nil
}
