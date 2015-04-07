package api

import (
	"github.com/revel/revel"
	"gotest1/app/modules/jwt"
	)

type App struct {
	*revel.Controller
	jwt.Security
}

const mySigningKey = "secret"

func (c App) Index() revel.Result {
	user, _ := c.GetUser()
	pagetitle := "Accueil"

	return c.Render(pagetitle, user)
}
