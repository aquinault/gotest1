package controllers

import (
		"github.com/revel/revel"
 		"myapp/app/models"
 		"fmt"
)

type Users struct {
	*revel.Controller
}

var users []models.User = []models.User{{0, "John Doo"}, {1, "Maria Luis"}}

func (c Users) Me() revel.Result {
	fmt.Println("Me()")
//	user := models.User{0, "John Doo"}	
	return c.RenderJson(users[0])
}

func (c Users) List() revel.Result {
	fmt.Println("List()")
	//users := c.Init()
	return c.RenderJson(users)
}

func (c Users) Get(id int) revel.Result {
	fmt.Println("Get(id int)")
	//users := c.Init()
	return c.RenderJson(users[id])
}




