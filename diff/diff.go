package diff

import (
    "eumorphic/diff/lineview"
    "eumorphic/diff/richtext"
    "fmt"
    "github.com/mattn/go-gtk/gtk"
    "gopkg.in/src-d/go-git.v4"
    "gopkg.in/src-d/go-git.v4/plumbing"
    "gopkg.in/src-d/go-git.v4/plumbing/format/diff"
    "gopkg.in/src-d/go-git.v4/plumbing/object"
    "strings"
)

type Diff struct {
    *gtk.HBox
    text *richtext.RichText
    line *lineview.LineView
}

func (d *Diff) Update(repo *git.Repository, hash string) {
    var parent *object.Commit
    var ctree, ptree *object.Tree
    var patch *object.Patch
    commit, err := repo.CommitObject(plumbing.NewHash(hash))
    if err == nil { parent, err = commit.Parent(0) }
    if err == nil { ctree,  err = commit.Tree() }
    if err == nil { ptree,  err = parent.Tree() }
    if err == nil { patch,  err = ptree.Patch(ctree) }
    if err != nil {
        return
    }

    d.text.Clear()
    for _, fp := range patch.FilePatches() {
        from, to := fp.Files()
        if from != nil && to != nil  {
            if from.Path() != to.Path() {
                d.text.Append("file", fmt.Sprintf("~ %s => %s", from.Path(), to.Path()))
            } else {
                d.text.Append("file", "~ " + from.Path())
            }
        } else if from != nil {
            d.text.Append("file", "- " + from.Path())
        } else if to != nil {
            d.text.Append("file", "+ " + to.Path())
        } else {
            d.text.Append("file", "???");
        }
        d.line.Add(0, 0)
        var ( oldline, newline = 0, 0 )
        for _, chunk := range fp.Chunks() {
            style := "normal"
            for _, l := range strings.Split(chunk.Content(), "\n") {
                switch chunk.Type() {
                    case diff.Add:
                        style = "insert"
                        newline++
                        d.line.Add(0, newline)
                    case diff.Delete:
                        style = "delete"
                        oldline++
                        d.line.Add(oldline, 0)
                    default:
                        oldline++
                        newline++
                        d.line.Add(oldline, newline)
                }
                d.text.Append(style, l)
            }
        }
    }
    d.line.Display()
}

func New() *Diff {
    hbox := gtk.NewHBox(false, 0)
    line := lineview.New()
    text := richtext.New()
    text.AddStyle("normal")
    text.AddStyle("file",   "background", "#d9d9d9")
    text.AddStyle("insert", "background", "#aaffaa")
    text.AddStyle("delete", "background", "#ffaaaa")
    hbox.PackStart(line, false, false, 0)
    hbox.PackEnd(text, true, true, 0)
    diff := &Diff {
        hbox,
        text,
        line,
    }
    return diff
}
