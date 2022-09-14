package render

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/ninostephen/booking/pkg/config"
	"github.com/ninostephen/booking/pkg/models"
)

var app *config.AppConfig

// NewTemplate sets the config for template package
func NewTemplates(a *config.AppConfig) {
	app = a
}

func AddDefaultData(td *models.TemplateData) *models.TemplateData {
	return td
}

// RenderTemplate renders html files using html/template
func RenderTemplate(w http.ResponseWriter, tmpl string, td *models.TemplateData) {
	var tc map[string]*template.Template
	var err error
	if app.UseCache {
		// get template cache from app config
		tc = app.TemplateCache
	} else {
		tc, err = CreateTemplateCache()
		if err != nil {
			log.Fatal(err)
		}
	}

	// get requested template from cache
	t, ok := tc[tmpl]
	if !ok {
		log.Fatal("failed to load template from cache")
	}

	buf := new(bytes.Buffer)

	td = AddDefaultData(td)
	err = t.Execute(buf, td)
	if err != nil {
		log.Println(err)
	}

	// render the template
	_, err = buf.WriteTo(w)
	if err != nil {
		log.Println(err)
	}
}

func CreateTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}
	// get all of the files named *.pages.tmpl from templates folder
	pages, err := filepath.Glob("./templates/*.page.tmpl")
	if err != nil {
		return myCache, err
	}
	// range through the pages
	for _, page := range pages {
		// parse the file name of the page
		name := filepath.Base(page)
		// create a new template set with the same name as the file and parse it.
		ts, err := template.New(name).ParseFiles(page)
		if err != nil {
			return myCache, err
		}
		// see whether there are layout files associated with it.
		matches, err := filepath.Glob("./templates/*.layout.tmpl")
		if err != nil {
			return myCache, err
		}
		// are there any layout files?
		if len(matches) > 0 {
			// if yes, then parse them and add it to the template set
			ts, err = ts.ParseGlob("./templates/*.layout.tmpl")
			if err != nil {

				return myCache, err
			}
		}
		// create a new entry to the cache

		myCache[name] = ts
	}
	return myCache, nil
}
