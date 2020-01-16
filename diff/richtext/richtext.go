package richtext

import (
    "github.com/mattn/go-gtk/gtk"
)

type RichText struct {
    *gtk.TextView
    styles map[string] *gtk.TextTag
}

func (r *RichText) AddStyle(label string, props ...string) {
    properties := map[string]interface{} { "family": "Monospace" }
    for i := 0; i < len(props) / 2; i++ {
        properties[props[i*2]] = props[i*2+1]
    }
    r.styles[label] = r.GetBuffer().CreateTag(label, properties)
}

func (r *RichText) Clear() {
    var start, end gtk.TextIter
    buffer := r.GetBuffer()
    buffer.GetStartIter(&start)
    buffer.GetEndIter(&end)
    buffer.Delete(&start, &end)
}

func (r *RichText) Append(style string,text string) {
    var end gtk.TextIter
    buffer := r.GetBuffer()
    buffer.GetEndIter(&end)
    r.GetBuffer().InsertWithTag(&end, text + "\n", r.styles[style])
}

func New() *RichText {
    text := gtk.NewTextView()
    text.SetEditable(false)
    return &RichText {
        text,
        make(map[string]*gtk.TextTag),
    }
}
