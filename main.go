package main

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"

	"github.com/dhinogz/discover-friendly/internal/handlers"
	_ "github.com/dhinogz/discover-friendly/migrations"
)

func main() {
	app := pocketbase.New()
	isGoRun := strings.HasPrefix(os.Args[0], os.TempDir())

	migratecmd.MustRegister(app, app.RootCmd, migratecmd.Config{
		Automigrate: isGoRun,
	})

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		currentEnv := os.Getenv("ENV")
		if currentEnv == "dev" {
			if err := godotenv.Load(); err != nil {
				e.App.Logger().Error("could not load env", "err", err)
				return err
			}
		}

		return handlers.NewAppRouter(e).SetupRoutes(isGoRun)
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
