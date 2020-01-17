package main

import (
	"eumorphic/diff"
	"eumorphic/history"
	"eumorphic/list"
	"os"

	"github.com/mattn/go-gtk/gtk"
	"gopkg.in/src-d/go-git.v4"
)

func main() {
	gtk.Init(&os.Args)
	dir, _ := os.Getwd()
	if len(os.Args) >= 2 {
		dir = os.Args[1]
	}
	window := gtk.NewWindow(gtk.WINDOW_TOPLEVEL)
	window.SetPosition(gtk.WIN_POS_CENTER)
	window.SetTitle(dir)
	window.Connect("destroy", gtk.MainQuit)
	window.SetSizeRequest(1100, 600)

	repo, _ := git.PlainOpen(dir)
	vbox := gtk.NewVBox(true, 8)
	hist := history.New()
	hist.Refresh(repo)
	scroll := gtk.NewScrolledWindow(nil, nil)
	scroll.SetPolicy(gtk.POLICY_NEVER, gtk.POLICY_AUTOMATIC)
	scroll.Add(hist)
	vbox.Add(scroll)

	hbox := gtk.NewHBox(false, 0)
	files := list.New()
	files.SetSizeRequest(250, 150)
	scroll = gtk.NewScrolledWindow(nil, nil)
	scroll.SetPolicy(gtk.POLICY_AUTOMATIC, gtk.POLICY_AUTOMATIC)
	scroll.Add(files)
	hbox.PackStart(scroll, false, false, 0)
	diff := diff.New()
	scroll = gtk.NewScrolledWindow(nil, nil)
	scroll.SetPolicy(gtk.POLICY_AUTOMATIC, gtk.POLICY_AUTOMATIC)
	scroll.AddWithViewPort(diff)
	hbox.PackEnd(scroll, true, true, 0)
	vbox.Add(hbox)

	window.Add(vbox)
	window.ShowAll()
	hist.SelectionChanged(func(h string) {
		files.Clear()
		diff.Update(repo, h, "", func(key string) { files.Add(key) })
	})
	files.SelectionChanged(func(k string) {
		diff.Update(repo, hist.GetSelected(history.Hash), k, func(key string) {})
	})
	gtk.Main()
}
