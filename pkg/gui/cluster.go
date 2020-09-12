package gui

import (
	"fmt"

	"github.com/jesseduffield/gocui"
	"github.com/yolossn/lazykubernetes/pkg/utils"
)

func (gui *Gui) getClusterInfoView() *gocui.View {
	v, _ := gui.g.View("cluster-info")
	return v
}

func (gui *Gui) onClusterInfoClick(g *gocui.Gui, v *gocui.View) error {

	if _, err := g.SetCurrentView(v.Name()); err != nil {
		return err
	}

	infoView := gui.getInfoView()
	out := utils.GetLazykubernetesArt()
	fmt.Fprintln(infoView, out)
	return nil
}

func (gui *Gui) reRenderClusterInfo() error {

	clusterView := gui.getClusterInfoView()

	info, err := gui.k8sClient.GetServerInfo()
	if err != nil {
		return nil
	}

	clusterView.Clear()
	fmt.Fprintf(clusterView, "Version:   %s.%s\n", info.Major, info.Minor)
	fmt.Fprintf(clusterView, "platform:   %s\n", info.Platform)

	return nil
}
