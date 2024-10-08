package handlers

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/apis"

	"fmt"
	"github.com/dhinogz/discover-friendly/internal/ui"
	"github.com/dhinogz/discover-friendly/internal/ui/auth"
	"github.com/dhinogz/discover-friendly/internal/ui/shared"
	"log/slog"
	"time"

	"github.com/dhinogz/discover-friendly/internal/models"
	"github.com/gorilla/csrf"
	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
)

// GetLogin returns the login page.
func (ar *AppRouter) GetLogin(c echo.Context) error {
	if c.Get(apis.ContextAuthRecordKey) != nil {
		return c.Redirect(302, "/")
	}

	return ui.Render(c, http.StatusOK, auth.LoginPage(shared.Context{}))
}

// PostLogout logs the user out by clearing the authentication cookie.
func (ar *AppRouter) PostLogout(c echo.Context) error {
	c.SetCookie(&http.Cookie{
		Name:     AuthCookieName,
		Value:    "",
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
		MaxAge:   -1,
	})

	return htmxRedirect(c, "/")
}

func (ar *AppRouter) HandleOAuthConnect(c echo.Context) error {
	state := csrf.Token(c.Request())
	cookie := newCookie(CookieOAuth, state)
	http.SetCookie(c.Response().Writer, cookie)

	url := ar.OAuth2.AuthCodeURL(state)
	slog.Info("state url", "state", state, "url", url)

	return c.Redirect(http.StatusFound, url)
}

func (ar *AppRouter) HandleOAuthCallback(c echo.Context) error {
	state := c.FormValue("state")
	cookie, err := c.Cookie(CookieOAuth)
	if err != nil {
		return fmt.Errorf("getting oauth cookie: %+v", err)
	} else if cookie == nil || cookie.Value != state {
		return fmt.Errorf("invalid state provided: %+v", err)
	}
	cookie.Value = ""
	cookie.Expires = time.Now()
	c.SetCookie(cookie)

	code := c.FormValue("code")
	token, err := ar.OAuth2.Exchange(c.Request().Context(), code)
	if err != nil {
		return err
	}
	spotifyClient := spotify.New(spotifyauth.New().Client(c.Request().Context(), token))
	currentSpotifyUser, err := spotifyClient.CurrentUser(c.Request().Context())
	if err != nil {
		slog.Error("could not get spotify user", "err", err)
		return err
	}

	tracks, err := spotifyClient.CurrentUsersTopTracks(c.Request().Context(), spotify.Limit(5))
	if err != nil {
		slog.Error("could not get current users top tracks", "err", err)
		return err
	}
	seeds := spotify.Seeds{}
	for _, t := range tracks.Tracks {
		seeds.Tracks = append(seeds.Tracks, t.ID)
	}

	rec, err := spotifyClient.GetRecommendations(c.Request().Context(), seeds, nil)
	if err != nil {
		slog.Error("could not get recommendations", "err", err)
		return err
	}
	slog.Info("recs", "r", rec)

	user, err := ar.App.Dao().FindAuthRecordByEmail("users", currentSpotifyUser.Email)
	if err != nil {
		user, err = models.CreateUser(ar.App.Dao(), currentSpotifyUser.Email, currentSpotifyUser.DisplayName)
		if err != nil {
			slog.Error("could not create user from spotify details", "err", err)
			return fmt.Errorf("Invalid credentials.")
		}
	}

	if err := ar.setAuthToken(c, user); err != nil {
		slog.Error("could not set auth token", "err", err)
		return err
	}

	oauth, err := models.GetOAuthByUserId(ar.App.Dao(), user.Id, "spotify")
	if err != nil {
		if err == models.ErrNoOAuthRows {
			slog.Info("no rows", "o", oauth)
		}
	}
	slog.Info("oauth update", "o", oauth)

	// expiry := types.DateTime{}
	// expiry.Scan(token.Expiry)
	//
	// userOAuth := models.OAuth{
	// 	Provider:     "spotify",
	// 	AccessToken:  token.AccessToken,
	// 	RefreshToken: token.RefreshToken,
	// 	TokenType:    token.TokenType,
	// 	Expiry:       expiry,
	// 	UserId:       user.Id,
	// }

	err = oauth.UpdateOAuth(ar.App.Dao(), token)
	if err != nil {
		return err
	}

	// err = oauth.Save(ar.App.Dao())
	// if err != nil {
	// 	return err
	// }

	return c.Redirect(http.StatusFound, "/")
}

const CookieOAuth = "oauth_state"

func newCookie(name, value string) *http.Cookie {
	cookie := http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		HttpOnly: true,
	}
	return &cookie
}
