package gui

import "github.com/jesseduffield/gocui"


func (gui *Gui) getInfoView() *gocui.View {
	v, _ := gui.g.View("info")
	return v
}

func (gui *Gui) onInfoClick(g *gocui.Gui, v *gocui.View) error {
	if _, err := g.SetCurrentView(v.Name()); err != nil {
		return err
	}
	return nil
}

func (gui *Gui) onInfoTabCick(tabIndex int) error {
	resourceView := gui.getInfoView()
	resourceView.TabIndex = tabIndex

	return nil
}

