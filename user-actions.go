package main

import (
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
	"log"
	"path/filepath"
	"time"
)

type UserAction struct {
	Keybind KeyChord
	Action  func()
}

type UserActions struct {
	FileOpen      UserAction
	FileClose     UserAction
	FileSaveImage UserAction
	Quit          UserAction

	Preferences UserAction
	About       UserAction

	ImageShrinkLarge      UserAction
	ImageEnlargeSmall     UserAction
	ImageZoomBestFit      UserAction
	ImageZoomOriginal     UserAction
	ImageZoomFitWidth     UserAction
	ImageZoomFitHeight    UserAction
	ImageFlipVertically   UserAction
	ImageFlipHorizontally UserAction

	ModeFullscreen UserAction
	ModeSeamless   UserAction
	ModeManga      UserAction
	ModeDouble     UserAction

	OrderRandom UserAction

	PageNext         UserAction
	PagePrevious     UserAction
	PageSkipForward  UserAction
	PageSkipBackward UserAction
	PageFirst        UserAction
	PageLast         UserAction
	PageGoto         UserAction
	ArchiveNext      UserAction
	ArchivePrevious  UserAction

	BookmarkAdd UserAction
}

func (ga *UserActions) Init(gui *GUI) {
	ga.FileOpen = UserAction{KeyChord{Key: gdk.KEY_o, Ctrl: true}, gui.FileOpen}
	ga.FileClose = UserAction{KeyChord{Key: gdk.KEY_w, Ctrl: true}, gui.FileClose}
	ga.FileSaveImage = UserAction{KeyChord{Key: gdk.KEY_s, Ctrl: true}, gui.FileSaveImage}
	ga.Quit = UserAction{KeyChord{Key: gdk.KEY_q, Ctrl: true}, gui.Quit}

	ga.Preferences = UserAction{KeyChord{Key: gdk.KEY_p, Ctrl: true}, gui.Quit}
	ga.About = UserAction{KeyChord{Key: gdk.KEY_F1}, gui.About}

	ga.PageNext = UserAction{KeyChord{Key: gdk.KEY_Down}, gui.PageNext}
	ga.PagePrevious = UserAction{KeyChord{Key: gdk.KEY_Up}, gui.PagePrevious}
	ga.PageSkipForward = UserAction{KeyChord{Key: gdk.KEY_Right}, gui.PageSkipForward}
	ga.PageSkipBackward = UserAction{KeyChord{Key: gdk.KEY_Left}, gui.PageSkipBackward}
	ga.PageFirst = UserAction{KeyChord{Key: gdk.KEY_Home}, gui.PageFirst}
	ga.PageLast = UserAction{KeyChord{Key: gdk.KEY_End}, gui.PageLast}
	ga.ArchiveNext = UserAction{KeyChord{Key: gdk.KEY_Page_Down, Ctrl: true}, gui.ArchiveNext}
	ga.ArchivePrevious = UserAction{KeyChord{Key: gdk.KEY_Page_Up, Ctrl: true}, gui.ArchivePrevious}

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
	gui.SetShrink(gui.MenuItemShrink.GetActive())
}

func (gui *GUI) ImageEnlargeSmall() {
	gui.SetEnlarge(gui.MenuItemEnlarge.GetActive())
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
	gui.SetVFlip(gui.MenuItemVFlip.GetActive())
}

func (gui *GUI) ImageFlipHorizontally() {
	gui.SetHFlip(gui.MenuItemHFlip.GetActive())
}

func (gui *GUI) ModeFullscreen() {
	gui.SetFullscreen(gui.MenuItemFullscreen.GetActive())
}

func (gui *GUI) ModeSeamless() {
	gui.SetSeamless(gui.MenuItemSeamless.GetActive())
}
func (gui *GUI) ModeManga() {
	gui.SetMangaMode(gui.MenuItemMangaMode.GetActive())
}

func (gui *GUI) ModeDouble() {
	gui.SetDoublePage(gui.MenuItemDoublePage.GetActive())
}

func (gui *GUI) OrderRandom() {
	gui.SetRandom(gui.MenuItemRandom.GetActive())
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
