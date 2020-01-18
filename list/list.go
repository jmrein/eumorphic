package list

import (
	"eumorphic/listview"
)

//List displays a list with a single column with no header
type List struct {
	*listview.ListView
}

//Add a value to the list
func (l *List) Add(value string) {
	l.AddRow(map[int]string{0: value})
}

//SelectionChanged adds a listener for when the selected item changes
func (l *List) SelectionChanged(onchanged func(value string)) {
	l.Connect("cursor_changed", func() { onchanged(l.GetSelected(0)) })
}

//New returns a new single-list displayer
func New() *List {
	tree := listview.New(1)
	tree.AddCol(0, "Value", 0).SetExpand(true)
	tree.SetHeadersVisible(false)
	return &List{tree}
}
