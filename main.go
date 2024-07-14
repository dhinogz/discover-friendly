package main

import (
	"log"
	"os"
	"strings"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"

	_ "github.com/dhinogz/spotify-test/migrations"
	"github.com/dhinogz/spotify-test/router"
)

func main() {
	app := pocketbase.New()
	isGoRun := strings.HasPrefix(os.Args[0], os.TempDir())

	migratecmd.MustRegister(app, app.RootCmd, migratecmd.Config{
		Automigrate: isGoRun,
	})

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		// if err := godotenv.Load(); err != nil {
		// 	slog.Error("could not load env", "err", err)
		// 	return err
		// }

		return router.NewAppRouter(e).SetupRoutes(isGoRun)
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
