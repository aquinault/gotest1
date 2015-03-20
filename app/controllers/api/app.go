package api

import "github.com/revel/revel"

type APIApp struct {
	*revel.Controller
}

func (c APIApp) Index() revel.Result {
	greeting := "A"
	return c.Render(greeting)
	//return c.Render()
}