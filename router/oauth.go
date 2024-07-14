package router

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/dhinogz/spotify-test/models"
	"github.com/gorilla/csrf"
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/tools/types"
	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
)

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

	expiry := types.DateTime{}
	expiry.Scan(token.Expiry)

	userOAuth := models.OAuth{
		Provider:     "spotify",
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		TokenType:    token.TokenType,
		Expiry:       expiry,
		UserId:       user.Id,
	}

	err = userOAuth.Save(ar.App.Dao())
	if err != nil {
		return err
	}

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
