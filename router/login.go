package router

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/apis"

	"github.com/dhinogz/spotify-test/components"
	"github.com/dhinogz/spotify-test/components/auth"
	"github.com/dhinogz/spotify-test/components/shared"
	"github.com/dhinogz/spotify-test/htmx"
)

// GetLogin returns the login page.
func (ar *AppRouter) GetLogin(c echo.Context) error {
	if c.Get(apis.ContextAuthRecordKey) != nil {
		return c.Redirect(302, "/")
	}

	return components.Render(c, http.StatusOK, auth.LoginPage(shared.Context{}))
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

	return htmx.Redirect(c, "/")
}
