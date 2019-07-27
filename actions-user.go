package main

import (
	"errors"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
	"log"
	"path/filepath"
	"time"
)

type UserActionBinding struct {
	Keybind KeyChord
	Call    func()
}

type UserAction int

const (
	Unknown UserAction = iota

	FileOpen
	FileClose
	FileSaveImage
	Quit

	Preferences
	About

	ImageShrinkLarge
	ImageEnlargeSmall
	ImageZoomBestFit
	ImageZoomOriginal
	ImageZoomFitWidth
	ImageZoomFitHeight
	ImageFlipVertically
	ImageFlipHorizontally

	ModeFullscreen
	ModeSeamless
	ModeManga
	ModeDouble

	OrderRandom

	PageNext
	PagePrevious
	PageSkipForward
	PageSkipBackward
	PageFirst
	PageLast
	PageGoto
	ArchiveNext
	ArchivePrevious

	BookmarkAdd
)

type UserActions map[UserAction]UserActionBinding

func (ua UserActions) Init(gui *GUI) {
	if ua == nil {
		log.Fatal("Cannot initialize nil UserActions map")
	}

	ua[FileOpen] = UserActionBinding{KeyChord{Key: gdk.KEY_o, Ctrl: true}, gui.FileOpen}
	ua[FileClose] = UserActionBinding{KeyChord{Key: gdk.KEY_w, Ctrl: true}, gui.FileClose}
	ua[FileSaveImage] = UserActionBinding{KeyChord{Key: gdk.KEY_s, Ctrl: true}, gui.FileSaveImage}
	ua[Quit] = UserActionBinding{KeyChord{Key: gdk.KEY_q, Ctrl: true}, gui.Quit}

	ua[Preferences] = UserActionBinding{KeyChord{Key: gdk.KEY_p, Ctrl: true}, gui.Quit}
	ua[About] = UserActionBinding{KeyChord{Key: gdk.KEY_F1}, gui.About}

	ua[ImageShrinkLarge] = UserActionBinding{KeyChord{Key: gdk.KEY_s}, gui.ImageShrinkLarge}
	ua[ImageEnlargeSmall] = UserActionBinding{KeyChord{Key: gdk.KEY_e}, gui.ImageEnlargeSmall}
	ua[ImageZoomBestFit] = UserActionBinding{KeyChord{Key: gdk.KEY_b}, gui.ImageZoomBestFit}
	ua[ImageZoomOriginal] = UserActionBinding{KeyChord{Key: gdk.KEY_o}, gui.ImageZoomOriginal}
	ua[ImageZoomFitWidth] = UserActionBinding{KeyChord{Key: gdk.KEY_w}, gui.ImageZoomFitWidth}
	ua[ImageZoomFitHeight] = UserActionBinding{KeyChord{Key: gdk.KEY_h}, gui.ImageZoomFitHeight}
	ua[ImageFlipVertically] = UserActionBinding{KeyChord{Key: gdk.KEY_v}, gui.ImageFlipVertically}
	ua[ImageFlipHorizontally] = UserActionBinding{KeyChord{Key: gdk.KEY_v, Shift: true}, gui.ImageFlipHorizontally}

	ua[OrderRandom] = UserActionBinding{KeyChord{Key: gdk.KEY_r}, gui.OrderRandom}

	ua[ModeFullscreen] = UserActionBinding{KeyChord{Key: gdk.KEY_f}, gui.ModeFullscreen}
	ua[ModeSeamless] = UserActionBinding{KeyChord{Key: gdk.KEY_s, Shift: true}, gui.ModeSeamless}
	ua[ModeManga] = UserActionBinding{KeyChord{Key: gdk.KEY_m}, gui.ModeManga}
	ua[ModeDouble] = UserActionBinding{KeyChord{Key: gdk.KEY_d}, gui.ModeDouble}

	ua[PageNext] = UserActionBinding{KeyChord{Key: gdk.KEY_Down}, gui.PageNext}
	ua[PagePrevious] = UserActionBinding{KeyChord{Key: gdk.KEY_Up}, gui.PagePrevious}
	ua[PageSkipForward] = UserActionBinding{KeyChord{Key: gdk.KEY_Right}, gui.PageSkipForward}
	ua[PageSkipBackward] = UserActionBinding{KeyChord{Key: gdk.KEY_Left}, gui.PageSkipBackward}
	ua[PageFirst] = UserActionBinding{KeyChord{Key: gdk.KEY_Home}, gui.PageFirst}
	ua[PageLast] = UserActionBinding{KeyChord{Key: gdk.KEY_End}, gui.PageLast}
	ua[ArchiveNext] = UserActionBinding{KeyChord{Key: gdk.KEY_Page_Down, Ctrl: true}, gui.ArchiveNext}
	ua[ArchivePrevious] = UserActionBinding{KeyChord{Key: gdk.KEY_Page_Up, Ctrl: true}, gui.ArchivePrevious}
}

func (ua UserActions) GetUserActionByKeyChord(kch KeyChord) (UserActionBinding, error) {
	// Returns first occurence, since there should only be one action per keybind.

	for _, uab := range ua {
		if uab.Keybind == kch {
			return uab, nil
		}
	}

	return UserActionBinding{}, errors.New("No user actions with that keybind.")
}

func MakeUserActions(gui *GUI) UserActions {
	ua := make(UserActions)
	ua.Init(gui)
	if ua == nil {
		log.Fatal("nil UserActions map")
	}
	return ua
}

func (gui *GUI) FileOpen() {
	res := gtk.ResponseType(gui.FileChooserDialogArchive.Run())
	gui.FileChooserDialogArchive.Hide()
	if res == gtk.RESPONSE_ACCEPT {
		filename := gui.FileChooserDialogArchive.GetFilename()
		gui.LoadArchive(filename)
	}
}

func (gui *GUI) FileClose() {
	if !gui.Loaded() {
		return
	}

	gui.State.Archive.Close()

	gui.State.Archive = nil
	gui.State.ArchiveName = ""
	gui.State.ArchivePath = ""
	gui.State.ArchivePos = 0

	gui.State.ImageHash = nil

	gui.ImageL.Clear()
	gui.ImageR.Clear()
	gui.State.PixbufL = nil
	gui.State.PixbufR = nil
	gui.State.CursorLastMoved = time.Now()
	gui.State.CursorHidden = false
	gui.State.CursorForceShown = false
	gui.SetStatus("")
	gui.MainWindow.SetTitle("Gomics")
	gc()
}

func (gui *GUI) FileSaveImage() {

}

func (gui *GUI) Quit() {
	gui.Config.WindowWidth, gui.Config.WindowHeight = gui.MainWindow.GetSize()

	if err := gui.Config.Save(filepath.Join(gui.State.ConfigPath, ConfigFile)); err != nil {
		log.Println(err)
	}
	gtk.MainQuit()
}

func (gui *GUI) Preferences() {
	gui.State.CursorForceShown = true
	res := gtk.ResponseType(gui.PreferencesDialog.Run())
	gui.PreferencesDialog.Hide()
	if res == gtk.RESPONSE_ACCEPT {
		// TODO save config
	}
	gui.State.CursorForceShown = false
}

func (gui *GUI) About() {
	gui.State.CursorForceShown = true
	gui.AboutDialog.Run()
	gui.AboutDialog.Hide()
	gui.State.CursorForceShown = false
}

func (gui *GUI) ImageShrinkLarge() {
	gui.SetShrink(!gui.Config.Shrink)
}

func (gui *GUI) ImageEnlargeSmall() {
	gui.SetEnlarge(!gui.Config.Enlarge)
}

func (gui *GUI) ImageZoomBestFit() {
	if gui.MenuItemBestFit.GetActive() {
		gui.SetZoomMode("BestFit")
	}
}

func (gui *GUI) ImageZoomOriginal() {
	if gui.MenuItemOriginal.GetActive() {
		gui.SetZoomMode("Original")
	}
}

func (gui *GUI) ImageZoomFitWidth() {
	if gui.MenuItemFitToWidth.GetActive() {
		gui.SetZoomMode("FitToWidth")
	}
}

func (gui *GUI) ImageZoomFitHeight() {
	if gui.MenuItemFitToHeight.GetActive() {
		gui.SetZoomMode("FitToHeight")
	}
}

func (gui *GUI) ImageFlipVertically() {
	gui.SetVFlip(gui.Config.VFlip)
}

func (gui *GUI) ImageFlipHorizontally() {
	gui.SetHFlip(gui.Config.HFlip)
}

func (gui *GUI) ModeFullscreen() {
	gui.SetFullscreen(!gui.Config.Fullscreen)
}

func (gui *GUI) ModeSeamless() {
	gui.SetSeamless(!gui.Config.Seamless)
}
func (gui *GUI) ModeManga() {
	gui.SetMangaMode(!gui.Config.MangaMode)
}

func (gui *GUI) ModeDouble() {
	gui.SetDoublePage(!gui.Config.DoublePage)
}

func (gui *GUI) OrderRandom() {
	gui.SetRandom(!gui.Config.Random)
}

func (gui *GUI) PageNext() {
	if !gui.Loaded() {
		if gui.Config.Seamless {
			gui.ArchiveNext()
		}
		return
	}

	if gui.Config.Random {
		gui.RandomPage()
		return
	}

	n := 1
	if gui.Config.DoublePage && gui.forceSinglePage() == false && gui.State.Archive.Len() > gui.State.ArchivePos+2 {
		n = 2
	}

	if gui.Config.Seamless && gui.State.Archive.Len()-gui.State.ArchivePos <= n {
		gui.ArchiveNext()
		return
	}

	gui.SetPage(gui.State.ArchivePos + n)
}

func (gui *GUI) PagePrevious() {
	if !gui.Loaded() {
		if gui.Config.Seamless {
			gui.ArchivePrevious()
		}
		return
	}

	if gui.Config.Random {
		gui.RandomPage()
		return
	}

	n := 1
	if gui.Config.DoublePage && gui.State.ArchivePos > 1 {
		n = 2
	}

	if gui.Config.Seamless && gui.State.ArchivePos+1 <= n {
		gui.ArchivePrevious()
		return
	}

	gui.SetPage(gui.State.ArchivePos - n)

	if (gui.Config.DoublePage && gui.forceSinglePage()) && gui.State.Archive.Len()-gui.State.ArchivePos > 1 {
		gui.PageNext()
	}
}

func (gui *GUI) PageSkipForward() {
	gui.SetPage(gui.State.ArchivePos + gui.Config.NSkip)
}

func (gui *GUI) PageSkipBackward() {
	gui.SetPage(gui.State.ArchivePos - gui.Config.NSkip)
}

func (gui *GUI) PageFirst() {
	if !gui.Loaded() {
		return
	}

	gui.SetPage(0)
}

func (gui *GUI) PageLast() {
	if !gui.Loaded() {
		return
	}

	if gui.Config.DoublePage && gui.State.Archive.Len() >= 2 {
		gui.SetPage(gui.State.Archive.Len() - 2)
	}
	gui.SetPage(gui.State.Archive.Len() - 1)
}

func (gui *GUI) PageGoto() {
	gui.RunGoToDialog()
}

func (gui *GUI) ArchiveNext() {
	newname, _ := gui.archiveNameRel(1)

	gui.LoadArchive(newname)
}

func (gui *GUI) ArchivePrevious() {
	newname, _ := gui.archiveNameRel(-1)

	gui.LoadArchive(newname)
	gui.PageLast()
}

func (gui *GUI) BookmarkAdd() {

}
