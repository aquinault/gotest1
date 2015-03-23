package ihm

import (
		"github.com/revel/revel"
		//"github.com/dgrijalva/jwt-go"
		//"github.com/nu7hatch/gouuid"
 		//"myapp/app/models"
 		"gotest1/app/utils"
 		"errors"
 		"fmt"
 		//"time"
)

type IHMAuth struct {
	*revel.Controller
}

const mySigningKey = "secret"

func (c IHMAuth) Testtoken(token string) revel.Result {
	fmt.Println("Testtoken() ", token)
	greeting := "Test Token"

	if token != "" {
		myToken := token
		json, _ := utils.ParseLoginToken(myToken, look)		
		return c.Render(greeting, json)
	}

	return c.Render(greeting)
}

func (c IHMAuth) Token(username string, signature string, token string) revel.Result {

	fmt.Println("Token()")
	greeting := "Generate Token"

	
	if username != ""  && signature != "" {

		tokenString := utils.GenerateToken(username, signature)
		/*
		u4, err := uuid.NewV4()
		if err != nil {
			fmt.Println("error:", err)
			//return
		}
		user := models.User{ u4.String(), username}

		// Create the token
		token := jwt.New(jwt.SigningMethodHS256)
		token.Header["kind"] = "login"
		// Set some claims
		token.Claims["user"] = user.Username
		token.Claims["id"] = user.Id
		token.Claims["foo"] = "bar"
		token.Claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
		// Sign and get the complete encoded token as a string

		tokenString, err := token.SignedString([]byte(mySigningKey))
		fmt.Println("err : ", err)
		if err != nil {
			return c.RenderText(err.Error())
		}
		*/
		fmt.Println("tokenString : ", tokenString)
		return c.Render(greeting, username, signature, tokenString)	
	}

	return c.Render(greeting, username, signature, token)
	//return c.Render()
}

/*
func parseLoginToken(myToken string, myLookupKey func(interface{}) (interface{}, error)) (models.User, error) {
	token, err := jwt.Parse(myToken, func(token *jwt.Token) (interface{}, error) {
		fmt.Println(myToken)
		return myLookupKey(token.Header["kind"])
	})

	fmt.Println("token", token)

	if token.Valid {
		fmt.Println("You look nice today")
		return models.User{ token.Claims["id"].(string), token.Claims["user"].(string)}, nil


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
*/
func look(kind interface{}) (interface{}, error) {
	if str, ok := kind.(string); ok {
		switch str {
		case "login":
			return []byte(mySigningKey), nil
		}
	}

	return "", errors.New("unknown jwt kind")
}
