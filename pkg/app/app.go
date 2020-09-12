package app

import (
	"github.com/yolossn/lazykubernetes/pkg/client"
	"github.com/yolossn/lazykubernetes/pkg/gui"
)

type App struct {
	Gui *gui.Gui
}

func (app *App) Run() error {
	return app.Gui.Run()
}

func NewApp(k8sClient *client.K8s) (*App, error) {

	gui, err := gui.NewGui(k8sClient)
	if err != nil {
		return nil, err
	}

	return &App{Gui: gui}, nil
}
