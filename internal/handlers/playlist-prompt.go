package handlers

import (
	"net/http"

	"github.com/dhinogz/discover-friendly/internal/ui"
	"github.com/dhinogz/discover-friendly/internal/ui/shared"
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/apis"
	pbmodels "github.com/pocketbase/pocketbase/models"
)

func (ar *AppRouter) HandleNewPlaylistPrompt(c echo.Context) error {
	// TODO: this should be middlware accessed by the context
	rec := c.Get(apis.ContextAuthRecordKey)
	if rec == nil {
		return ui.Render(c, http.StatusOK, ui.Home(shared.Context{}))
	}

	user := c.Get(apis.ContextAuthRecordKey).(*pbmodels.Record)

	return ui.Render(c, http.StatusOK, ui.NewPlaylistPrompt(shared.Context{User: user}))
}
