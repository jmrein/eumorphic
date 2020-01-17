package history

import (
    "eumorphic/listview"
    "fmt"
    "gopkg.in/src-d/go-git.v4"
    "gopkg.in/src-d/go-git.v4/plumbing/object"
    "strings"
)

const (
    Hash = iota
    Commit
    Subject
    Message
    Author
    Date
    NumCols
)

type History struct {
    *listview.ListView
}

func (h *History) add(c *object.Commit) error {
    message := strings.ReplaceAll(c.Message, "&", "&amp;")
    h.AddRow(map[int]string {
            Hash:    c.Hash.String(),
            Commit:  c.Hash.String()[0:7],
            Subject: strings.Split(c.Message, "\n")[0],
            Message: message,
            Author:  fmt.Sprintf("%s <%s>", c.Author.Name, c.Author.Email),
            Date:    c.Author.When.Format("2006-01-02 15:04:05"),
        })
    return nil
}

func (h *History) Refresh(repo *git.Repository) error {
    head, err := repo.Head()
    if err != nil {
        return err
    }
    iter, err := repo.Log(&git.LogOptions {
        From:  head.Hash(),
        Order: git.LogOrderCommitterTime,
    })
    h.Clear()
    h.AddRow(map[int]string { Hash: ":working:", Subject: "(Working directory)" })
    h.AddRow(map[int]string { Hash: ":staged:", Subject: "(Staged)" })
    iter.ForEach(h.add)
    return nil
}

func (h *History) SelectionChanged(onchanged func(hash string)) {
    h.Connect("cursor_changed", func() { onchanged(h.GetSelected(Hash)) })
}

func New() *History {
    tree := listview.New(NumCols)
    tree.AddCol(Commit, "Commit", 0)
    tree.AddCol(Subject, "Subject", 500).SetExpand(true)
    tree.AddCol(Author, "Author", 300)
    tree.AddCol(Date, "Date", 0)
    return &History { tree }
}

