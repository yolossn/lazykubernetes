package main

import (
	"github.com/yolossn/lazykubernetes/pkg/app"
	"github.com/yolossn/lazykubernetes/pkg/client"
)

func main() {
	// Setup k8sClient
	k8sClient, err := client.Newk8s()
	if err != nil {
		panic(err)
	}
	// _, _ = k8sClient.GetServerInfo()
	ui, err := app.NewApp(k8sClient)
	if err != nil {
		panic(err)
	}

	err = ui.Run()
	if err != nil {
		panic(err)
	}
}
