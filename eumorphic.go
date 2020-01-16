package main

import (
    "eumorphic/diff"
    "eumorphic/history"
    "github.com/mattn/go-gtk/gtk"
    "gopkg.in/src-d/go-git.v4"
    "os"
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
    history := history.New(repo)
    history.Refresh()
    vbox.Add(history)
    diff := diff.New()
    vbox.Add(diff)
    window.Add(vbox)
    window.ShowAll()
    history.SelectionChanged(diff.Update)
    gtk.Main()
}
