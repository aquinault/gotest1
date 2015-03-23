package api

import (
		"github.com/revel/revel"
 		"gotest1/app/models"
 		"strconv"
 		"fmt"
)

type APIUsers struct {
	*revel.Controller
}

var users []models.User = []models.User{{"0", "John Doo"}, {"1", "Maria Luis"}}

func (c APIUsers) Me() revel.Result {
	fmt.Println("Me()")
//	user := models.User{0, "John Doo"}	
	return c.RenderJson(users[0])
}

func (c APIUsers) List() revel.Result {
	fmt.Println("List()")

	newUser := models.User{string(len(users)), "user" + string(strconv.Itoa(len(users)))}
	users = append(users, newUser)	

	//users := c.Init()
	return c.RenderJson(users)
}

func (c APIUsers) Get(index int) revel.Result {
	fmt.Println("Get(index string)")
	//users := c.Init()
	return c.RenderJson(users[index])
}




