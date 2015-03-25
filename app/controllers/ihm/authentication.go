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
}

func (c Authentication) Login() revel.Result {
	fmt.Println("Login()")
	return c.Render()
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

		tokenString := jwt.GenerateToken(username, signature)

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

