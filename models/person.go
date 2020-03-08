package model

import (
	"gopkg.in/mgo.v2/bson"
)

type (
	Person struct {
		ID    bson.ObjectId `json:"id" bson:"_id,omitempty"`
		Email string        `json:"email" bson:"email"`
		Name  string        `json:"name,omitempty" bson:"name"`
	}
)
