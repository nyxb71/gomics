// Copyright (c) 2013-2018 Utkan Güngördü <utkan@freeconsole.org>
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

//go:generate go-bindata about.jpg icon.png gomics.glade

package main

import (
	"fmt"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"log"
	"reflect"
	"runtime"
	"time"
)

type GUI struct {
	MainWindow                     *gtk.Window            `build:"MainWindow"`
	VBox                           *gtk.Box               `build:"VBox"`
	Menubar                        *gtk.MenuBar           `build:"Menubar"`
	ScrolledWindow                 *gtk.ScrolledWindow    `build:"ScrolledWindow"`
	Viewport                       *gtk.Viewport          `build:"Viewport"`
	ImageBox                       *gtk.Box               `build:"ImageBox"`
	ImageL                         *gtk.Image             `build:"ImageL"`
	ImageR                         *gtk.Image             `build:"ImageR"`
	Statusbar                      *gtk.Statusbar         `build:"Statusbar"`
	AboutDialog                    *gtk.AboutDialog       `build:"AboutDialog"`
	MenuItemAbout                  *gtk.MenuItem          `build:"MenuItemAbout"`
	MenuItemOpen                   *gtk.MenuItem          `build:"MenuItemOpen"`
	MenuItemClose                  *gtk.MenuItem          `build:"MenuItemClose"`
	MenuItemQuit                   *gtk.MenuItem          `build:"MenuItemQuit"`
	MenuItemSaveImage              *gtk.MenuItem          `build:"MenuItemSaveImage"`
	FileChooserDialogArchive       *gtk.FileChooserDialog `build:"FileChooserDialogArchive"`
	Toolbar                        *gtk.Toolbar           `build:"Toolbar"`
	ButtonNextPage                 *gtk.ToolButton        `build:"ButtonNextPage"`
	ButtonPreviousPage             *gtk.ToolButton        `build:"ButtonPreviousPage"`
	ButtonLastPage                 *gtk.ToolButton        `build:"ButtonLastPage"`
	ButtonFirstPage                *gtk.ToolButton        `build:"ButtonFirstPage"`
	ButtonNextArchive              *gtk.ToolButton        `build:"ButtonNextArchive"`
	ButtonPreviousArchive          *gtk.ToolButton        `build:"ButtonPreviousArchive"`
	ButtonNextScene                *gtk.ToolButton        `build:"ButtonNextScene"`
	ButtonPreviousScene            *gtk.ToolButton        `build:"ButtonPreviousScene"`
	ButtonSkipForward              *gtk.ToolButton        `build:"ButtonSkipForward"`
	ButtonSkipBackward             *gtk.ToolButton        `build:"ButtonSkipBackward"`
	MenuItemNextPage               *gtk.MenuItem          `build:"MenuItemNextPage"`
	MenuItemPreviousPage           *gtk.MenuItem          `build:"MenuItemPreviousPage"`
	MenuItemLastPage               *gtk.MenuItem          `build:"MenuItemLastPage"`
	MenuItemFirstPage              *gtk.MenuItem          `build:"MenuItemFirstPage"`
	MenuItemNextArchive            *gtk.MenuItem          `build:"MenuItemNextArchive"`
	MenuItemPreviousArchive        *gtk.MenuItem          `build:"MenuItemPreviousArchive"`
	MenuItemSkipForward            *gtk.MenuItem          `build:"MenuItemSkipForward"`
	MenuItemSkipBackward           *gtk.MenuItem          `build:"MenuItemSkipBackward"`
	MenuItemEnlarge                *gtk.CheckMenuItem     `build:"MenuItemEnlarge"`
	MenuItemShrink                 *gtk.CheckMenuItem     `build:"MenuItemShrink"`
	MenuItemFullscreen             *gtk.CheckMenuItem     `build:"MenuItemFullscreen"`
	MenuItemSeamless               *gtk.CheckMenuItem     `build:"MenuItemSeamless"`
	MenuItemRandom                 *gtk.CheckMenuItem     `build:"MenuItemRandom"`
	MenuItemPreferences            *gtk.MenuItem          `build:"MenuItemPreferences"`
	MenuItemHFlip                  *gtk.CheckMenuItem     `build:"MenuItemHFlip"`
	MenuItemVFlip                  *gtk.CheckMenuItem     `build:"MenuItemVFlip"`
	MenuItemMangaMode              *gtk.CheckMenuItem     `build:"MenuItemMangaMode"`
	MenuItemDoublePage             *gtk.CheckMenuItem     `build:"MenuItemDoublePage"`
	MenuItemGoTo                   *gtk.MenuItem          `build:"MenuItemGoTo"`
	GoToThumbnailImage             *gtk.Image             `build:"GoToThumbnailImage"`
	MenuItemBestFit                *gtk.RadioMenuItem     `build:"MenuItemBestFit"`
	MenuItemOriginal               *gtk.RadioMenuItem     `build:"MenuItemOriginal"`
	MenuItemFitToWidth             *gtk.RadioMenuItem     `build:"MenuItemFitToWidth"`
	MenuItemFitToHeight            *gtk.RadioMenuItem     `build:"MenuItemFitToHeight"`
	PreferencesDialog              *gtk.Dialog            `build:"PreferencesDialog"`
	PagesToSkipSpinButton          *gtk.SpinButton        `build:"PagesToSkipSpinButton"`
	GoToDialog                     *gtk.Dialog            `build:"GoToDialog"`
	GoToSpinButton                 *gtk.SpinButton        `build:"GoToSpinButton"`
	GoToScrollbar                  *gtk.Scrollbar         `build:"GoToScrollbar"`
	InterpolationComboBoxText      *gtk.ComboBoxText      `build:"InterpolationComboBoxText"`
	OneWideCheckButton             *gtk.CheckButton       `build:"OneWideCheckButton"`
	SmartScrollCheckButton         *gtk.CheckButton       `build:"SmartScrollCheckButton"`
	EmbeddedOrientationCheckButton *gtk.CheckButton       `build:"EmbeddedOrientationCheckButton"`
	HideIdleCursorCheckButton      *gtk.CheckButton       `build:"HideIdleCursorCheckButton"`
	AddBookmarkMenuItem            *gtk.MenuItem          `build:"AddBookmarkMenuItem"`
	MenuBookmarks                  *gtk.Menu              `build:"MenuBookmarks"`
	RecentChooserMenu              *gtk.RecentChooserMenu `build:"RecentChooserMenu"`
	Config                         Config
	State                          State
	RecentManager                  *gtk.RecentManager
}

// LoadWidgets() fills the GUI struct with widgets built from the
// glade UI file at the specified location
func (gui *GUI) LoadWidgets() (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()

	builder, err := gtk.BuilderNew()
	if err != nil {
		return err
	}

	gomics_glade, err := Asset("gomics.glade")
	if err != nil {
		panic(err.Error())
	}
	if err = builder.AddFromString(string(gomics_glade)); err != nil {
		return err
	}

	guiStruct := reflect.ValueOf(gui).Elem()

	for i := 0; i < guiStruct.NumField(); i++ {
		field := guiStruct.Field(i)
		widget := guiStruct.Type().Field(i).Tag.Get("build")
		if widget == "" {
			continue
		}

		obj, err := builder.GetObject(widget)
		if err != nil {
			return err
		}

		w := reflect.ValueOf(obj).Convert(field.Type())
		field.Set(w)
	}

	return nil
}

func (gui *GUI) SetCursor(cursorName string) {
	disp, _ := gdk.DisplayGetDefault()
	win, _ := gui.Viewport.GetWindow()
	newCursor, _ := gdk.CursorNewFromName(disp, cursorName)
	win.SetCursor(newCursor)
}

func (gui *GUI) HideCursor() {
	// TODO Fix cursor not hiding
	gui.SetCursor("")
	gui.State.CursorHidden = true
}

func (gui *GUI) ShowCursor() {
	gui.SetCursor("default")
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

func (gui *GUI) initUI() {
	// Load UI
	if err := gui.LoadWidgets(); err != nil {
		log.Fatal(err)
	}

	about, err := Asset("about.jpg")
	if err != nil {
		panic(err.Error())
	}
	gui.AboutDialog.SetLogo(mustLoadPixbuf(about))
	icon, err := Asset("icon.png")
	gui.MainWindow.SetIcon(mustLoadPixbuf(icon))
	if err != nil {
		panic(err.Error())
	}

	if len(gitVersion) >= 7 {
		version := fmt.Sprintf("Version: %s (built: %s)\nCompiler version: %s", gitVersion[:7], buildDate, runtime.Version())
		gui.AboutDialog.SetVersion(version)
	}

	gui.FileChooserDialogArchive.AddButton("_Open", gtk.RESPONSE_ACCEPT)
	gui.FileChooserDialogArchive.AddButton("_Cancel", gtk.RESPONSE_CANCEL)

	gui.PreferencesDialog.AddButton("_OK", gtk.RESPONSE_ACCEPT)

	gui.GoToDialog.AddButton("_Cancel", gtk.RESPONSE_CANCEL)
	gui.GoToDialog.AddButton("_Go", gtk.RESPONSE_ACCEPT)
	//gui.GoToDialog.SetDefaultResponse(gtk.RESPONSE_ACCEPT)

	gui.syncUI()

	// Connect signals
	gui.MenuItemAbout.Connect("activate", gui.About)
	gui.MenuItemOpen.Connect("activate", gui.FileOpen)

	gui.MenuItemSaveImage.Connect("activate", gui.SavePNG)

	gui.MenuItemQuit.Connect("activate", gui.Quit)
	gui.MenuItemClose.Connect("activate", gui.FileClose)
	gui.MainWindow.Connect("delete-event", gui.Quit) // destroy

	var oldW, oldH int
	gui.MainWindow.Connect("size-allocate", func() {
		// Avoid unnecessary redraws
		w, h := gui.GetSize() // FIXME slow? use GdkRectangle *allocation passed in the signal
		if w == oldW && h == oldH {
			return
		}
		oldW, oldH = w, h
		gui.ResizeEvent()
	})

	gui.ButtonNextPage.Connect("clicked", gui.PageNext)
	gui.ButtonPreviousPage.Connect("clicked", gui.PagePrevious)
	gui.ButtonFirstPage.Connect("clicked", gui.PageFirst)
	gui.ButtonLastPage.Connect("clicked", gui.PageLast)
	gui.ButtonNextArchive.Connect("clicked", gui.ArchiveNext)
	gui.ButtonPreviousArchive.Connect("clicked", gui.ArchivePrevious)
	gui.ButtonNextScene.Connect("clicked", gui.NextScene)
	gui.ButtonPreviousScene.Connect("clicked", gui.PreviousScene)
	gui.ButtonSkipForward.Connect("clicked", gui.PageSkipForward)
	gui.ButtonSkipBackward.Connect("clicked", gui.PageSkipBack)

	gui.MenuItemNextPage.Connect("activate", gui.PageNext)
	gui.MenuItemPreviousPage.Connect("activate", gui.PagePrevious)
	gui.MenuItemFirstPage.Connect("activate", gui.PageFirst)
	gui.MenuItemLastPage.Connect("activate", gui.PageLast)
	gui.MenuItemNextArchive.Connect("activate", gui.ArchiveNext)
	gui.MenuItemPreviousArchive.Connect("activate", gui.ArchivePrevious)
	gui.MenuItemSkipForward.Connect("activate", gui.PageSkipForward)
	gui.MenuItemSkipBackward.Connect("activate", gui.PageSkipBack)

	gui.MenuItemEnlarge.Connect("toggled", gui.ImageEnlargeSmall)
	gui.MenuItemShrink.Connect("toggled", gui.ImageShrinkLarge)

	gui.MenuItemFullscreen.Connect("toggled", gui.ModeFullscreen)

	gui.MenuItemSeamless.Connect("toggled", gui.ModeSeamless)

	gui.MenuItemRandom.Connect("toggled", gui.OrderRandom)

	gui.MenuItemHFlip.Connect("toggled", gui.ImageFlipHorizontally)

	gui.MenuItemVFlip.Connect("toggled", gui.ImageFlipVertically)

	gui.MenuItemMangaMode.Connect("toggled", gui.ModeManga)

	gui.MenuItemDoublePage.Connect("toggled", gui.ModeDouble)

	gui.MenuItemOriginal.Connect("toggled", gui.ImageZoomOriginal)

	gui.MenuItemBestFit.Connect("toggled", gui.ImageZoomBestFit)

	gui.MenuItemFitToWidth.Connect("toggled", gui.ImageZoomFitWidth)

	gui.MenuItemFitToHeight.Connect("toggled", gui.ImageZoomFitHeight)

	gui.MenuItemPreferences.Connect("activate", gui.Preferences)

	gui.MenuItemGoTo.Connect("activate", gui.PageGoto)

	gui.GoToSpinButton.Connect("value-changed", func() {
		gui.GoToScrollbar.SetValue(gui.GoToSpinButton.GetValue())
		// TODO load & display the thumbnail image
	})

	gui.GoToScrollbar.Connect("value-changed", func() {
		gui.GoToSpinButton.SetValue(gui.GoToScrollbar.GetValue())
		gui.goToDialogLoadSetThumbnail()
		// load & display the thumbnail image
	})

	gui.RecentChooserMenu.Connect("item-activated", func() {
		uri := gui.RecentChooserMenu.GetCurrentUri()
		gui.LoadArchive(uri)
	})

	gui.PagesToSkipSpinButton.SetRange(1, 100)
	gui.PagesToSkipSpinButton.SetIncrements(1, 10)
	gui.PagesToSkipSpinButton.SetValue(float64(gui.Config.NSkip))

	gui.PagesToSkipSpinButton.Connect("value-changed", func() {
		gui.Config.NSkip = int(gui.PagesToSkipSpinButton.GetValue())
		gui.goToDialogLoadSetThumbnail()
	})

	gui.InterpolationComboBoxText.Connect("changed", func() {
		gui.SetInterpolation(gui.InterpolationComboBoxText.GetActive())
	})

	gui.OneWideCheckButton.Connect("toggled", func() {
		gui.SetOneWide(gui.OneWideCheckButton.GetActive())
	})

	gui.SmartScrollCheckButton.Connect("toggled", func() {
		gui.SetSmartScroll(gui.SmartScrollCheckButton.GetActive())
	})

	gui.EmbeddedOrientationCheckButton.Connect("toggled", func() {
		gui.SetEmbeddedOrientation(gui.EmbeddedOrientationCheckButton.GetActive())
	})

	gui.HideIdleCursorCheckButton.Connect("toggled", func() {
		gui.SetHideIdleCursor(gui.HideIdleCursorCheckButton.GetActive())

	})

	gui.AddBookmarkMenuItem.Connect("activate", func() {
		gui.AddBookmark()
	})

	gui.ScrolledWindow.SetEvents(gui.ScrolledWindow.GetEvents() | int(gdk.BUTTON_PRESS_MASK))

	gui.ScrolledWindow.Connect("scroll-event", func(w *gtk.ScrolledWindow, e *gdk.Event) {
		se := &gdk.EventScroll{e}

		gui.Scroll(se.DeltaX(), se.DeltaY())
	})

	// FIXME
	gui.ScrolledWindow.Connect("button-press-event", func(_ *gtk.ScrolledWindow, e *gdk.Event) bool {
		//log.Println(w)
		be := &gdk.EventButton{e}
		switch be.Button() {
		case 1:
			gui.PageNext()
		case 3:
			gui.PagePrevious()
		case 2:
			gui.ArchiveNext()
		}
		return true
	})

	gui.MainWindow.Connect("motion-notify-event", func(_ *gtk.Window, _ *gdk.Event) bool {
		gui.State.CursorLastMoved = time.Now()
		return true
	})

	glib.TimeoutAdd(250, gui.UpdateCursorVisibility)

	gui.MainWindow.Connect("key-press-event", func(_ *gtk.Window, e *gdk.Event) {
		kp := KeyPress{gdk.EventKey{e}}

		switch kp.Key() {
		case gdk.KEY_Down:
			if kp.Ctrl() {
				gui.ArchiveNext()
			} else if kp.Shift() {
				gui.Scroll(0, 1)
			} else {
				gui.PageNext()
			}
		case gdk.KEY_Up:
			if kp.Ctrl() {
				gui.ArchivePrevious()
			} else if kp.Shift() {
				gui.Scroll(0, -1)
			} else {
				gui.PagePrevious()
			}
		case gdk.KEY_Right:
			if kp.Ctrl() {
				gui.NextScene()
			} else if kp.Shift() {
				gui.Scroll(1, 0)
			} else {
				gui.PageSkipForward()
			}
		case gdk.KEY_Left:
			if kp.Ctrl() {
				gui.PreviousScene()
			} else if kp.Shift() {
				gui.Scroll(-1, 0)
			} else {
				gui.PageSkipBack()
			}
		}
	})

	gui.RebuildBookmarksMenu()

	gui.MainWindow.SetDefaultSize(gui.Config.WindowWidth, gui.Config.WindowHeight)
	gui.MainWindow.ShowAll()

	// Tiny hack
	mw, mh := gui.MainWindow.GetSize()
	va := gui.Viewport.GetAllocation()
	gui.State.DeltaW, gui.State.DeltaH = mw-va.GetWidth(), mh-va.GetHeight()

	gui.SetFullscreen(gui.Config.Fullscreen)
	gui.SetZoomMode(gui.Config.ZoomMode)
	gui.SetDoublePage(gui.Config.DoublePage)
	gui.SetMangaMode(gui.Config.MangaMode)

	gui.fixFocus()
}

func (gui *GUI) goToDialogLoadSetThumbnail() {
	n := int(gui.GoToSpinButton.GetValue() - 1)
	pixbuf, err := gui.State.Archive.Load(n, gui.Config.EmbeddedOrientation)
	if err != nil {
		gui.ShowError(err.Error())
		return
	}

	w, h := fit(pixbuf.GetWidth(), pixbuf.GetHeight(), 128, 128)

	scaled, err := pixbuf.ScaleSimple(w, h, interpolations[gui.Config.Interpolation])
	if err != nil {
		gui.ShowError(err.Error())
		return
	}

	gui.State.GoToThumnailPixbuf = scaled
	gui.GoToThumbnailImage.SetFromPixbuf(scaled)

	gc()
}

func (gui *GUI) syncUI() {
	// Sync config & UI
	gui.MenuItemEnlarge.SetActive(gui.Config.Enlarge)
	gui.MenuItemShrink.SetActive(gui.Config.Shrink)
	gui.MenuItemHFlip.SetActive(gui.Config.HFlip)
	gui.MenuItemVFlip.SetActive(gui.Config.VFlip)
	gui.MenuItemRandom.SetActive(gui.Config.Random)
	gui.MenuItemSeamless.SetActive(gui.Config.Seamless)
	gui.MenuItemDoublePage.SetActive(gui.Config.DoublePage)
	gui.MenuItemMangaMode.SetActive(gui.Config.MangaMode)

	switch gui.Config.ZoomMode {
	case "FitToWidth":
		gui.MenuItemFitToWidth.SetActive(true)
	case "FitToHeight":
		gui.MenuItemFitToHeight.SetActive(true)
	case "BestFit":
		gui.MenuItemBestFit.SetActive(true)
	default:
		gui.MenuItemOriginal.SetActive(true)
	}

	gui.InterpolationComboBoxText.SetActive(gui.Config.Interpolation)
	gui.OneWideCheckButton.SetActive(gui.Config.OneWide)
	gui.EmbeddedOrientationCheckButton.SetActive(gui.Config.EmbeddedOrientation)
	gui.HideIdleCursorCheckButton.SetActive(gui.Config.HideIdleCursor)
}

func (gui *GUI) RunGoToDialog() {
	if !gui.Loaded() {
		return
	}

	gui.GoToSpinButton.SetRange(1, float64(gui.State.Archive.Len()))
	gui.GoToSpinButton.SetValue(float64(gui.State.ArchivePos) + 1)
	gui.GoToSpinButton.SetIncrements(1, float64(gui.Config.NSkip))

	gui.GoToScrollbar.SetRange(1, float64(gui.State.Archive.Len()))
	gui.GoToScrollbar.SetValue(float64(gui.State.ArchivePos) + 1)
	gui.GoToScrollbar.SetIncrements(1, float64(gui.State.Archive.Len()))

	gui.goToDialogLoadSetThumbnail()

	res := gtk.ResponseType(gui.GoToDialog.Run())
	gui.GoToDialog.Hide()
	if res == gtk.RESPONSE_ACCEPT {
		gui.SetPage(int(gui.GoToSpinButton.GetValue()) - 1)

		gui.GoToThumbnailImage.Clear()
		gui.State.GoToThumnailPixbuf = nil
		gc()
	}
}
