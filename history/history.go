package history

import (
    "fmt"
    "github.com/mattn/go-gtk/glib"
    "github.com/mattn/go-gtk/gtk"
    "gopkg.in/src-d/go-git.v4"
    "gopkg.in/src-d/go-git.v4/plumbing/object"
    "strings"
)

const (
    Hash = iota
    Commit
    Subject
    Author
    Date
)

func add_col(tree *gtk.TreeView, column int, title string, width int) *gtk.TreeViewColumn {
    col := gtk.NewTreeViewColumnWithAttributes(title, gtk.NewCellRendererText(), "text", column)
    if width > 0 {
        col.SetSizing(gtk.TREE_VIEW_COLUMN_FIXED)
        col.SetFixedWidth(width)
    }
    tree.AppendColumn(col)
    return col
}

type History struct {
    *gtk.TreeView
    repo  *git.Repository
    store *gtk.ListStore
}

func (h *History) add_special(key string, label string) {
    var iter gtk.TreeIter
    h.store.Append(&iter)
    h.store.Set(&iter,
        Hash, key,
        Subject, label,
    )
}

func (h *History) add(c *object.Commit) error {
    lines := strings.Split(c.Message, "\n")
    var iter gtk.TreeIter
    h.store.Append(&iter)
    h.store.Set(&iter,
        Hash, c.Hash.String(),
        Commit, c.Hash.String()[0:7],
        Subject, lines[0],
        Author, fmt.Sprintf("%s <%s>", c.Author.Name, c.Author.Email),
        Date, c.Author.When.Format("2006-01-02 15:04:05"),
    )
    return nil
}

func (h *History) Refresh() error {
    head, err := h.repo.Head()
    if err != nil {
        return err
    }
    iter, err := h.repo.Log(&git.LogOptions {
        From:  head.Hash(),
        Order: git.LogOrderCommitterTime,
    })
    h.store.Clear()
    h.add_special(":working:", "(Working directory)")
    h.add_special(":staged:", "(Staged)")
    iter.ForEach(h.add)
    return nil
}

func (h *History) SelectionChanged(onchanged func(r *git.Repository, hash string)) {
    h.Connect("cursor_changed", func() {
        var iter gtk.TreeIter
        if h.GetSelection().GetSelected(&iter) {
            var value glib.GValue
            h.store.GetValue(&iter, Hash, &value)
            onchanged(h.repo, value.GetString())
        }
    })
}

func New(repo *git.Repository) *History {
    store  := gtk.NewListStore(glib.G_TYPE_STRING, glib.G_TYPE_STRING, glib.G_TYPE_STRING,
        glib.G_TYPE_STRING, glib.G_TYPE_STRING)
    tree := gtk.NewTreeView()
    tree.SetModel(store)
    add_col(tree, Commit, "Commit", 0)
    add_col(tree, Subject, "Subject", 500).SetExpand(true)
    add_col(tree, Author, "Author", 300)
    add_col(tree, Date, "Date", 0)
    return &History { tree, repo, store }
}

