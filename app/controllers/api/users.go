package api

import (
		"github.com/revel/revel"
        "gotest1/app/models"
        "gotest1/app/modules/mongo"
        "gotest1/app/modules/jwt"
 		"strconv"
 		"fmt"
 		"log"
        "net/http"
        //"errors"
        //"encoding/json"
        "gopkg.in/mgo.v2"
        "gopkg.in/mgo.v2/bson"
)

type APIUsers struct {
	*revel.Controller
    mongo.Mongo
    jwt.Security
}

var users []models.User = []models.User{{"0", "John Doo", "0","0","0","0","0","0"}, {"1", "Maria Luis","0","0","0","0","0","0"}}

func (c APIUsers) List2() revel.Result {
    fmt.Println("List2()")

    // Verification du Token
    // Si invalide, retourne un 401
    //
    res, err := c.CheckToken();
    fmt.Println(err) 
    if res {
       fmt.Println("Token OK") 
    } else {
        fmt.Println("Token KO")       
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
    for i := 0; i < 3; i++ {
        //results2[i] = models.PublicUser{User: &results[i], Token: "tokenString"}
        results2[i] = models.PublicUser{User: &results[i]}
    }

    return c.RenderJson(results2)
}

func (c APIUsers) Login(username string, password string) revel.Result {
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

func (c APIUsers) CreateUsers() revel.Result {
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


func (c APIUsers) Me() revel.Result {
	fmt.Println("Me()")
//	user := models.User{0, "John Doo"}	
	return c.RenderJson(users[0])
}

func (c APIUsers) List() revel.Result {
	fmt.Println("List()")

	newUser := models.User{string(len(users)), "user" + string(strconv.Itoa(len(users))),"0","0","0","0","0","0"}
	users = append(users, newUser)	

	//users := c.Init()
	return c.RenderJson(users)
}

func (c APIUsers) Get(index int) revel.Result {
	fmt.Println("Get(index string)")
	//users := c.Init()
	return c.RenderJson(users[index])
}




