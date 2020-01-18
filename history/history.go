package history

import (
	"eumorphic/listview"
	"fmt"
	"strings"

	git "gopkg.in/libgit2/git2go.v24"
)

//Values stored in the list store
const (
	Hash = iota
	Commit
	Subject
	Message
	Author
	Date
	NumCols
)

//History displays the git history
type History struct {
	*listview.ListView
}

func (h *History) add(c *git.Commit) bool {
	message := strings.TrimSpace(strings.ReplaceAll(c.Message(), "&", "&amp;"))
	hash := c.Id().String()
	h.AddRow(map[int]string{
		Hash:    hash,
		Commit:  hash[0:7],
		Subject: strings.Split(c.Message(), "\n")[0],
		Message: message,
		Author:  fmt.Sprintf("%s <%s>", c.Author().Name, c.Author().Email),
		Date:    c.Author().When.Format("2006-01-02 15:04:05"),
	})
	return true
}

//Refresh refreshes the display
func (h *History) Refresh(repo *git.Repository) {
	h.Clear()
	h.AddRow(map[int]string{Hash: ":working:", Subject: "(Working directory)"})
	h.AddRow(map[int]string{Hash: ":staged:", Subject: "(Staged)"})
	walk, err := repo.Walk()
	if err == nil {
		walk.Sorting(git.SortTime)
		err = walk.PushHead()
	}
	if err == nil {
		err = walk.Iterate(h.add)
	}
	if err != nil {
		fmt.Println(err)
	}
}

//SelectionChanged adds a listener for the selected commit
func (h *History) SelectionChanged(onchanged func(hash string)) {
	h.Connect("cursor_changed", func() { onchanged(h.GetSelected(Hash)) })
}

//New returns a new history viewer
func New() *History {
	tree := listview.New(NumCols)
	tree.AddCol(Commit, "Commit", 0)
	tree.AddCol(Subject, "Subject", 500).SetExpand(true)
	tree.AddCol(Author, "Author", 300)
	tree.AddCol(Date, "Date", 0)
	return &History{tree}
}
