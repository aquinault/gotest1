package api

import (
		"github.com/revel/revel"
 		"gotest1/app/models"
 		"strconv"
 		"fmt"
 		"log"
        "gopkg.in/mgo.v2"
        "gopkg.in/mgo.v2/bson"
)

type APIUsers struct {
	*revel.Controller
}

var users []models.User = []models.User{{"0", "John Doo", "0","0","0","0","0","0"}, {"1", "Maria Luis","0","0","0","0","0","0"}}

func (c APIUsers) Login() revel.Result {
    session, err := mgo.Dial("localhost")
    if err != nil {
        panic(err)
    }
    defer session.Close()

    // Optional. Switch the session to a monotonic behavior.
    session.SetMode(mgo.Monotonic, true)

    c1 := session.DB("test").C("users")
    /*err = c1.Insert(&user1, &user2, &user3)
    if err != nil {
            log.Fatal(err)
    }
	*/

	result := models.User{}
    err = c1.Find(bson.M{"username": "jdoo", "password" : "password"}).One(&result)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("User:", result)

	return c.RenderJson(result)
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




