package diffview

import (
	"eumorphic/diffview/lineview"
	"eumorphic/diffview/richtext"

	"github.com/mattn/go-gtk/gtk"
	git "gopkg.in/libgit2/git2go.v24"
)

//DiffView displays the diff for a commit to its parent
type DiffView struct {
	*gtk.HBox
	text *richtext.RichText
	line *lineview.LineView
}

func getTree(repo *git.Repository, hash string) (*git.Commit, *git.Tree, error) {
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

func getHead(repo *git.Repository) (*git.Tree, error) {
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

func getDiff(repo *git.Repository, hash string, files []string) (*git.Diff, error) {
	options, err := git.DefaultDiffOptions()
	if err != nil {
		return nil, err
	}
	if files != nil {
		options.Pathspec = files
	}
	switch hash {
	case ":working:":
		return repo.DiffIndexToWorkdir(nil, &options)
	case ":staged:":
		tree, err := getHead(repo)
		if err != nil {
			return nil, err
		}
		return repo.DiffTreeToIndex(tree, nil, &options)
	}
	commit, tree, err := getTree(repo, hash)
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

/*Update shows the diff for a commit to its parent
hash - the focus commit hash, or a special string
    :working: will compare the working directory to the staged directory
    :staged: will compare the staged version to head
file - if not blank, show only this file*/
func (d *DiffView) Update(repo *git.Repository, hash string, files []string) ([]git.DiffDelta, error) {
	diff, err := getDiff(repo, hash, files)
	if err != nil {
		return nil, err
	}
	options, err := git.DefaultDiffFindOptions()
	if err != nil {
		return nil, err
	}
	err = diff.FindSimilar(&options)
	if err != nil {
		return nil, err
	}
	styles := map[git.DiffLineType]string{
		git.DiffLineAddition: "insert",
		git.DiffLineDeletion: "delete",
		git.DiffLineAddEOFNL: "insert",
		git.DiffLineDelEOFNL: "delete",
	}

	d.text.Clear()
	deltas := make([]git.DiffDelta, 0)
	diff.ForEach(func(file git.DiffDelta, progress float64) (git.DiffForEachHunkCallback, error) {
		deltas = append(deltas, file)
		//fmt.Printf("%s => %s (%d)\n", file.OldFile.Path, file.NewFile.Path, file.Status)
		if file.OldFile.Path != file.NewFile.Path {
			d.text.Append("file", file.OldFile.Path+" => ")
		}
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
	return deltas, nil
}

//New returns a DiffView
func New() *DiffView {
	hbox := gtk.NewHBox(false, 0)
	line := lineview.New()
	text := richtext.New()
	text.AddStyle("normal")
	text.AddStyle("file", "background", "#d9d9d9")
	text.AddStyle("insert", "background", "#aaffaa")
	text.AddStyle("delete", "background", "#ffaaaa")
	hbox.PackStart(line, false, false, 0)
	hbox.PackEnd(text, true, true, 0)
	diff := &DiffView{
		hbox,
		text,
		line,
	}
	return diff
}
