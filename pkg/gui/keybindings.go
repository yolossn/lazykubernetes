package gui

import "github.com/jesseduffield/gocui"

func (gui *Gui) SetKeybindings(g *gocui.Gui) error {

	// Exit Keybinding Ctrl+C
	if err := g.SetKeybinding("", nil, gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}

	// MouseClick
	if err := g.SetKeybinding("cluster-info", nil, gocui.MouseLeft, gocui.ModNone, gui.onClusterInfoClick); err != nil {
		return err
	}

	if err := g.SetKeybinding("namespace", nil, gocui.MouseLeft, gocui.ModNone, gui.onNamespaceClick); err != nil {
		return err
	}

	if err := g.SetKeybinding("resource", nil, gocui.MouseLeft, gocui.ModNone, gui.onResourceClick); err != nil {
		return err
	}

	if err := g.SetKeybinding("info", nil, gocui.MouseLeft, gocui.ModNone, gui.onInfoClick); err != nil {
		return err
	}

	// Arrow keys
	if err := g.SetKeybinding("namespace", nil, gocui.KeyArrowDown, gocui.ModNone, gui.handleNSKeyUp); err != nil {
		return err
	}

	if err := g.SetKeybinding("namespace", nil, gocui.KeyArrowUp, gocui.ModNone, gui.handleNSKeyDown); err != nil {
		return err
	}

	if err := g.SetKeybinding("resource", nil, gocui.KeyArrowDown, gocui.ModNone, gui.handleResourceKeyUp); err != nil {
		return err
	}

	if err := g.SetKeybinding("resource", nil, gocui.KeyArrowUp, gocui.ModNone, gui.handleResourceKeyDown); err != nil {
		return err
	}

	// Tab click
	if err := g.SetTabClickBinding("resource", gui.onResourceTabClick); err != nil {
		return err
	}

	if err := g.SetTabClickBinding("info", gui.onInfoTabCick); err != nil {
		return err
	}

	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
