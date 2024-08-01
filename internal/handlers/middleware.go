package handlers

import (
	"time"

	"github.com/dhinogz/discover-friendly/internal/models"
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/apis"
	pbmodels "github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tokens"
	"github.com/pocketbase/pocketbase/tools/security"
	"github.com/spf13/cast"
	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2"
)

const AuthCookieName = "Auth"

func (ar *AppRouter) LoadAuthContextFromCookie() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			tokenCookie, err := c.Request().Cookie(AuthCookieName)
			if err != nil || tokenCookie.Value == "" {
				return next(c) // no token cookie
			}

			token := tokenCookie.Value

			claims, _ := security.ParseUnverifiedJWT(token)
			tokenType := cast.ToString(claims["type"])

			switch tokenType {
			case tokens.TypeAdmin:
				admin, err := ar.App.Dao().FindAdminByToken(
					token,
					ar.App.Settings().AdminAuthToken.Secret,
				)
				if err == nil && admin != nil {
					c.Set(apis.ContextAdminKey, admin)
				}

			case tokens.TypeAuthRecord:
				record, err := ar.App.Dao().FindAuthRecordByToken(
					token,
					ar.App.Settings().RecordAuthToken.Secret,
				)
				if err == nil && record != nil {
					c.Set(apis.ContextAuthRecordKey, record)
				}
			}

			return next(c)
		}
	}
}

// TODO: middleware to check for oauth token expiry. If expired, refresh it
func (ar *AppRouter) LoadSpotifyAuthMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			rec := c.Get(apis.ContextAuthRecordKey)
			if rec == nil {
				return next(c)
			}
			user := rec.(*pbmodels.Record)

			o, err := models.GetOAuthByUserId(ar.App.Dao(), user.Id, "spotify")
			if err != nil {
				return err
			}

			token := &oauth2.Token{
				AccessToken:  o.AccessToken,
				TokenType:    o.TokenType,
				RefreshToken: o.RefreshToken,
				Expiry:       o.Expiry.Time(),
			}

			// Create a new Spotify client
			config := spotifyauth.New()
			client := config.Client(c.Request().Context(), token)
			spotifyClient := spotify.New(client)

			// Check if token needs refresh
			if token.Expiry.Before(time.Now()) {
				newToken, err := config.RefreshToken(c.Request().Context(), token)
				if err != nil {
					return err
				}

				if err := o.UpdateOAuth(ar.App.Dao(), newToken); err != nil {
					return err
				}

				// Update the Spotify client with the new token
				client = config.Client(c.Request().Context(), newToken)
				spotifyClient = spotify.New(client)
			}

			// Set the Spotify client in the context
			c.Set("spotifyClient", spotifyClient)

			return next(c)
		}
	}
}
