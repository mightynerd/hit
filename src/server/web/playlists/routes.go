package web_playlists

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/mightynerd/hit/web"
)

func RegisterRoutes(web *web.Web) {
	group := huma.NewGroup(web.API, "/playlists")
	group.UseMiddleware(web.AuthMiddleware)
	tags := []string{"Playlists"}

	huma.Register(group,
		huma.Operation{
			Method: http.MethodGet,
			Path:   "",
			Tags:   tags,
			Security: []map[string][]string{
				{"bearerAuth": {}},
			},
		}, getPlaylists(web))

	huma.Register(group,
		huma.Operation{
			Method: http.MethodPost,
			Path:   "",
			Tags:   tags,
			Security: []map[string][]string{
				{"bearerAuth": {}},
			},
		}, createPlaylist(web))

	huma.Register(group,
		huma.Operation{
			Method: http.MethodDelete,
			Path:   "/:playlistId",
			Tags:   tags,
			Security: []map[string][]string{
				{"bearerAuth": {}},
			},
		}, deletePlaylist(web))

}
