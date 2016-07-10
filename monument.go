package main

import (
	"net/http"
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
)

func render(c echo.Context, message string) error {
	return c.String(http.StatusOK, message)
}

func index(c echo.Context) error {
	return render(c, "I am a fresh start for the Monument web app.")
}

func getMemories(c echo.Context) error {
	return render(c, "FIXME:  render list of memories")
}

func getMemory(c echo.Context) error {
	id := c.Param("id")
	return render(c, "FIXME:  get memory " + id)
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

	e.GET("/", index)
	e.GET("/memories", getMemories)
	e.GET("/memories/:id", getMemory)
	e.POST("/memories", createMemory)
	e.PUT("/memories", updateMemory)
	e.GET("/memories/submitted", getMemorySubmitted)
	e.GET("/memories/approve/:uuid", approveMemory)
	e.GET("/memories/add", getAddMemory)
	e.GET("/about", getAbout)

	e.Run(standard.New(":4000"))
}