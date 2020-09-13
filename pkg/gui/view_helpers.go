package gui

import (
	"time"

	"github.com/jesseduffield/gocui"
)

// This function is a modified version of
// https://github.com/jesseduffield/lazydocker/blob/fa6460b8ab3486b7e84c3a7d4c64fbd8e3f4be21/pkg/gui/gui.go#L227
func (gui *Gui) goEvery(interval time.Duration, function func() error) {
	_ = function() // time.Tick doesn't run immediately so we'll do that here // TODO: maybe change
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for range ticker.C {
			_ = function()
		}
	}()
}

// This function is a modified version of
// https://github.com/jesseduffield/lazydocker/blob/a14e6400cbbd7e2aa9ef22166b085b0678f9ca3a/pkg/gui/view_helpers.go#L361
func (gui *Gui) FindSelectedLine(v *gocui.View, itemCount int) int {
	_, cy := v.Cursor()
	_, oy := v.Origin()

	selectedLine := cy - oy

	if selectedLine < 0 {
		return 0
	}

	if selectedLine > itemCount-1 {
		return itemCount - 1
	}
	return selectedLine
}

// This function was copied from
// https://github.com/jesseduffield/lazydocker/blob/a14e6400cbbd7e2aa9ef22166b085b0678f9ca3a/pkg/gui/view_helpers.go#L319
func (gui *Gui) changeSelectedLine(line *int, total int, up bool) {
	if up {
		if *line == -1 || *line == 0 {
			return
		}

		*line -= 1
	} else {
		if *line == -1 || *line == total-1 {
			return
		}

		*line += 1
	}
}

// This function was copied from
// https://github.com/jesseduffield/lazydocker/blob/a14e6400cbbd7e2aa9ef22166b085b0678f9ca3a/pkg/gui/view_helpers.go#L168
func (gui *Gui) focusPoint(selectedX int, selectedY int, lineCount int, v *gocui.View) error {
	if selectedY < 0 || selectedY > lineCount {
		return nil
	}
	ox, oy := v.Origin()
	originalOy := oy
	cx, cy := v.Cursor()
	originalCy := cy
	_, height := v.Size()

	ly := Max(height-1, 0)

	windowStart := oy
	windowEnd := oy + ly

	if selectedY < windowStart {
		oy = Max(oy-(windowStart-selectedY), 0)
	} else if selectedY > windowEnd {
		oy += (selectedY - windowEnd)
	}

	if windowEnd > lineCount-1 {
		shiftAmount := (windowEnd - (lineCount - 1))
		oy = Max(oy-shiftAmount, 0)
	}

	if originalOy != oy {
		_ = v.SetOrigin(ox, oy)
	}

	cy = selectedY - oy
	if originalCy != cy {
		_ = v.SetCursor(cx, selectedY-oy)
	}
	return nil
}

func Max(x, y int) int {
	if x > y {
		return x
	}
	return y
}
