package spotify

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/mightynerd/hit/db"
)

type GetPlaylistItemsResponse struct {
	Href     string `json:"href"`
	Limit    int    `json:"limit"`
	Next     string `json:"next"`
	Offset   int    `json:"offset"`
	Previous string `json:"previous"`
	Total    int    `json:"total"`
	Items    []struct {
		AddedAt string `json:"added_at"`
		AddedBy struct {
			ExternalUrls struct {
				Spotify string `json:"spotify"`
			} `json:"external_urls"`
			Followers struct {
				Href  string `json:"href"`
				Total int    `json:"total"`
			} `json:"followers"`
			Href string `json:"href"`
			ID   string `json:"id"`
			Type string `json:"type"`
			URI  string `json:"uri"`
		} `json:"added_by"`
		IsLocal bool `json:"is_local"`
		Track   struct {
			Album struct {
				AlbumType        string   `json:"album_type"`
				TotalTracks      int      `json:"total_tracks"`
				AvailableMarkets []string `json:"available_markets"`
				ExternalUrls     struct {
					Spotify string `json:"spotify"`
				} `json:"external_urls"`
				Href   string `json:"href"`
				ID     string `json:"id"`
				Images []struct {
					URL    string `json:"url"`
					Height int    `json:"height"`
					Width  int    `json:"width"`
				} `json:"images"`
				Name                 string `json:"name"`
				ReleaseDate          string `json:"release_date"`
				ReleaseDatePrecision string `json:"release_date_precision"`
				Restrictions         struct {
					Reason string `json:"reason"`
				} `json:"restrictions"`
				Type    string `json:"type"`
				URI     string `json:"uri"`
				Artists []struct {
					ExternalUrls struct {
						Spotify string `json:"spotify"`
					} `json:"external_urls"`
					Href string `json:"href"`
					ID   string `json:"id"`
					Name string `json:"name"`
					Type string `json:"type"`
					URI  string `json:"uri"`
				} `json:"artists"`
			} `json:"album"`
			Artists []struct {
				ExternalUrls struct {
					Spotify string `json:"spotify"`
				} `json:"external_urls"`
				Href string `json:"href"`
				ID   string `json:"id"`
				Name string `json:"name"`
				Type string `json:"type"`
				URI  string `json:"uri"`
			} `json:"artists"`
			AvailableMarkets []string `json:"available_markets"`
			DiscNumber       int      `json:"disc_number"`
			DurationMs       int      `json:"duration_ms"`
			Explicit         bool     `json:"explicit"`
			ExternalIds      struct {
				Isrc string `json:"isrc"`
				Ean  string `json:"ean"`
				Upc  string `json:"upc"`
			} `json:"external_ids"`
			ExternalUrls struct {
				Spotify string `json:"spotify"`
			} `json:"external_urls"`
			Href       string `json:"href"`
			ID         string `json:"id"`
			IsPlayable bool   `json:"is_playable"`
			LinkedFrom struct {
			} `json:"linked_from"`
			Restrictions struct {
				Reason string `json:"reason"`
			} `json:"restrictions"`
			Name        string `json:"name"`
			Popularity  int    `json:"popularity"`
			PreviewURL  string `json:"preview_url"`
			TrackNumber int    `json:"track_number"`
			Type        string `json:"type"`
			URI         string `json:"uri"`
			IsLocal     bool   `json:"is_local"`
		} `json:"track"`
	} `json:"items"`
}

func (spotify *Spotify) GetPlaylistItems(playlistId string) (*[]db.Track, error) {
	tracks := []db.Track{}
	limit := 20
	offset := 0

	for {
		fmt.Println("Requesting spotify playlist", playlistId, " offset ", offset, " limit ", limit)
		response := &GetPlaylistItemsResponse{}
		_, err := spotify.newRequest().
			SetQueryParam("market", "SE").
			SetQueryParam("limit", strconv.Itoa(limit)).
			SetQueryParam("offset", strconv.Itoa(offset)).
			SetPathParam("playlist_id", playlistId).
			SetResult(response).
			Get("/playlists/{playlist_id}/tracks")

		if err != nil {
			return nil, err
		}

		for _, item := range response.Items {
			yearParseResult, err := parseYear(item.Track.Album.ReleaseDate)
			if err != nil {
				yearParseResult = 0
			}
			track := db.Track{
				Title:      item.Track.Name,
				Artist:     item.Track.Artists[0].Name,
				Year:       yearParseResult,
				SpotifyURI: item.Track.URI,
			}
			tracks = append(tracks, track)
		}

		if offset >= response.Total {
			break
		}

		offset += limit
	}

	return &tracks, nil
}

func parseYear(date string) (int, error) {
	yearStr := strings.Split(date, "-")[0]
	year, err := strconv.Atoi(yearStr)

	if err != nil {
		return 0, err
	}

	return year, nil
}
