package api

import (
		"github.com/revel/revel"
        "gotest1/app/models"
        "gotest1/app/modules/mongo"
        "gotest1/app/modules/jwt"
 		"fmt"
 		"log"
        "net/http"
        //"errors"
        //"encoding/json"
        "gopkg.in/mgo.v2"
        "gopkg.in/mgo.v2/bson"
)

type Users struct {
	*revel.Controller
    mongo.Mongo
    jwt.Security
}

func (c Users) Me() revel.Result {
    fmt.Println("Me()")

    // Verification du Token
    // Si invalide, retourne un 401
    //
    res, err := c.CheckToken();
    if res == nil {
        c.Response.Status = http.StatusUnauthorized       
        return c.RenderError(&revel.Error{
            Title:       "Not authorized",
            Description: "Token not valid for url " + string(c.Request.RequestURI ),
        })
    }

    //
    //
    c1 := c.MongoDatabase.C("users")
    result := models.User{}
    err = c1.Find(bson.M{"username": (*res).Username}).One(&result)
    if err != nil {
        log.Fatal(err)
    }
    
    return c.RenderJson(result)
}

func (c Users) Get(id string) revel.Result {
    //fmt.Println("Get2(index string)")
    // Verification du Token
    // Si invalide, retourne un 401
    //
    res, err := c.CheckToken();
    if res == nil {
        c.Response.Status = http.StatusUnauthorized       
        return c.RenderError(&revel.Error{
            Title:       "Not authorized",
            Description: "Token not valid for url " + string(c.Request.RequestURI ),
        })
    }

    //
    //
    c1 := c.MongoDatabase.C("users")
    result := models.User{}
    err = c1.Find(bson.M{"id": id}).One(&result)
    if err != nil {
        log.Fatal(err)
    }
    
    return c.RenderJson(result)
}


func (c Users) List() revel.Result {
    fmt.Println("List2()")

    // Verification du Token
    // Si invalide, retourne un 401
    //
    res, err := c.CheckToken();
    if res == nil {
        c.Response.Status = http.StatusUnauthorized       
        return c.RenderError(&revel.Error{
            Title:       "Not authorized",
            Description: "Token not valid for url " + string(c.Request.RequestURI ),
        })
    }
    
    //
    //
    c1 := c.MongoDatabase.C("users")
    results := []models.User{}
    err = c1.Find(bson.M{}).All(&results)
    if err != nil {
        log.Fatal(err)
    }

    results2 := make([]models.PublicUser, len(results))
    for i := 0; i < len(results); i++ {
        //results2[i] = models.PublicUser{User: &results[i], Token: "tokenString"}
        results2[i] = models.PublicUser{User: &results[i]}
    }

    return c.RenderJson(results2)
}

func (c Users) Login(username string, password string) revel.Result {
    fmt.Println("username:", username)
    fmt.Println("password:", password)

    c1 := c.MongoDatabase.C("users")

    result := models.User{}
    //err = c1.Find(bson.M{"username": "jdoo", "password" : "password"}).One(&result)
    err := c1.Find(bson.M{"username": username, "password" : password}).One(&result)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("User:", result)

    signingKey, _ := revel.Config.String("app.signingKey")

    tokenString := jwt.GenerateToken(username, signingKey)
    fmt.Println("tokenString : ", tokenString)

    result2 := models.PublicUser{User: &result, Token: tokenString}

    fmt.Println("User2:", result2)

    c.Session["Token"] = string(tokenString)

    return c.RenderJson(result2)
}

func (c Users) Create(username string, firstname string, lastname string, email string, id string, twitteruid string, facebookuid string, password string) revel.Result {
    fmt.Println("username:", username)
    fmt.Println("firstname:", firstname)
    fmt.Println("lastname:", lastname)
    fmt.Println("email:", email)
    fmt.Println("id:", id)
    fmt.Println("twitteruid:", twitteruid)
    fmt.Println("facebookuid:", facebookuid)
    fmt.Println("password:", password)

    c1 := c.MongoDatabase.C("users")

    user := models.User{id, username, firstname, lastname, email, twitteruid, facebookuid, password}

    err := c1.Insert(&user)
    if err != nil {
        log.Fatal(err)
    }

    return c.RenderJson(user)
}

func (c Users) Delete(id string) revel.Result {
    fmt.Println("id:", id)
    c1 := c.MongoDatabase.C("users")

    //user := models.User{id, username, firstname, lastname, email, twitteruid, facebookuid, password}

    err := c1.Remove(bson.M{"id": id})
    if err != nil {
        log.Fatal(err)
    }

    return c.RenderJson(bson.M{"id": id})
}


func (c Users) CreateUsers() revel.Result {
	fmt.Println("CreateUsers()")
	user1 := models.User{"1", "jdoo","John","Doo","john@doo","0","0","password"}
	user2 := models.User{"2", "mluis","Maria","Luis","maria@luis","0","0","password"}
	user3 := models.User{"3", "test1","firstn","lastn","test1@test1","0","0","password"}

	var users [3]models.User
	users[0] = user1
	users[1] = user2
	users[2] = user3

    session, err := mgo.Dial("localhost")
    if err != nil {
        panic(err)
    }
    defer session.Close()

    // Optional. Switch the session to a monotonic behavior.
    session.SetMode(mgo.Monotonic, true)

    c1 := session.DB("test").C("users")
    err = c1.Insert(&user1, &user2, &user3)
    if err != nil {
            log.Fatal(err)
    }

	return c.RenderJson(users)
}
