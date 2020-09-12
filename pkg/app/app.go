package app

import (
	"github.com/yolossn/lazykubernetes/pkg/gui"
)

type App struct {
	Gui *gui.Gui
}

func (app *App) Run() error {
	return app.Gui.Run()
}

func NewApp() (*App, error) {

	gui, err := gui.NewGui()
	if err != nil {
		return nil, err
	}

	return &App{Gui: gui}, nil
}
