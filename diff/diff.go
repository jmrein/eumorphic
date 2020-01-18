package diff

import (
    "eumorphic/diff/lineview"
    "eumorphic/diff/richtext"
    "fmt"
    "github.com/mattn/go-gtk/gtk"
    "gopkg.in/libgit2/git2go.v24"
    _ "strings"
)

type Diff struct {
    *gtk.HBox
    text *richtext.RichText
    line *lineview.LineView
}

func (d *Diff) Update(repo *git.Repository, hash string, file string, file_encountered func(key string)) {
    var commit, parent *git.Commit
    var ctree, ptree *git.Tree
    var options git.DiffOptions
    var diff *git.Diff
    oid, err := git.NewOid(hash)
    if err == nil { commit, err = repo.LookupCommit(oid) }
    if err == nil {
        parent = commit.Parent(0)
        ctree, err = commit.Tree()
    }
    if err == nil { ptree, err = parent.Tree() }
    if err == nil { options, err = git.DefaultDiffOptions() }
    if err == nil {
        if file != "" {
            options.Pathspec = []string{file}
        }
        diff,  err = repo.DiffTreeToTree(ptree, ctree, &options)
    }
    if err != nil {
        fmt.Println(err)
        return
    }
    styles := map[git.DiffLineType]string {
        git.DiffLineAddition: "insert",
        git.DiffLineDeletion: "delete",
        git.DiffLineAddEOFNL: "insert",
        git.DiffLineDelEOFNL: "delete",
    }

    d.text.Clear()
    diff.ForEach(func(file git.DiffDelta, progress float64) (git.DiffForEachHunkCallback, error) {
        file_encountered(file.NewFile.Path)
        d.text.Append("file", file.NewFile.Path + "\n")
        d.line.Add(0, 0)
        return func(hunk git.DiffHunk) (git.DiffForEachLineCallback, error) {
            d.text.Append("file", hunk.Header)
            d.line.Add(0, 0)
            return func(line git.DiffLine) error {
                if style, ok := styles[line.Origin]; ok {
                    d.text.Append(style, line.Content)
                } else {
                    d.text.Append("normal", line.Content)
                }
                d.line.Add(line.OldLineno, line.NewLineno)
                return nil
            }, nil
        }, nil
    }, git.DiffDetailLines)
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
