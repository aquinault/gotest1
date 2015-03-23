package models


type User struct {
  Id string 
  Username string
  Firstname string
  Lastname string
  Email string
  TwitterUid string
  FacebookUid string
  Password string `json:"password" bson:"password"`
}