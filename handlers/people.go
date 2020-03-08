package handler

import (
	"net/http"

	"github.com/labstack/echo"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	model "github.com/m-russ/go-template/models"
)

// CreatePerson inserts a new person model into the collection.
func (h *Handler) CreatePerson(c echo.Context) (err error) {
	// Bind
	p := &model.Person{ID: bson.NewObjectId()}
	if err = c.Bind(p); err != nil {
		return
	}

	// Validate
	if p.Email == "" || p.Name == "" {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: "invalid email or name"}
	}

	// Save person
	db := h.DB.Clone()
	defer db.Close()
	if err = db.DB("go-template").C("people").Insert(p); err != nil {
		return
	}

	return c.JSON(http.StatusCreated, p)
}

//ReadPerson returns a person from the collection.
func (h *Handler) ReadPerson(c echo.Context) (err error) {

	// Bind route params.
	email := c.Param("email")
	name := c.Param("name")

	// Validate
	if email == "" || name == "" {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: "invalid email or name"}
	}

	p := new(model.Person)

	//Find Person
	db := h.DB.Clone()
	defer db.Close()
	if err = db.DB("go-template").C("people").Find(bson.M{"email": email, "name": name}).One(p); err != nil {
		if err == mgo.ErrNotFound {
			return &echo.HTTPError{Code: http.StatusNoContent, Message: "no person with that email and name"}
		}
		return
	}

	return c.JSON(http.StatusOK, p)
}

//UpdatePerson overwrites email and name in the collection.
func (h *Handler) UpdatePerson(c echo.Context) (err error) {

	// Bind form data and route param.
	p := new(model.Person)
	id := c.Param("id")
	if err = c.Bind(p); err != nil {
		return
	}

	// Update Person
	db := h.DB.Clone()
	defer db.Close()
	if err = db.DB("go-template").C("people").
		UpdateId(bson.ObjectIdHex(id), bson.M{"email": p.Email, "name": p.Name}); err != nil {
		if err == mgo.ErrNotFound {
			return echo.ErrNotFound
		}
	}
	return c.String(http.StatusOK, "Updated person successfully")
}

// RemovePerson deletes a person from the collection.
func (h *Handler) RemovePerson(c echo.Context) (err error) {

	//Bind route param.
	id := c.Param("id")

	if id == "" {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: "id cannot be null or empty"}
	}

	//Delete Person
	db := h.DB.Clone()
	defer db.Clone()
	if err = db.DB("go-template").C("people").RemoveId(bson.ObjectIdHex(id)); err != nil {
		if err == mgo.ErrNotFound {
			return echo.ErrNotFound
		}
	}

	return c.String(http.StatusOK, "Deleted person successfully")
}
