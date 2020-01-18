package richtext

import (
	"github.com/mattn/go-gtk/gtk"
)

//RichText is a more convenient interface to a rich TextView
type RichText struct {
	*gtk.TextView
	styles map[string]*gtk.TextTag
}

//AddStyle creates a style and adds it to a map
func (r *RichText) AddStyle(label string, props ...string) {
	properties := map[string]interface{}{"family": "Monospace"}
	for i := 0; i < len(props)/2; i++ {
		properties[props[i*2]] = props[i*2+1]
	}
	r.styles[label] = r.GetBuffer().CreateTag(label, properties)
}

//Clear clears all text
func (r *RichText) Clear() {
	var start, end gtk.TextIter
	buffer := r.GetBuffer()
	buffer.GetStartIter(&start)
	buffer.GetEndIter(&end)
	buffer.Delete(&start, &end)
}

//Append appends text to the RichText
func (r *RichText) Append(style string, text string) {
	var end gtk.TextIter
	buffer := r.GetBuffer()
	buffer.GetEndIter(&end)
	r.GetBuffer().InsertWithTag(&end, text, r.styles[style])
}

//New returns a new RichText
func New() *RichText {
	text := gtk.NewTextView()
	text.SetEditable(false)
	return &RichText{
		text,
		make(map[string]*gtk.TextTag),
	}
}
