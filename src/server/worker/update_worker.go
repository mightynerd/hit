package worker

import (
	"context"
	"fmt"
	"time"

	"github.com/mightynerd/hit/db"
	"github.com/mightynerd/hit/discogs"
)

type UpdateWorker struct {
	DB      db.DB
	Discogs discogs.DiscogsConfig
}

func (w UpdateWorker) RunUpdateWorker(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			track, err := w.DB.GetTrackToEnhance()

			if err != nil {
				fmt.Printf("failed to get track to enhance %v\n", err)
				continue
			}

			if track == nil {
				continue
			}

			err = w.Discogs.EnhanceYear(track)
			if err != nil {
				fmt.Printf("failed to enhance track %s\n", track.ID)
			}

			track.Enhanced = true
			w.DB.UpdateTrack(track)

			unenhanced, err := w.DB.CountUnenhancedTracks(track.PlaylistID)
			if err != nil {
				continue
			}

			if unenhanced == 0 {
				w.DB.UpdatePlaylistStatus(track.PlaylistID, db.PlaylistStatusActive)
			}
		}
	}
}
