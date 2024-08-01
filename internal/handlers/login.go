package handlers

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/apis"

	"github.com/dhinogz/discover-friendly/internal/components"
	"github.com/dhinogz/discover-friendly/internal/components/auth"
	"github.com/dhinogz/discover-friendly/internal/components/shared"
	"github.com/dhinogz/discover-friendly/internal/htmx"
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
