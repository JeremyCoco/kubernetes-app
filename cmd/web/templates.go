package main

import (
	"html/template"
	"path/filepath"
)

func createTemplatesMap(dir string) (map[string]*template.Template, error) {
	views, err := filepath.Glob(filepath.Join(dir, "views/*.go.html"))
	if err != nil {
		return nil, err
	}

	templates := make(map[string]*template.Template)

	for _, view := range views {
		name := filepath.Base(view)

		ts, err := template.New(name).ParseFiles(view)
		if err != nil {
			return nil, err

		}

		ts, err = ts.ParseGlob(filepath.Join(dir, "layout/*.go.html"))
		if err != nil {
			return nil, err

		}

		ts, err = ts.ParseGlob(filepath.Join(dir, "components/*.go.html"))
		if err != nil {
			return nil, err

		}

		templates[name] = ts
	}

	return templates, nil
}
