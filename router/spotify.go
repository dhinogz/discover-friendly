package router

import (
	"log/slog"
	"net/http"

	"github.com/dhinogz/spotify-test/components"
	"github.com/dhinogz/spotify-test/components/shared"
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/zmb3/spotify/v2"
)

func (ar *AppRouter) HandleSpotifySearch(c echo.Context) error {
	rec := c.Get(apis.ContextAuthRecordKey)
	if rec == nil {
		return components.Render(c, http.StatusOK, components.Home(shared.Context{}))
	}

	spotifyClient := c.Get("spotifyClient").(*spotify.Client)

	q := c.FormValue("query")
	if q == "" {
		return c.String(http.StatusOK, "")
	}
	slog.Info("query val", "q", q)
	searchRes, err := spotifyClient.Search(c.Request().Context(), q, spotify.SearchTypeTrack, spotify.Limit(5))
	if err != nil {
		return err
	}
	return components.Render(
		c,
		http.StatusOK,
		components.SpotifySearchResults(searchRes.Tracks.Tracks),
	)
}

func (ar *AppRouter) HandlePlayTrack(c echo.Context) error {
	rec := c.Get(apis.ContextAuthRecordKey)
	if rec == nil {
		return components.Render(c, http.StatusOK, components.Home(shared.Context{}))
	}

	spotifyClient := c.Get("spotifyClient").(*spotify.Client)
	uri := c.PathParam("uri")

	var uris []spotify.URI
	uris = append(uris, spotify.URI(uri))

	err := spotifyClient.PlayOpt(c.Request().Context(), &spotify.PlayOptions{
		URIs: uris,
	})
	if err != nil {
		return c.String(http.StatusOK, err.Error())
	}

	return c.String(http.StatusOK, "now playing")

}

// func (ar *AppRouter) HandleSpotifySearch(c echo.Context) error {
// 	rec := c.Get(apis.ContextAuthRecordKey)
// 	if rec == nil {
// 		return components.Render(c, http.StatusOK, components.Home(shared.Context{}))
// 	}
//
// 	user := c.Get(apis.ContextAuthRecordKey).(*pbmodels.Record)
//
// 	o, err := models.GetOAuthByUserId(ar.App.Dao(), user.Id, "spotify")
// 	if err != nil {
// 		return err
// 	}
//
// 	token := &oauth2.Token{
// 		AccessToken:  o.AccessToken,
// 		TokenType:    o.TokenType,
// 		RefreshToken: o.RefreshToken,
// 		Expiry:       o.Expiry.Time(),
// 	}
//
// 	spotifyClient := spotify.New(spotifyauth.New().Client(c.Request().Context(), token))
// 	slog.Info("token", "t", token)
//
// 	q := c.FormValue("query")
// 	if q == "" {
// 		return c.String(http.StatusOK, "")
// 	}
// 	slog.Info("query val", "q", q)
//
// 	searchRes, err := spotifyClient.Search(c.Request().Context(), q, spotify.SearchTypeTrack)
// 	if err != nil {
// 		return err
// 	}
//
// 	return components.Render(
// 		c,
// 		http.StatusOK,
// 		components.SpotifySearchResults(searchRes.Tracks.Tracks),
// 	)
//
// }
