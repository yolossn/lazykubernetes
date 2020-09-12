package main

import "github.com/yolossn/lazykubernetes/pkg/app"

func main() {
	ui, err := app.NewApp()
	if err != nil {
		panic(err)
	}

	err = ui.Run()
	if err != nil {
		panic(err)
	}
}
