package controllers

import (
		"github.com/revel/revel"
		"github.com/dgrijalva/jwt-go"
 		"myapp/app/models"
 		"errors"
 		"fmt"
 		"time"
)

type Auth struct {
	*revel.Controller
}

const mySigningKey = "secret"


func (c Auth) Token() revel.Result {
	fmt.Println("Token()")

	/*user := new(models.User)
	user.Id = 0
	user.Name = "John Doo"
	*/
	user := models.User{0, "John Doo"}


	// Create the token
    token := jwt.New(jwt.SigningMethodHS256)
    token.Header["kind"] = "login"
    // Set some claims
    token.Claims["user"] = user.Username
    token.Claims["id"] = user.Id
    token.Claims["foo"] = "bar"
    token.Claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
    // Sign and get the complete encoded token as a string
	fmt.Println("mySigningKey : ", mySigningKey)



    tokenString, err := token.SignedString([]byte(mySigningKey))
	fmt.Println("err : ", err)
    if err != nil {
    	return c.RenderText(err.Error())
	}
	
	fmt.Println("token : ", tokenString)

    return c.RenderText(tokenString)
}


func (c Auth) TestToken() revel.Result {
	fmt.Println("TestToken()")

	myToken := "eyJhbGciOiJIUzI1NiIsImtpbmQiOiJsb2dpbiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE0MjcwMzc2NzIsImZvbyI6ImJhciIsImlkIjowLCJ1c2VyIjoiSm9obiBEb28ifQ.kRX5jZ4C9GcDJE0vHX0Ezbs-F_-KjT2bHNGqbIrNy0c"

	value, _ := parseLoginToken(myToken, look)
	fmt.Println("value", value)

	return c.RenderJson(value)
}

func parseLoginToken(myToken string, myLookupKey func(interface{}) (interface{}, error)) (models.User, error) {
	token, err := jwt.Parse(myToken, func(token *jwt.Token) (interface{}, error) {
		fmt.Println(myToken)
		return myLookupKey(token.Header["kind"])
	})

	fmt.Println("token", token)

	if token.Valid {
		fmt.Println("You look nice today")
		return models.User{ int32(token.Claims["id"].(float64)), token.Claims["user"].(string)}, nil


	} else if ve, ok := err.(*jwt.ValidationError); ok {
		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			fmt.Println("That's not even a token")
			return models.User{}, err
		} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			// Token is either expired or not active yet
			fmt.Println("Timing is everything")
			return models.User{}, err
		} else {
			fmt.Println("Couldn't handle this token:", err)
			return models.User{}, err
		}
	} else {
		fmt.Println("Couldn't handle this token:", err)
		return models.User{}, err
	}


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
