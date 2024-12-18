package internal

import (
	"html/template"
	"io"
	"path/filepath"

	"github.com/labstack/echo/v4"
)

type Template struct {
	TemplateDir string
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	// Check if the request is an HTMX request
	isPartial := c.Request().Header.Get("HX-Request") == "true"

	var tmpl *template.Template
	var err error

	if isPartial {
		// Load only the partial template
		tmpl, err = template.ParseFiles(filepath.Join(t.TemplateDir, name+".html"))
	} else {
		// Load base template and content template
		tmpl, err = template.ParseFiles(
			filepath.Join(t.TemplateDir, "base.html"),             // Base layout
			filepath.Join(t.TemplateDir, "nav.html"),              // Nav template
			filepath.Join(t.TemplateDir, "releases_partial.html"), // Include releases_partial
			filepath.Join(t.TemplateDir, name+".html"),            // Content file
		)
	}

	if err != nil {
		return err
	}

	// Render the template
	if isPartial {
		return tmpl.Execute(w, data)
	}
	return tmpl.ExecuteTemplate(w, "base.html", data)
}
