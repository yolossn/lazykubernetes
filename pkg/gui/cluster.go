package gui

import "github.com/jesseduffield/gocui"

func (gui *Gui) onClusterInfoClick(g *gocui.Gui, v *gocui.View) error {

	if _, err := g.SetCurrentView(v.Name()); err != nil {
		return err
	}

	infoView := gui.getInfoView()
	infoView.Tabs = []string{"Info"}
	return nil
}
