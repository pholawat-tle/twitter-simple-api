package main

import (
	"net/http"
	transportHTTP "twitter/internal/transport/http"

	log "github.com/sirupsen/logrus"
)

type App struct {
	Name    string
	Version string
}

func NewApp(name string, version string) *App {
	return &App{Name: name, Version: version}
}

func (app *App) Run() error {
	log.Info("Starting the application")
	h := transportHTTP.NewHandler()
	h.SetUpRoutes()

	if err := http.ListenAndServe(":8081", h.Router); err != nil {
		return err
	}

	return nil
}

func main() {
	app := NewApp("Twitter Clone", "1.0.0")

	if err := app.Run(); err != nil {
		log.Error(err)
		log.Fatal("The application can't be started!")
	}
}
