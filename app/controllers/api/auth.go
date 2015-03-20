package api

import (
		"github.com/revel/revel"
		//"github.com/dgrijalva/jwt-go"
		//"github.com/nu7hatch/gouuid"
 		//"myapp/app/models"
 		"myapp/app/utils"
 		"errors"
 		"fmt"
 		//"time"
)

type APIAuth struct {
	*revel.Controller
}

const mySigningKey = "secret"


func (c APIAuth) Token(username string, signature string) revel.Result {
	fmt.Println("Token()")

	fmt.Println("username ", username)
	fmt.Println("signature ", signature)

	return c.RenderText(utils.GenerateToken(username, signature))
}


func (c APIAuth) TestToken() revel.Result {
	fmt.Println("TestToken()")

	myToken := "eyJhbGciOiJIUzI1NiIsImtpbmQiOiJsb2dpbiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE0MjcxMjMzNjMsImZvbyI6ImJhciIsImlkIjoiMGUzZWQ3NmItYmI3Ni00NTUyLTQ2ZDktNTkwOWU0NzcwMWE2IiwidXNlciI6ImFxdWluYXVsdCJ9.k-WxkoTV3Vo1ziFa_V8dobCLCksMIqT-f4TImvLQqoY"

	value, _ := utils.ParseLoginToken(myToken, look)
	//fmt.Println("value", value)

	return c.RenderJson(value)
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
