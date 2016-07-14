package main

import (
	"bytes"
	"database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	_ "github.com/lib/pq"
	"html/template"
	"io"
	"log"
	"net/http"
	"time"
)

var db *sqlx.DB

type Template struct {
	templates *template.Template
}

type Memory struct {
	Id           uint64     `db:"id"`
	Title        string     `db:"title"`
	Details      string     `db:"details"`
	Latitude     float64    `db:"latitude"`
	Longitude    float64    `db:"longitude"`
	Author       string     `db:"author"`
	IsApproved   bool       `db:"is_approved"`
	ApprovalUuid sql.NullString     `db:"approval_uuid"`
	InsertedAt   time.Time  `db:"inserted_at"`
	UpdatedAt    time.Time  `db:"updated_at"`
}

func init() {
	var err error
	
	db, err = sqlx.Connect("postgres", "user=monument dbname=monument sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
}

func render(c echo.Context, message string) error {
	return c.String(http.StatusOK, message)
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func index(c echo.Context) error {
	var memories []*Memory
	err := db.Select(&memories, "select * from memory")

	if err == nil {
		var buffer bytes.Buffer
		
		buffer.WriteString("Memories:\n\n")
		
		for _, mem := range memories {
			buffer.WriteString(mem.Title)
			buffer.WriteString("\n")
		}
		
		return c.Render(http.StatusOK, "index.html", memories)
	} else {
		return render(c, err.Error())
	}
}

func getMemories(c echo.Context) error {
	return render(c, "FIXME:  render list of memories")
}

func getMemory(c echo.Context) error {
	id := c.Param("id")
	return render(c, "FIXME:  get memory "+id)
}

func createMemory(c echo.Context) error {
	return render(c, "FIXME:  create memory")
}

func updateMemory(c echo.Context) error {
	return render(c, "FIXME:  update memory")
}

func getMemorySubmitted(c echo.Context) error {
	return render(c, "FIXME:  get memory submitted")
}

func approveMemory(c echo.Context) error {
	return render(c, "FIXME:  approve memory")
}

func getAddMemory(c echo.Context) error {
	return render(c, "FIXME:  get add memory")
}

func getAbout(c echo.Context) error {
	return render(c, "FIXME:  get about")
}

func main() {
	e := echo.New()

	t := &Template{
		templates: template.Must(template.ParseGlob("templates/*.html")),
	}

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
