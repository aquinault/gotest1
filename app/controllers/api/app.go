package api

import (
	"github.com/revel/revel"
	"gotest1/app/modules/jwt"
	)

type APIApp struct {
	*revel.Controller
	jwt.Security
}

const mySigningKey = "secret"

func (c APIApp) Index() revel.Result {
	var username string = ""
	var firstname string = ""
	var lastname string = ""
	var email string = ""
	user, err := c.GetUser()
	if err == nil {
		username = user.Username
		firstname = user.Firstname
		lastname = user.Lastname
		email = user.Email
	}
	pagetitle := "Accueil"


	return c.Render(pagetitle, username, firstname, lastname, email)
}
