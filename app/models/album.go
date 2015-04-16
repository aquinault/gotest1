package models

import (
	"gopkg.in/mgo.v2/bson"
)
//type omit *struct {}

type Album struct {
	Id bson.ObjectId `json:"id" bson:"_id"`
	Name string `json:"name"`
  	Username string `json:"username"`
	Images []string `json:"images"`
}

type PublicAlbum struct {
  *Album
}
