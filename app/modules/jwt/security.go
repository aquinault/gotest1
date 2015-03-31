package jwt

import (
        "github.com/revel/revel"
        "gotest1/app/models"
 		"fmt"
        "errors"
    )

// Extension du controlleur
type Security struct {
    *revel.Controller
}

// Recupere le token dans le cookie, decode le token et recupere le user
func (c *Security) GetUser() (*models.User, error) {
    tokenString := c.Session["Token"]
    if tokenString == "" {
        return nil, errors.New("unknown jwt token")

    } else {
        user, err := ParseLoginToken(tokenString, look)
        return &user, err
    }

}

// Verifie la validity du token dans le cookie et la session
func (c *Security) CheckToken() (*models.User, error) {
    fmt.Println("CheckToken")
    // regarde sur le token est dans le Header sinon dans le cookie
    var token string
    if len(c.Request.Header["Token"]) != 0 {
        token = c.Request.Header.Get("Token")
    } else {
        token = c.Session["Token"]
    }
    if token != "" {
        user, err := ParseLoginToken(token, look)
        return &user, err
    }

    return nil, errors.New("unknown jwt token")
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






