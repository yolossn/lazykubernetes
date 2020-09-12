package gui

import (
	"fmt"

	"github.com/jesseduffield/gocui"
)

func (gui *Gui) getNamespaceView() *gocui.View {
	v, _ := gui.g.View("namespace")
	return v
}

func (gui *Gui) onNamespaceClick(g *gocui.Gui, v *gocui.View) error {

	if _, err := g.SetCurrentView(v.Name()); err != nil {
		return err
	}

	infoView := gui.getInfoView()
	infoView.Tabs = getNamespaceInfoTabs()

	return nil
}

func (gui *Gui) reRenderNamespace() error {
	nsView := gui.getNamespaceView()
	if nsView == nil {
		return nil
	}

	ns, err := gui.k8sClient.ListNamespace()
	if err != nil {
		return err
	}

	gui.g.Update(func(*gocui.Gui) error {
		nsView.Clear()
		for _, n := range ns {
			fmt.Fprintln(nsView, n.Name)
			fmt.Println(n.Name)
		}
		return nil
	})

	return nil
}
