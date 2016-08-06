package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/inchingforward/mnmnt/models"
	"github.com/inchingforward/mnmnt/utils"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	_ "github.com/lib/pq"
	"github.com/russross/blackfriday"
)

var (
	t     *Template
	debug = false
)

// Template represents the parsed templates from the "templates" directory.
type Template struct {
	templates *template.Template
}

// A TemplateContext holds data and an error that are passed in to a
// template for rendering.  Either Data or Err can be nil.
type TemplateContext struct {
	Data interface{}
	Err  error
}

func init() {
	var err error

	models.DB, err = sqlx.Connect("postgres", "user=monument dbname=monument sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
}

// RenderMessage displays message using the message.html template.
func RenderMessage(c echo.Context, status int, message string) error {
	return c.Render(status, "message.html", message)
}

// Render will render the given templ template passing in data if err
// is nil, or will render a 404 or 500 error page depending on the err.
func Render(c echo.Context, templ string, data interface{}, err error) error {
	if err == nil {
		return c.Render(http.StatusOK, templ, data)
	} else if err == sql.ErrNoRows {
		return c.Render(http.StatusNotFound, "404.html", nil)
	} else {
		log.Println(err)
		return c.Render(http.StatusInternalServerError, "500.html", err)
	}
}

// RenderContext will render the given templ template passing in the ctx
// TemplateContext.
func RenderContext(c echo.Context, templ string, ctx TemplateContext) error {
	return c.Render(http.StatusOK, templ, ctx)
}

// Render renders the template referenced by name and passes the data value
// into the template.  If main was run with the debug argument, the templates
// are re-parsed on each call to Render.
func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	if debug {
		funcMap := template.FuncMap{
			"mdb": MarkDownBasic,
		}

		t.templates = template.Must(template.New("main").Funcs(funcMap).ParseGlob("templates/*.html"))
	}

	return t.templates.ExecuteTemplate(w, name, data)
}

// Index renders the home page.
func Index(c echo.Context) error {
	memories, err := models.GetRecentMemories()

	return Render(c, "index.html", memories, err)
}

// GetMemories renders all approved memories.
func GetMemories(c echo.Context) error {
	memories, err := models.GetAllMemories()

	return Render(c, "memories.html", memories, err)
}

// GetMemory renders the memory details page for the corresponding
// memory id parameter.
func GetMemory(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return RenderMessage(c, http.StatusBadRequest, fmt.Sprintf("Invalid id: '%v'", c.Param("id")))
	}

	memory, err := models.GetMemory(id)

	return Render(c, "memory.html", memory, err)
}

// CreateMemory creates a new Memory using values from a submitted form.
func CreateMemory(c echo.Context) error {
	m := models.Memory{}
	if err := c.Bind(&m); err != nil {
		return Render(c, "memory.html", m, err)
	}

	err := models.AddMemory(&m)
	if err != nil {
		return RenderContext(c, "memory_form.html", TemplateContext{c.FormParams(), err})
	}

	utils.SendEmail(m)
	utils.Tweet(m)

	return Render(c, "memory_submitted.html", m, err)
}

// GetMemorySubmitted renders the memory submitted success page.
func GetMemorySubmitted(c echo.Context) error {
	return Render(c, "memory_submitted.html", nil, nil)
}

// ApproveMemory approves the Memory corresponding to the uuid parameter.
func ApproveMemory(c echo.Context) error {
	uuid := c.Param("uuid")

	if uuid == "" {
		return c.Render(http.StatusBadRequest, "message.html", "Missing UUID")
	}

	memory, err := models.GetMemoryByUuid(uuid)
	if err != nil {
		return Render(c, "", nil, err)
	}

	models.ApproveMemory(memory)
	if err != nil {
		return Render(c, "", memory, err)
	}

	return Render(c, "memory_approved.html", memory, nil)
}

// GetAddMemory renders the create memory form.
func GetAddMemory(c echo.Context) error {
	return Render(c, "memory_form.html", nil, nil)
}

// GetAbout renders the About page.
func GetAbout(c echo.Context) error {
	return c.Render(http.StatusOK, "about.html", nil)
}

// MarkDownBasic passes the given data to the MarkdownBasic formatter.
func MarkDownBasic(args ...interface{}) template.HTML {
	s := blackfriday.MarkdownBasic([]byte(fmt.Sprintf("%s", args...)))
	return template.HTML(s)
}

func main() {
	e := echo.New()

	funcMap := template.FuncMap{
		"mdb": MarkDownBasic,
	}

	t = &Template{
		templates: template.Must(template.New("main").Funcs(funcMap).ParseGlob("templates/*.html")),
	}

	if len(os.Args) > 1 && os.Args[1] == "debug" {
		debug = true
	}

	log.Printf("debug: %v\n", debug)

	e.SetRenderer(t)

	e.Static("/static", "static")
	e.GET("/", Index)
	e.GET("/memories", GetMemories)
	e.GET("/memories/:id", GetMemory)
	e.POST("/memories", CreateMemory)
	e.GET("/memories/submitted", GetMemorySubmitted)
	e.GET("/memories/approve/:uuid", ApproveMemory)
	e.GET("/memories/add", GetAddMemory)
	e.GET("/about", GetAbout)

	log.Println("Listening...")
	e.Run(standard.New(":4000"))
}
