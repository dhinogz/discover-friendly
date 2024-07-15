package router

import (
	"os"
	"strings"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"github.com/pocketbase/pocketbase/core"
	"golang.org/x/oauth2"

	"github.com/dhinogz/spotify-test/assets"
	"github.com/dhinogz/spotify-test/htmx"
)

type OAuthSettings struct {
	ClientID     string
	ClientSecret string
	AuthURL      string
	TokenURL     string
	RedirectURL  string
	Scopes       string
}

type AppRouter struct {
	App    core.App
	Router *echo.Echo
	OAuth2 *oauth2.Config
}

func buildConfig(cfg *OAuthSettings) *oauth2.Config {
	endpoint := oauth2.Endpoint{
		AuthURL:  cfg.AuthURL,
		TokenURL: cfg.TokenURL,
	}
	return &oauth2.Config{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		Endpoint:     endpoint,
		RedirectURL:  cfg.RedirectURL,
		Scopes:       strings.Split(cfg.Scopes, ","),
	}
}

func NewAppRouter(e *core.ServeEvent) *AppRouter {

	oauthSettings := &OAuthSettings{
		ClientID:     os.Getenv("SPOTIFY_CLIENT_ID"),
		ClientSecret: os.Getenv("SPOTIFY_CLIENT_SECRET"),
		AuthURL:      os.Getenv("SPOTIFY_AUTH_URL"),
		TokenURL:     os.Getenv("SPOTIFY_TOKEN_URL"),
		RedirectURL:  os.Getenv("SPOTIFY_REDIRECT_URL"),
		Scopes:       os.Getenv("SPOTIFY_SCOPES"),
	}

	o := buildConfig(oauthSettings)

	return &AppRouter{
		App:    e.App,
		Router: e.Router,
		OAuth2: o,
	}
}

func (ar *AppRouter) SetupRoutes(live bool) error {
	ar.Router.Use(middleware.Logger())
	ar.Router.HTTPErrorHandler = htmx.WrapDefaultErrorHandler(ar.Router.HTTPErrorHandler)
	ar.Router.GET("/static/*", assets.AssetsHandler(ar.App.Logger(), live), middleware.Gzip())

	ar.Router.Use(ar.LoadAuthContextFromCookie())
	ar.Router.GET("/", ar.GetHome)
	ar.Router.GET("/search", ar.GetHome, ar.LoadAuthContextFromCookie(), ar.LoadSpotifyAuthMiddleware())

	// Search htmx response
	ar.Router.POST("/spotify-search", ar.HandleSpotifySearch, ar.LoadAuthContextFromCookie(), ar.LoadSpotifyAuthMiddleware())
	ar.Router.GET("/tracks/:uri", ar.HandlePlayTrack, ar.LoadAuthContextFromCookie(), ar.LoadSpotifyAuthMiddleware())

	// Auth handlers
	ar.Router.GET("/login", ar.GetLogin)

	// Spotify oauth
	ar.Router.GET("/oauth/spotify/connect", ar.HandleOAuthConnect)
	ar.Router.GET("/oauth/spotify/callback", ar.HandleOAuthCallback)

	// Logout handlers
	ar.Router.POST("/logout", ar.PostLogout)

	// Errors handlers. Middleware reroutes to here
	ar.Router.GET("/error", ar.GetError)

	return nil
}
