package main

import (
	"flag"
	"log"

	"github.com/inchingforward/mnmnt/handlers"
	"github.com/inchingforward/mnmnt/models"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	_ "github.com/lib/pq"
)

func init() {
	var err error

	models.DB, err = sqlx.Connect("postgres", "user=monument dbname=monument sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	e := echo.New()

	debug := false
	flag.BoolVar(&debug, "debug", false, "true to enable debug")
	flag.Parse()

	log.Printf("debug: %v\n", debug)

	handlers.Setup(e, debug)

	log.Println("Listening...")
	e.Run(standard.New(":4000"))
}
