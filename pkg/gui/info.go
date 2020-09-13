package gui

import (
	"math"

	"github.com/jesseduffield/gocui"
)

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

func (gui *Gui) onInfoTabClick(tabIndex int) error {

	infoView := gui.getInfoView()
	infoView.TabIndex = tabIndex
	gui.panelStates.Info.TabIndex = tabIndex

	return nil
}

// The following scroll functions are modified version of code from lazydocker
// https://github.com/jesseduffield/lazydocker/blob/fa6460b8ab3486b7e84c3a7d4c64fbd8e3f4be21/pkg/gui/main_panel.go
func (gui *Gui) scrollLeftInfo(g *gocui.Gui, v *gocui.View) error {
	infoView := gui.getInfoView()
	ox, oy := infoView.Origin()
	newOx := int(math.Max(0, float64(ox-20)))

	return infoView.SetOrigin(newOx, oy)
}

func (gui *Gui) scrollRightInfo(g *gocui.Gui, v *gocui.View) error {
	infoView := gui.getInfoView()
	ox, oy := infoView.Origin()

	content := infoView.ViewBufferLines()
	var largestNumberOfCharacters int
	for _, txt := range content {
		if len(txt) > largestNumberOfCharacters {
			largestNumberOfCharacters = len(txt)
		}
	}

	sizeX, _ := infoView.Size()
	if ox+sizeX >= largestNumberOfCharacters {
		return nil
	}

	return infoView.SetOrigin(ox+20, oy)
}

func (gui *Gui) scrollUpInfo(g *gocui.Gui, v *gocui.View) error {
	mainView := gui.getInfoView()
	mainView.Autoscroll = false
	ox, oy := mainView.Origin()
	newOy := int(math.Max(0, float64(oy-20)))
	return mainView.SetOrigin(ox, newOy)
}

func (gui *Gui) scrollDownInfo(g *gocui.Gui, v *gocui.View) error {
	mainView := gui.getInfoView()
	mainView.Autoscroll = false
	ox, oy := mainView.Origin()

	reservedLines := 0
	_, sizeY := mainView.Size()
	reservedLines = sizeY

	totalLines := mainView.ViewLinesHeight()
	if oy+reservedLines >= totalLines {
		return nil
	}

	return mainView.SetOrigin(ox, oy+20)
}
