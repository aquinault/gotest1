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

type PublicAlbum struct {
  *Album
  StatusType string `json:"statustype,omitempty"`
  StatusMsg string `json:"statusmsg,omitempty"`
}
