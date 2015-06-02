package views

import (
	"github.com/joncalhoun/viewcon"
	"html/template"
	"log"
	"path/filepath"
)

type MetricsView struct {
	viewcon.View
	Feed   viewcon.Page // A custom page for the cats view
	Browse viewcon.Page
}

var Metrics MetricsView

func MetricsFiles() []string {
	files, err := filepath.Glob("templates/metrics/includes/*.tmpl")
	if err != nil {
		log.Panic(err)
	}
	files = append(files, viewcon.LayoutFiles()...)
	return files
}

func init() {
	indexFiles := append(MetricsFiles(), "templates/metrics/index.tmpl")
	Metrics.Index = viewcon.Page{
		Template: template.Must(template.New("index").ParseFiles(indexFiles...)),
		Layout:   "my_layout",
	}

	browseFiles := append(MetricsFiles(), "templates/metrics/browse.tmpl")
	Metrics.Browse = viewcon.Page{
		Template: template.Must(template.New("browse").ParseFiles(browseFiles...)),
		Layout:   "my_layout",
	}

	feedFiles := append(MetricsFiles(), "templates/metrics/feed.tmpl")
	Metrics.Feed = viewcon.Page{
		Template: template.Must(template.New("feed").ParseFiles(feedFiles...)),
		Layout:   "other_layout",
	}
}
