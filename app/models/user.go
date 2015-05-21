package models

import (
  "gopkg.in/mgo.v2/bson"
)

type omit *struct {}

type User struct {
  //Id string `json:"id"`
  Id bson.ObjectId `json:"id" bson:"_id"`
  Username string `json:"username"`
  Firstname string `json:"firstname"`
  Lastname string `json:"lastname,omitempty"`
  Email string `json:"email"`
  TwitterUid string `json:"twitteruid,omitempty"`
  FacebookUid string `json:"facebookuid,omitempty"`
  Password string `json:"password"`
  AvatarId string `json:"avatarid,omitempty"`
}

type PublicUser struct {
  *User
  Password omit `json:"password,omitempty"`
  Token string `json:"token,omitempty"`
  Code Code `json:"code,omitempty"`
}

type PublicUser2 struct {
  Data *User `json:"data,omitempty"`
  Code Code `json:"code,omitempty"`
}