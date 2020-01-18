package list

import (
	"eumorphic/listview"
)

type List struct {
	*listview.ListView
}

func (l *List) Add(value string) {
	l.AddRow(map[int]string{0: value})
}

func (l *List) SelectionChanged(onchanged func(value string)) {
	l.Connect("cursor_changed", func() { onchanged(l.GetSelected(0)) })
}

func New() *List {
	tree := listview.New(1)
	tree.AddCol(0, "Value", 0).SetExpand(true)
	tree.SetHeadersVisible(false)
	return &List{tree}
}
