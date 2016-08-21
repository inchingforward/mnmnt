package main

import (
	"flag"
	"html/template"
	"io"
	"log"

	"github.com/inchingforward/mnmnt/handlers"
	"github.com/inchingforward/mnmnt/models"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	_ "github.com/lib/pq"
)

var (
	t     *Template
	debug = false
)

// Template represents the parsed templates from the "templates" directory.
type Template struct {
	templates *template.Template
}

func init() {
	var err error

	models.DB, err = sqlx.Connect("postgres", "user=monument dbname=monument sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
}

// Render renders the template referenced by name and passes the data value
// into the template.  If main was run with the debug argument, the templates
// are re-parsed on each call to Render.
func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	if debug {
		funcMap := template.FuncMap{
			"mdb": handlers.MarkDownBasic,
		}

		t.templates = template.Must(template.New("main").Funcs(funcMap).ParseGlob("templates/*.html"))
	}

	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	e := echo.New()

	funcMap := template.FuncMap{
		"mdb": handlers.MarkDownBasic,
	}

	t = &Template{
		templates: template.Must(template.New("main").Funcs(funcMap).ParseGlob("templates/*.html")),
	}

	flag.BoolVar(&debug, "debug", false, "true to enable debug")
	flag.Parse()

	log.Printf("debug: %v\n", debug)

	e.SetRenderer(t)

	handlers.SetHandlers(e)

	log.Println("Listening...")
	e.Run(standard.New(":4000"))
}
