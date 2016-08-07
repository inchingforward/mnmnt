package handlers

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/inchingforward/mnmnt/models"
	"github.com/inchingforward/mnmnt/utils"
	"github.com/labstack/echo"
	"github.com/russross/blackfriday"
)

// A TemplateContext holds data and an error that are passed in to a
// template for rendering.  Either Data or Err can be nil.
type TemplateContext struct {
	Data interface{}
	Err  error
}

// MarkDownBasic passes the given data to the MarkdownBasic formatter.
func MarkDownBasic(args ...interface{}) template.HTML {
	s := blackfriday.MarkdownBasic([]byte(fmt.Sprintf("%s", args...)))
	return template.HTML(s)
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

// RenderMessage displays message using the message.html template.
func RenderMessage(c echo.Context, status int, message string) error {
	return c.Render(status, "message.html", message)
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
