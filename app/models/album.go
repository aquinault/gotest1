package models

import (
	"gopkg.in/mgo.v2/bson"
)
//type omit *struct {}

type Album struct {
	Id bson.ObjectId `json:"id,omitempty" bson:"_id"`
	Name string `json:"name,omitempty"`
  	UserId string `json:"userid,omitempty"`
	Images []string `json:"images,omitempty"`
}

type Code struct {
  Type string `json:"type,omitempty"`
  Msg string `json:"msg,omitempty"`
}

type PublicError struct {
	Code Code `json:"code,omitempty"`
}

type PublicAlbum struct {
  Data *Album `json:"data,omitempty"`
  Code Code `json:"code,omitempty"`
}

type PublicAlbums struct {
  Data *[]Album `json:"data,omitempty"`
  Code Code `json:"code,omitempty"`
}