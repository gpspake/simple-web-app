package main

import (
	"html/template"
	"io"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/labstack/echo/v4"
	_ "github.com/mattn/go-sqlite3"
)

// Template struct implements the Echo Renderer interface
type Template struct{}

// Render implements the Echo Renderer interface and renders an HTML template.
// It combines a base layout template ("base.html") with a specific content template
// (e.g., "index.html" or "about.html"), allowing dynamic data to be injected.
//
// Parameters:
//   - w: An `io.Writer` where the rendered template output is written (e.g., HTTP response writer).
//   - name: A string specifying the name of the content template to render (without the ".html" extension).
//   - data: An interface{} containing the dynamic data to pass to the templates.
//   - c: The Echo context, which allows access to request/response details and logging.
//
// Behavior:
//   - Loads the base layout template ("templates/base.html").
//   - Dynamically loads the content template specified by the `name` parameter ("templates/{name}.html").
//   - Renders the templates using the provided `data`, embedding the content into the base layout.
//
// Returns:
//   - An error if template parsing or execution fails.
//   - Otherwise, nil to indicate successful rendering.
func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	tmpl, err := template.ParseFiles(
		"templates/base.html",
		"templates/"+name+".html",
	)
	if err != nil {
		return err
	}
	return tmpl.ExecuteTemplate(w, "base.html", data)
}
