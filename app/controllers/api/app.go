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
	user, err := c.GetUser()
	if err == nil {
		username = user.Username
	}
	pagetitle := "Accueil"

	return c.Render(pagetitle, username)
}
