package web_tracks

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/mightynerd/hit/web"
)

func RegisterRoutes(web *web.Web) {
	group := huma.NewGroup(web.API, "/playlists/:playlistId/tracks")
	group.UseMiddleware(web.AuthMiddleware)
	tags := []string{"Tracks"}

	huma.Register(group,
		huma.Operation{
			Method: http.MethodGet,
			Path:   "",
			Tags:   tags,
			Security: []map[string][]string{
				{"bearerAuth": {}},
			},
		}, getTracks(web))

	huma.Register(group,
		huma.Operation{
			Method: http.MethodPatch,
			Path:   "",
			Tags:   tags,
			Security: []map[string][]string{
				{"bearerAuth": {}},
			},
		}, updateTrack(web))

	huma.Register(group,
		huma.Operation{
			Method: http.MethodDelete,
			Path:   "/:trackId",
			Tags:   tags,
			Security: []map[string][]string{
				{"bearerAuth": {}},
			},
		}, deleteTrack(web))

}
