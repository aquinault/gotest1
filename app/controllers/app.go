package controllers

import "github.com/revel/revel"

type App struct {
	*revel.Controller
}

func (c App) Index() revel.Result {
	greeting := "A"
	return c.Render(greeting)
	//return c.Render()
}