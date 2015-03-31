package api

import (
		"github.com/revel/revel"
		//"github.com/dgrijalva/jwt-go"
		//"github.com/nu7hatch/gouuid"
 		//"myapp/app/models"
 		//"gotest1/app/utils"
 		"gotest1/app/modules/jwt"
 		"errors"
 		"fmt"
 		//"time"
)

type Auth struct {
	*revel.Controller
}

//const mySigningKey = "secret"


func (c Auth) Token(username string, signature string) revel.Result {
	fmt.Println("Token()")

	fmt.Println("username ", username)
	fmt.Println("signature ", signature)

	//return c.RenderText(jwt.GenerateToken(username, signature))
	return c.Render()
}


func (c Auth) TestToken() revel.Result {
	fmt.Println("TestToken()")

	myToken := "eyJhbGciOiJIUzI1NiIsImtpbmQiOiJsb2dpbiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE0MjcxMjMzNjMsImZvbyI6ImJhciIsImlkIjoiMGUzZWQ3NmItYmI3Ni00NTUyLTQ2ZDktNTkwOWU0NzcwMWE2IiwidXNlciI6ImFxdWluYXVsdCJ9.k-WxkoTV3Vo1ziFa_V8dobCLCksMIqT-f4TImvLQqoY"

	value, _ := jwt.ParseLoginToken(myToken, look)
	//fmt.Println("value", value)

	return c.RenderJson(value)
}

func look(kind interface{}) (interface{}, error) {
	signingKey, _ := revel.Config.String("app.signingKey")
	if str, ok := kind.(string); ok {
		switch str {
		case "login":
			return []byte(signingKey), nil
		}
	}

	return "", errors.New("unknown jwt kind")
}
