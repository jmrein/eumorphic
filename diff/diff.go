package diff

import (
	"eumorphic/diff/lineview"
	"eumorphic/diff/richtext"
	"fmt"
	"github.com/mattn/go-gtk/gtk"
	"gopkg.in/libgit2/git2go.v24"
)

type Diff struct {
	*gtk.HBox
	text *richtext.RichText
	line *lineview.LineView
}

func get_tree(repo *git.Repository, hash string) (*git.Commit, *git.Tree, error) {
	oid, err := git.NewOid(hash)
	if err != nil {
		return nil, nil, err
	}
	commit, err := repo.LookupCommit(oid)
	if err != nil {
		return nil, nil, err
	}
	tree, err := commit.Tree()
	return commit, tree, err
}

func get_head(repo *git.Repository) (*git.Tree, error) {
	head, err := repo.Head()
	if err != nil {
		return nil, err
	}
	object, err := head.Peel(git.ObjectTree)
	if err != nil {
		return nil, err
	}
	return object.AsTree()
}

func get_diff(repo *git.Repository, hash string, file string) (*git.Diff, error) {
	options, err := git.DefaultDiffOptions()
	if err != nil {
		return nil, err
	}
	if file != "" {
		options.Pathspec = []string{file}
	}
	switch hash {
	case ":working:":
		return repo.DiffIndexToWorkdir(nil, &options)
	case ":staged:":
		tree, err := get_head(repo)
		if err != nil {
			return nil, err
		}
		return repo.DiffTreeToIndex(tree, nil, &options)
	}
	commit, tree, err := get_tree(repo, hash)
	if err != nil {
		return nil, err
	}
	var parent *git.Tree
	if commit.ParentCount() > 0 {
		parent, err = commit.Parent(0).Tree()
	}
	if err != nil {
		return nil, err
	}
	return repo.DiffTreeToTree(parent, tree, &options)
}

func (d *Diff) Update(repo *git.Repository, hash string, file string, file_encountered func(key string)) {
	diff, err := get_diff(repo, hash, file)
	if err != nil {
		fmt.Println(err)
		return
	}
	styles := map[git.DiffLineType]string{
		git.DiffLineAddition: "insert",
		git.DiffLineDeletion: "delete",
		git.DiffLineAddEOFNL: "insert",
		git.DiffLineDelEOFNL: "delete",
	}

	d.text.Clear()
	diff.ForEach(func(file git.DiffDelta, progress float64) (git.DiffForEachHunkCallback, error) {
		file_encountered(file.NewFile.Path)
		d.text.Append("file", file.NewFile.Path+"\n")
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
	text.AddStyle("file", "background", "#d9d9d9")
	text.AddStyle("insert", "background", "#aaffaa")
	text.AddStyle("delete", "background", "#ffaaaa")
	hbox.PackStart(line, false, false, 0)
	hbox.PackEnd(text, true, true, 0)
	diff := &Diff{
		hbox,
		text,
		line,
	}
	return diff
}
