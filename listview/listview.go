package listview

import (
	"github.com/mattn/go-gtk/glib"
	"github.com/mattn/go-gtk/gtk"
)

type ListView struct {
	*gtk.TreeView
	store *gtk.ListStore
}

func (lv *ListView) AddCol(column int, title string, width int) *gtk.TreeViewColumn {
	col := gtk.NewTreeViewColumnWithAttributes(title, gtk.NewCellRendererText(), "text", column)
	if width > 0 {
		col.SetSizing(gtk.TREE_VIEW_COLUMN_FIXED)
		col.SetFixedWidth(width)
	}
	lv.AppendColumn(col)
	return col
}

func (lv *ListView) AddRow(data map[int]string) {
	var iter gtk.TreeIter
	lv.store.Append(&iter)
	for k, v := range data {
		lv.store.Set(&iter, k, v)
	}
}

func (lv *ListView) GetSelected(column int) string {
	var iter gtk.TreeIter
	if lv.GetSelection().GetSelected(&iter) {
		var value glib.GValue
		lv.store.GetValue(&iter, column, &value)
		return value.GetString()
	}
	return ""
}

func (lv *ListView) Clear() {
	lv.store.Clear()
}

func New(columns int) *ListView {
	types := make([]interface{}, columns)
	for i := 0; i < columns; i++ {
		types[i] = glib.G_TYPE_STRING
	}
	store := gtk.NewListStore(types...)
	tree := gtk.NewTreeView()
	tree.SetModel(store)
	return &ListView{tree, store}
}
