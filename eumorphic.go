package main

import (
	"eumorphic/diffview"
	"eumorphic/filelist"
	"eumorphic/history"
	"fmt"
	"os"
	"github.com/mattn/go-gtk/gtk"
	"gopkg.in/libgit2/git2go.v24"
)

type eumorphic struct {
	window *gtk.Window
	hist   *history.History
	repo   *git.Repository
}

func (e *eumorphic) open(dir string) error {
	repo, err := git.OpenRepository(dir)
	if err != nil {
		return err
	}
	e.window.SetTitle(dir)
	e.hist.Refresh(repo)
	e.repo = repo
	return nil
}

func (e *eumorphic) showError(text string) {
	dlg := gtk.NewMessageDialog(
		e.window,
		gtk.DIALOG_MODAL,
		gtk.MESSAGE_ERROR,
		gtk.BUTTONS_OK,
		text)
	dlg.Run()
	dlg.Response(dlg.Destroy)
}

func (e *eumorphic) openClicked() {
	dlg := gtk.NewFileChooserDialog(
		"Open repository...",
		e.window,
		gtk.FILE_CHOOSER_ACTION_SELECT_FOLDER,
		gtk.STOCK_CANCEL,
		gtk.RESPONSE_CANCEL,
		gtk.STOCK_OK,
		gtk.RESPONSE_ACCEPT)
	if response := dlg.Run(); response == gtk.RESPONSE_ACCEPT {
		if err := e.open(dlg.GetFilename()); err != nil {
			e.showError(err.Error())
		}
	}
	dlg.Destroy()
}

func main() {
	gtk.Init(&os.Args)
	window := gtk.NewWindow(gtk.WINDOW_TOPLEVEL)
	window.SetPosition(gtk.WIN_POS_CENTER)
	window.Connect("destroy", gtk.MainQuit)
	window.SetSizeRequest(1100, 600)

	vbox := gtk.NewVBox(false, 0)
	toolbar := gtk.NewToolbar()
	toolbar.SetStyle(gtk.TOOLBAR_ICONS)
	btnOpen := gtk.NewToolButtonFromStock(gtk.STOCK_OPEN)
	btnRefresh := gtk.NewToolButtonFromStock(gtk.STOCK_REFRESH)
	toolbar.Insert(btnOpen, -1)
	toolbar.Insert(btnRefresh, -1)
	vbox.PackStart(toolbar, false, false, 0)

	hist := history.New()
	scroll := gtk.NewScrolledWindow(nil, nil)
	scroll.SetPolicy(gtk.POLICY_NEVER, gtk.POLICY_AUTOMATIC)
	scroll.Add(hist)
	vbox.PackStart(scroll, true, true, 0)

	hbox := gtk.NewHBox(false, 0)
	files := filelist.New()
	files.SetSizeRequest(250, 150)
	scroll = gtk.NewScrolledWindow(nil, nil)
	scroll.SetPolicy(gtk.POLICY_AUTOMATIC, gtk.POLICY_AUTOMATIC)
	scroll.Add(files)
	hbox.PackStart(scroll, false, false, 0)
	diff := diffview.New()
	scroll = gtk.NewScrolledWindow(nil, nil)
	scroll.SetPolicy(gtk.POLICY_AUTOMATIC, gtk.POLICY_AUTOMATIC)
	scroll.AddWithViewPort(diff)
	hbox.PackEnd(scroll, true, true, 0)
	vbox.PackEnd(hbox, true, true, 4)

	window.Add(vbox)
	self := &eumorphic{
		window,
		hist,
		nil,
	}

	hist.SelectionChanged(func(h string) {
		deltas, err := diff.Update(self.repo, h, nil)
		if err == nil {
			files.Clear()
			for _, d := range deltas {
				files.Add(d)
			}
		} else {
			fmt.Println(err)
		}
	})
	files.SelectionChanged(func(o, n string) {
		diff.Update(self.repo, hist.GetSelected(history.Hash), []string{o, n})
	})

	dir, _ := os.Getwd()
	if len(os.Args) >= 2 {
		dir = os.Args[1]
	}
	self.open(dir)
	btnOpen.OnClicked(self.openClicked)
	btnRefresh.OnClicked(func() { hist.Refresh(self.repo) })
	window.ShowAll()
	gtk.Main()
}
