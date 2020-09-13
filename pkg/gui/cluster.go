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
	infoView.Clear()
	out := utils.GetLazykubernetesArt()
	fmt.Fprintln(infoView, out)
	return nil
}

func (gui *Gui) reRenderClusterInfo() error {

	clusterView := gui.getClusterInfoView()

	info, err := gui.k8sClient.GetServerInfo()
	if err != nil {
		clusterView.Clear()
		fmt.Fprintf(clusterView, "Health: %s", "🔴")
		return nil
	}

	clusterView.Clear()
	fmt.Fprintf(clusterView, "Version:   %s.%s\nplatform:  %s\nHealth:  %s", info.Major, info.Minor, info.Platform, "🟢")
	return nil
}
