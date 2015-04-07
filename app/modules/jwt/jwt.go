package jwt

import (
		"github.com/dgrijalva/jwt-go"
		//"github.com/nu7hatch/gouuid"
 		"gotest1/app/models"
 		"fmt"
 		"time"
)

func Test(myToken string) string {
	return string("parseLoginToken2")
}

//func GenerateToken(username string, signature string) string {
func GenerateToken(user models.User, signature string) string {
	fmt.Println("GenerateToken()")

/*	u4, err := uuid.NewV4()
	if err != nil {
		fmt.Println("error:", err)
		//return;
	}

	fmt.Println("uuid ", u4)
	user := models.User{ u4.String(), myuser.Username, myuser.Firstname, myuser.Lastname, "0", "0", "0", "0"}
*/

	// Create the token
    token := jwt.New(jwt.SigningMethodHS256)
    token.Header["kind"] = "login"
    // Set some claims
    token.Claims["username"] = user.Username
    token.Claims["firstname"] = user.Firstname
    token.Claims["lastname"] = user.Lastname
    token.Claims["email"] = user.Email
    token.Claims["id"] = user.Id
    token.Claims["foo"] = "bar"
    token.Claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
    // Sign and get the complete encoded token as a string

    tokenString, err := token.SignedString([]byte(signature))
	fmt.Println("err : ", err)
    if err != nil {
    	return string(err.Error())
	}
	
	fmt.Println("token : ", tokenString)

    return tokenString
}


func ParseLoginToken(myToken string, myLookupKey func(interface{}) (interface{}, error)) (models.User, error) {

	fmt.Println("jwt.parseLoginToken()")

	token, err := jwt.Parse(myToken, func(token *jwt.Token) (interface{}, error) {
		//fmt.Println(myToken)
		return myLookupKey(token.Header["kind"])
	})

	//fmt.Println("token", token)

	if token.Valid {
		fmt.Println("You look nice today")
		return models.User{ token.Claims["id"].(string), 
							token.Claims["username"].(string), 
							token.Claims["firstname"].(string), 
							token.Claims["lastname"].(string), 
							token.Claims["email"].(string), 
							"0", "0", "0"}, nil


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
