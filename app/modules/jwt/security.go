package jwt

import (
        "github.com/revel/revel"
 		"fmt"
        "errors"
    )

// Extension du controlleur
type Security struct {
    *revel.Controller
}

// Verifie la validity du token dans le cookie et la session
func (c *Security) CheckToken() (bool, error) {
    fmt.Println("CheckToken")
    // regarde sur le token est dans le Header sinon dans le cookie
    var token string
    if len(c.Request.Header["Token"]) != 0 {
        token = c.Request.Header.Get("Token")
    } else {
        token = c.Session["Token"]
    }
    if token != "" {
        _, err := ParseLoginToken(token, look)
        return true, err
    }

    return false, errors.New("unknown jwt token")
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







