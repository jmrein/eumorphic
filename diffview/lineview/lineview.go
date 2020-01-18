package lineview

import (
	"eumorphic/diffview/richtext"
	"fmt"
	"strconv"
)

//LineView displays line numbers
type LineView struct {
	*richtext.RichText
	oldlines []string
	newlines []string
	max      int
}

//Add to the displayed lines
func (lv *LineView) Add(oldline int, newline int) {
	var (
		o, n = "", ""
	)
	if oldline > 0 {
		o = strconv.Itoa(oldline)
	}
	if newline > 0 {
		n = strconv.Itoa(newline)
	}
	lv.oldlines = append(lv.oldlines, o)
	lv.newlines = append(lv.newlines, n)
	if len(o) > lv.max {
		lv.max = len(o)
	}
	if len(n) > lv.max {
		lv.max = len(n)
	}
}

//Display refreshes the line display
func (lv *LineView) Display() {
	lv.Clear()
	padding := strconv.Itoa(lv.max)
	format := " %" + padding + "s %" + padding + "s \n"
	for i, o := range lv.oldlines {
		n := lv.newlines[i]
		lv.Append("normal", fmt.Sprintf(format, o, n))
	}
	lv.max = 0
	lv.oldlines = make([]string, 0, 1000)
	lv.newlines = make([]string, 0, 1000)
}

//New returns a new line display
func New() *LineView {
	text := richtext.New()
	text.SetCanDefault(false)
	text.SetCanFocus(false)
	text.SetSensitive(false)
	text.AddStyle("normal", "background", "#d9d9d9", "foreground", "blue")
	return &LineView{
		text,
		make([]string, 0, 1000),
		make([]string, 0, 1000),
		0,
	}
}
