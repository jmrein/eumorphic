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
    d.styles[label] = r.text.GetBuffer().CreateTag(label, properties)
}


func New() *RichText {
    text := gtk.NewTextView()
    text.SetEditable(false)
    return &RichText {
        text,
        make(map[string]*gtk.TextTag)
    }
}
