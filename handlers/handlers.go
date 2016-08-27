package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"

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

// Template represents the parsed templates from the "templates" directory.
type Template struct {
	templates *template.Template
}

var (
	t     *Template
	debug = false
)

// Setup sets the memory handlers and renderer on the Echo instance.
func Setup(e *echo.Echo, isDebug bool) {
	debug = isDebug

	funcMap := template.FuncMap{
		"mdb": MarkDownBasic,
	}

	t = &Template{
		templates: template.Must(template.New("main").Funcs(funcMap).ParseGlob("templates/*.html")),
	}

	e.Static("/static", "static")
	e.GET("/", Index)
	e.GET("/memories", GetMemories)
	e.GET("/memories/:uuid/edit", GetEditMemory)
	e.POST("/memories/edit", EditMemory)
	e.GET("/memories/:slug", GetMemory)
	e.POST("/memories", CreateMemory)
	e.GET("/memories/submitted", GetMemorySubmitted)
	e.GET("/memories/approve/:uuid", ApproveMemory)
	e.GET("/memories/add", GetAddMemory)
	e.GET("/about", GetAbout)

	e.SetRenderer(t)
}

// MarkDownBasic passes the given data to the MarkdownBasic formatter.
func MarkDownBasic(args ...interface{}) template.HTML {
	s := blackfriday.MarkdownBasic([]byte(fmt.Sprintf("%s", args...)))
	return template.HTML(s)
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
	slug := c.Param("slug")

	memory, err := models.GetMemory(slug)

	return Render(c, "memory.html", memory, err)
}

// GetEditMemory renders the memory form using the memory
// for the given uuid parameter.
func GetEditMemory(c echo.Context) error {
	uuid := c.Param("uuid")
	if uuid == "" {
		return RenderMessage(c, http.StatusBadRequest, fmt.Sprintf("Missing uuid"))
	}

	memory, err := models.GetMemoryByEditUUID(uuid)
	if err == sql.ErrNoRows {
		return c.Render(http.StatusNotFound, "404.html", nil)
	}

	return RenderContext(c, "memory_edit_form.html", TemplateContext{memory, err})
}

// CreateMemory creates a new Memory using values from a submitted form.
func CreateMemory(c echo.Context) error {
	m := models.Memory{}
	if err := c.Bind(&m); err != nil {
		return Render(c, "memory.html", m, err)
	}

	err := models.AddMemory(&m)
	if err != nil {
		return RenderContext(c, "memory_form.html", TemplateContext{m, err})
	}

	utils.SendEmail(m)

	return c.Redirect(http.StatusFound, "/memories/submitted")
}

// EditMemory updates the details of a previously saved memory.
func EditMemory(c echo.Context) error {
	uuid := c.FormValue("uuid")

	if uuid == "" {
		return RenderMessage(c, http.StatusBadRequest, fmt.Sprintf("Missing uuid"))
	}

	memory, err := models.GetMemoryByEditUUID(uuid)
	if err == sql.ErrNoRows {
		return c.Render(http.StatusNotFound, "404.html", nil)
	}

	details := c.FormValue("details")
	if details == "" {
		err = errors.New("Missing details")
		return RenderContext(c, "memory_edit_form.html", TemplateContext{memory, err})
	}
	memory.Details = details

	err = models.UpdateDetails(memory)
	if err != nil {
		return RenderContext(c, "memory_edit_form.html", TemplateContext{memory, err})
	}

	url := fmt.Sprintf("/memories/%v", memory.ID)
	return c.Redirect(http.StatusFound, url)
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

	memory, err := models.GetMemoryByApprovalUUID(uuid)
	if err != nil {
		return Render(c, "", nil, err)
	}

	models.ApproveMemory(memory)
	if err != nil {
		return Render(c, "", memory, err)
	}

	utils.Tweet(memory)

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
