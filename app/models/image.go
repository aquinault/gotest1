package models

import (
	"gopkg.in/mgo.v2/bson"
)
/*
"small":  {"w": 340, "h": 358},
"medium": {"w": 600, "h": 631},
"large":  {"w": 913, "h": 960},
"thumb":  {"w": 150, "h": 150}
*/
type Image struct {
  Origin string  `json:"origin,omitempty"`     // "MGO", "S3"
  Small bson.ObjectId `json:"small,omitempty"`
  Medium bson.ObjectId `json:"medium,omitempty"`
  Large bson.ObjectId `json:"large,omitempty"`
  Thumb bson.ObjectId `json:"thumb,omitempty"`
}

type PublicImage struct {
  Data *Image `json:"data,omitempty"`
  Code Code `json:"code,omitempty"`
}
