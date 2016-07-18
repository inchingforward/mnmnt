package main

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	_ "github.com/lib/pq"
	"github.com/russross/blackfriday"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

var (
	db    *sqlx.DB
	t     *Template
	debug = false
)

type Template struct {
	templates *template.Template
}

type Memory struct {
	Id           uint64         `db:"id"`
	Title        string         `db:"title"`
	Details      string         `db:"details"`
	Latitude     float64        `db:"latitude"`
	Longitude    float64        `db:"longitude"`
	Author       string         `db:"author"`
	IsApproved   bool           `db:"is_approved"`
	ApprovalUuid sql.NullString `db:"approval_uuid"`
	InsertedAt   time.Time      `db:"inserted_at"`
	UpdatedAt    time.Time      `db:"updated_at"`
}

func init() {
	var err error

	db, err = sqlx.Connect("postgres", "user=monument dbname=monument sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
}

func renderFixMe(c echo.Context, message string) error {
	return c.String(http.StatusOK, message)
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
		funcMap := template.FuncMap {
			"mdb": markDownBasic,
		}

		t.templates = template.Must(template.New("main").Funcs(funcMap).ParseGlob("templates/*.html"))
	}

	return t.templates.ExecuteTemplate(w, name, data)
}

func index(c echo.Context) error {
	var memories []*Memory
	err := db.Select(&memories, "select * from memory order by id desc limit 5")

	return render(c, "index.html", memories, err)
}

func getMemories(c echo.Context) error {
	return renderFixMe(c, "FIXME:  render list of memories")
}

func getMemory(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return renderFixMe(c, err.Error())
	}	

	memory := Memory{}
	err = db.Get(&memory, "select * from memory where id = $1", id)
	
	return render(c, "memory.html", memory, err)
}

func createMemory(c echo.Context) error {
	return renderFixMe(c, "FIXME:  create memory")
}

func updateMemory(c echo.Context) error {
	return renderFixMe(c, "FIXME:  update memory")
}

func getMemorySubmitted(c echo.Context) error {
	return renderFixMe(c, "FIXME:  get memory submitted")
}

func approveMemory(c echo.Context) error {
	return renderFixMe(c, "FIXME:  approve memory")
}

func getAddMemory(c echo.Context) error {
	return renderFixMe(c, "FIXME:  get add memory")
}

func getAbout(c echo.Context) error {
	return c.Render(http.StatusOK, "about.html", nil)
}

func markDownBasic(args ...interface{}) template.HTML {
    s := blackfriday.MarkdownCommon([]byte(fmt.Sprintf("%s", args...)))
    return template.HTML(s)
}

func main() {
	e := echo.New()

	funcMap := template.FuncMap {
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
	e.PUT("/memories", updateMemory)
	e.GET("/memories/submitted", getMemorySubmitted)
	e.GET("/memories/approve/:uuid", approveMemory)
	e.GET("/memories/add", getAddMemory)
	e.GET("/about", getAbout)

	log.Println("Listening...")
	e.Run(standard.New(":4000"))
}
