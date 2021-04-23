package render

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"server/everydaymuslimappserver/internal/config"
	"server/everydaymuslimappserver/internal/models"
	"time"

	"github.com/justinas/nosurf"
)

var app *config.AppConfig
var functions = template.FuncMap{
	"humanDate":    HumanDate,
	"dateWithTime": DateWithTime,
}

var pathToTemplates = "./templates"

//NewRenderer sets the config of the template pkg
func NewRenderer(a *config.AppConfig) {
	app = a
}

//HumanDate returns time in YYYY-MM-DD
func HumanDate(t time.Time) string {
	return t.Format("2006-01-02")
}

//DateWithTime returns the date and hours and minutes
func DateWithTime(t time.Time) string {
	return t.Format("2006-01-02 15:04")
}

func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {

	td.Flash = app.Session.PopString(r.Context(), "flash")
	td.Warning = app.Session.PopString(r.Context(), "warning")
	td.Error = app.Session.PopString(r.Context(), "error")
	td.CSRFToken = nosurf.Token(r)

	if app.Session.Exists(r.Context(), "userId") {
		td.IsAuthenticated = 1
	}
	return td
}

//Template function
func Templates(w http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) error {

	var tc map[string]*template.Template

	if app.UseCache {
		//get the template cache from AppConfig
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()
	}

	t, ok := tc[tmpl]
	if !ok {
		return errors.New("Could not get template from the cache")

	}

	buf := new(bytes.Buffer)

	td = AddDefaultData(td, r)
	_ = t.Execute(buf, td)

	_, err := buf.WriteTo(w)
	if err != nil {
		fmt.Println("error writing template to browser", err)
		return err
	}
	return nil
}

//CreateTemplateCache creates a map of templateCache
func CreateTemplateCache() (map[string]*template.Template, error) {

	myCache := map[string]*template.Template{}

	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.html", pathToTemplates))

	if err != nil {
		return myCache, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.html", pathToTemplates))
		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.html", pathToTemplates))
			if err != nil {
				return myCache, err
			}
		}
		myCache[name] = ts
	}
	return myCache, nil
}
