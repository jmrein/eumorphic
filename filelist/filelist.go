package filelist

import (
	"eumorphic/listview"

	git "gopkg.in/libgit2/git2go.v24"
)

//Values stored in the list store
const (
	OldFile = iota
	NewFile
	NumCols
)

//FileList displays a list of files changed in the selected commit
type FileList struct {
	*listview.ListView
}

//Add a value to the list
func (l *FileList) Add(delta git.DiffDelta) {
	l.AddRow(map[int]string{OldFile: delta.OldFile.Path, NewFile: delta.NewFile.Path})
}

//SelectionChanged adds a listener for when the selected item changes
func (l *FileList) SelectionChanged(onchanged func(oldfile, newfile string)) {
	l.Connect("cursor_changed", func() { onchanged(l.GetSelected(OldFile), l.GetSelected(NewFile)) })
}

//New returns a new single-list displayer
func New() *FileList {
	tree := listview.New(NumCols)
	tree.AddCol(NewFile, "New", 0).SetExpand(true)
	tree.SetHeadersVisible(false)
	return &FileList{tree}
}
