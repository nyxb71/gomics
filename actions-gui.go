package main

import (
	"github.com/gotk3/gotk3/gdk"
	"log"
	"time"
)


func (gui *GUI) SetCursor(cursorName string) (err error) {
	disp, err := gdk.DisplayGetDefault()
	if err != nil {
		log.Println("Error setting cursor, could not get default display.")
		return err
	}

	win, err := gui.Viewport.GetWindow()
	if err != nil {
		log.Println("Error setting cursor, could not get viewport window.")
		return err
	}

	newCursor, err := gdk.CursorNewFromName(disp, cursorName)
	if err != nil {
		log.Println("Error setting cursor, could not get cursor: ", cursorName)
		return err
	}

	win.SetCursor(newCursor)
	return nil
}

func (gui *GUI) HideCursor() {
	if err := gui.SetCursor("none"); err != nil {
		log.Print("Error hiding cursor")
		return
	}

	gui.State.CursorHidden = true
}

func (gui *GUI) ShowCursor() {
	if err := gui.SetCursor("default"); err != nil {
		log.Print("Error showing cursor")
		return
	}

	gui.State.CursorHidden = false
}

func (gui *GUI) UpdateCursorVisibility() bool {
	cursorShouldBeHidden := false

	if gui.Config.HideIdleCursor && !gui.State.CursorForceShown {
		cursorShouldBeHidden = time.Since(gui.State.CursorLastMoved).Seconds() > 1
	}

	if cursorShouldBeHidden && !gui.State.CursorHidden {
		gui.HideCursor()
	}

	if !cursorShouldBeHidden && gui.State.CursorHidden {
		gui.ShowCursor()
	}

	return true
}
