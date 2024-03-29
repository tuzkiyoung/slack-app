package views

import (
	"bytes"
	"html/template"
	"io/fs"
)

func renderTemplate(f fs.FS, file string, args interface{}) bytes.Buffer {

	var tpl bytes.Buffer

	// read the block-kit definition as a go template
	t, err := template.ParseFS(f, file)
	if err != nil {
		panic(err)
	}

	// we render the view
	err = t.Execute(&tpl, args)
	if err != nil {
		panic(err)
	}

	return tpl
}
