package main

import (
	"os"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
	"gopkg.in/mgo.v2"

	handler "github.com/m-russ/go-template/handlers"
)

func main() {
	e := echo.New()

	mongoDBUri := "localhost"
	env, isSet := os.LookupEnv("MONGODB_URI")
	if isSet {
		mongoDBUri = env
	}

	// Logger Config
	e.Logger.SetLevel(log.INFO)
	e.Use(middleware.Logger())

	// Database connection
	db, err := mgo.Dial(mongoDBUri)
	if err != nil {
		e.Logger.Fatal(err)
	}

	// Create indicies
	if err = db.Copy().DB("go-template").C("people").EnsureIndex(mgo.Index{
		Key:    []string{"email"},
		Unique: true,
	}); err != nil {
		log.Fatal(err)
	}

	// Initialize handler
	h := &handler.Handler{DB: db}

	// Routes
	e.POST("/person", h.CreatePerson)
	e.GET("/person/:email/:name", h.ReadPerson)
	e.PUT("/person/:id", h.UpdatePerson)
	e.DELETE("/person/:id", h.RemovePerson)

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}
