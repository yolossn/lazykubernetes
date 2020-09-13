package main

import (
	"log"
	"os"

	"github.com/jesseduffield/gocui"
	"github.com/yolossn/lazykubernetes/pkg/app"
	"github.com/yolossn/lazykubernetes/pkg/client"
)

func main() {
	// Setup k8sClient
	k8sClient, err := client.Newk8s()
	if err != nil {
		log.Fatal("Couldn't connect to the k8s cluster")
	}

	ui, err := app.NewApp(k8sClient)
	if err != nil {
		log.Fatal("Something went wrong")
	}

	err = ui.Run()
	if err != nil {
		if err == gocui.ErrQuit {
			os.Exit(0)
		}
		log.Fatal("Something went wrong")
	}
}
