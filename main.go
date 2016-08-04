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
	"github.com/satori/go.uuid"
)

var (
	t     *Template
	debug = false
)

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

func renderMessage(c echo.Context, status int, message string) error {
	return c.Render(status, "message.html", message)
}

func render(c echo.Context, templ string, data interface{}, err error) error {
	if err == nil {
		return c.Render(http.StatusOK, templ, data)
	} else if err == sql.ErrNoRows {
		return c.Render(http.StatusNotFound, "404.html", nil)
	} else {
		log.Println(err)
		return c.Render(http.StatusInternalServerError, "500.html", nil)
	}
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	if debug {
		funcMap := template.FuncMap{
			"mdb": markDownBasic,
		}

		t.templates = template.Must(template.New("main").Funcs(funcMap).ParseGlob("templates/*.html"))
	}

	return t.templates.ExecuteTemplate(w, name, data)
}

func index(c echo.Context) error {
	memories, err := models.GetRecentMemories()

	return render(c, "index.html", memories, err)
}

func getMemories(c echo.Context) error {
	memories, err := models.GetAllMemories()

	return render(c, "memories.html", memories, err)
}

func getMemory(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return renderMessage(c, http.StatusBadRequest, fmt.Sprintf("Invalid id: '%v'", c.Param("id")))
	}

	memory, err := models.GetMemory(id)

	return render(c, "memory.html", memory, err)
}

func createMemory(c echo.Context) error {
	m := models.Memory{}
	if err := c.Bind(m); err != nil {
		return err
	}

	if m.Author == "" {
		m.Author = "Anonymous"
	}

	m.ApprovalUuid = uuid.NewV4().String()

	id, err := models.AddMemory(m)
	if err != nil {
		return render(c, "memory.html", m, err)
	}

	log.Printf("New memory \"%v\" (id: %v) created.\n", m.Title, id)

	m.Id = id

	utils.SendEmail(m)
	utils.Tweet(m)

	return render(c, "memory_submitted.html", m, err)
}

func getMemorySubmitted(c echo.Context) error {
	return render(c, "memory_submitted.html", nil, nil)
}

func approveMemory(c echo.Context) error {
	uuid := c.Param("uuid")

	if uuid == "" {
		return c.Render(http.StatusBadRequest, "message.html", "Missing UUID")
	}

	memory, err := models.GetMemoryByUuid(uuid)
	if err != nil {
		return render(c, "", nil, err)
	}

	models.ApproveMemory(memory)
	if err != nil {
		return render(c, "", memory, err)
	}

	return render(c, "memory_approved.html", memory, nil)
}

func getAddMemory(c echo.Context) error {
	return render(c, "memory_form.html", nil, nil)
}

func getAbout(c echo.Context) error {
	return c.Render(http.StatusOK, "about.html", nil)
}

func markDownBasic(args ...interface{}) template.HTML {
	s := blackfriday.MarkdownBasic([]byte(fmt.Sprintf("%s", args...)))
	return template.HTML(s)
}

func main() {
	e := echo.New()

	funcMap := template.FuncMap{
		"mdb": markDownBasic,
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
	e.GET("/", index)
	e.GET("/memories", getMemories)
	e.GET("/memories/:id", getMemory)
	e.POST("/memories", createMemory)
	e.GET("/memories/submitted", getMemorySubmitted)
	e.GET("/memories/approve/:uuid", approveMemory)
	e.GET("/memories/add", getAddMemory)
	e.GET("/about", getAbout)

	log.Println("Listening...")
	e.Run(standard.New(":4000"))
}
