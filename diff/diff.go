package diff

import (
    "github.com/mattn/go-gtk/gtk"
    "gopkg.in/src-d/go-git.v4"
    "gopkg.in/src-d/go-git.v4/plumbing"
    "gopkg.in/src-d/go-git.v4/plumbing/object"
)

type Diff struct {
    *gtk.ScrolledWindow
    text *gtk.TextView
    styles map[string]*gtk.TextTag
}

func (d *Diff) add_style(name string, props map[string]interface{}) {
    props["family"] = "Monospace"
    d.styles[name] = d.text.GetBuffer().CreateTag(name, props)
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

    buffer := d.text.GetBuffer()
    var start, end gtk.TextIter
    buffer.GetStartIter(&start)
    buffer.GetEndIter(&end)
    buffer.Delete(&start, &end)
    //buffer.GetStartIter(&start)
    buffer.InsertWithTag(&start, patch.String(), d.styles["normal"])
}

func New() *Diff {
    scroll := gtk.NewScrolledWindow(nil, nil)
    scroll.SetPolicy(gtk.POLICY_AUTOMATIC, gtk.POLICY_AUTOMATIC)
    text := gtk.NewTextView()
    text.SetEditable(false)
    scroll.Add(text)
    diff := &Diff {
        scroll,
        text,
        make(map[string]*gtk.TextTag),
    }
    diff.add_style("normal", map[string]interface{}{})
    diff.add_style("file",   map[string]interface{}{"background": "#d9d9d9"})
    diff.add_style("insert", map[string]interface{}{"background": "#aaffaa"})
    diff.add_style("delete", map[string]interface{}{"background": "#ffaaaa"})
    return diff
}
