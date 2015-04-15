package ihm

import (
		"github.com/revel/revel"
 		"gotest1/app/modules/jwt"
 		//"gotest1/app/utils"
 		"errors"
 		"fmt"
)

const mySigningKey = "secret"

type Authentication struct {
	*revel.Controller
	jwt.Security
}

func (c Authentication) Albums() revel.Result {
	fmt.Println("Albums")
	user, _ := c.GetUser()
	pagetitle := "Albums"
	return c.Render(pagetitle, user)
}

func (c Authentication) Album() revel.Result {
	fmt.Println("Album")
	user, _ := c.GetUser()
	pagetitle := "Album"
	return c.Render(pagetitle, user)
}

func (c Authentication) Login() revel.Result {
	fmt.Println("Login()")

	user, _ := c.GetUser()
	pagetitle := "Login"

	return c.Render(pagetitle, user)
}

func (c Authentication) Logout() revel.Result {
	fmt.Println("Logout()")
	c.Session["Token"] = ""
	return c.Redirect("/login")
}

func (c Authentication) UsersLogin() revel.Result {
	fmt.Println("UsersLogin()")
	return c.Render()
}

func (c Authentication) UsersCreate() revel.Result {
	fmt.Println("UsersCreate()")

	user, _ := c.GetUser()
	pagetitle := "Users Management"

	return c.Render(pagetitle, user)
}

func (c Authentication) UsersList() revel.Result {
	fmt.Println("UsersCreate()")

	user, _ := c.GetUser()
	pagetitle := "Users Management"

	return c.Render(pagetitle, user)
}

func (c Authentication) UsersMe() revel.Result {
	fmt.Println("UsersMe()")

	user, _ := c.GetUser()
	pagetitle := "Me"

	return c.Render(pagetitle, user)
}

func (c Authentication) Testtoken(token string) revel.Result {
	fmt.Println("Testtoken() ", token)
	greeting := "Test Token"

	if token != "" {
		myToken := token
		json, _ := jwt.ParseLoginToken(myToken, look)		
		return c.Render(greeting, json)
	}
	return c.Render(greeting)
}

func (c Authentication) Token(username string, signature string, token string) revel.Result {
	fmt.Println("Token()")
	greeting := "Generate Token"
	
	if username != ""  && signature != "" {

		//tokenString := jwt.GenerateToken(username, signature)
		tokenString := "123"
		fmt.Println("tokenString : ", tokenString)
		return c.Render(greeting, username, signature, tokenString)	
	}
	return c.Render(greeting, username, signature, token)
}

func look(kind interface{}) (interface{}, error) {
	if str, ok := kind.(string); ok {
		switch str {
		case "login":
			return []byte(mySigningKey), nil
		}
	}
	return "", errors.New("unknown jwt kind")
}

