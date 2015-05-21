package api

import (
		"github.com/revel/revel"
        "gotest1/app/models"
        "gotest1/app/modules/mongo"
        "gotest1/app/modules/jwt"
 		"fmt"
 		"log"
        "net/http"
        "encoding/json"
        "io/ioutil"
        "gopkg.in/mgo.v2"
        "gopkg.in/mgo.v2/bson"
)

type Users struct {
	*revel.Controller
    mongo.Mongo
    jwt.Security
}

func (c Users) RenderPublicErrorJson(statusType string, statusMsg string) revel.Result {    
    result2 := models.PublicError{}
    result2.Code.Type = statusType
    result2.Code.Msg = statusMsg
    return c.RenderJson(result2)
}

func (c Users) RenderPublicUser2Json(user *models.User, statusType string, statusMsg string) revel.Result {
    result2 := models.PublicUser2{Data: user}
    result2.Code.Type = statusType
    result2.Code.Msg = statusMsg
    return c.RenderJson(result2)
}



func (c Users) internalError() revel.Result {
    c.Response.Status = http.StatusUnauthorized       
    return c.RenderError(&revel.Error{
        Title:       "Not authorized",
        Description: "Token not valid for url " + string(c.Request.RequestURI ),
    })
}

func (c Users) parseUserItem() (models.User, error) {
    body, err := ioutil.ReadAll(c.Request.Body)
    if err != nil {
        log.Fatal(err)
    }
    useritem := models.User{}
    err = json.Unmarshal([]byte(body), &useritem)
    return useritem, err
}

func (c Users) Me() revel.Result {
    fmt.Println("Me()")

    // Verification du Token si invalide, retourne un 401
    //
    res, err := c.CheckToken();
    if err != nil {
        c.internalError()
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
    // Verification du Token si invalide, retourne un 401
    //
    _, err := c.CheckToken();
    if err != nil {
        c.internalError()
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

    // Verification du Token si invalide, retourne un 401
    //
    _, err := c.CheckToken();
    if err != nil {
        c.internalError()
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
    err := c1.Find(bson.M{"username": username, "password" : password}).One(&result)
    if err != nil {
        return c.RenderPublicErrorJson("KO", "Authenticated: " + err.Error()) 
        //log.Fatal(err)
        /*return c.RenderError(&revel.Error{
            Title:       "Not authenticated",
            Description: "Username or Password are not valid for url " + string(c.Request.RequestURI ),
        })*/

    }

    fmt.Println("User:", result)

    signingKey, _ := revel.Config.String("app.signingKey")

    //tokenString := jwt.GenerateToken(username, signingKey)
    tokenString := jwt.GenerateToken(result, signingKey)
    fmt.Println("tokenString : ", tokenString)
    c.Session["Token"] = string(tokenString)

/*    result2 := models.PublicUser{User: &result, Token: tokenString}
    fmt.Println("User2:", result2)

    // Todo à encapsuler comme dans albums
    result2.Code.Type = "OK"
    result2.Code.Msg = "Authenticated successfull"

    return c.RenderJson(result2)
*/
    return c.RenderPublicUser2Json(&result,"OK", "Authenticated successfull")
}

func (c Users) Create(username string, firstname string, lastname string, email string, id bson.ObjectId, twitteruid string, facebookuid string, password string) revel.Result {
    fmt.Println("username:", username)
    fmt.Println("firstname:", firstname)
    fmt.Println("lastname:", lastname)
    fmt.Println("email:", email)
    fmt.Println("id:", id)
    fmt.Println("twitteruid:", twitteruid)
    fmt.Println("facebookuid:", facebookuid)
    fmt.Println("password:", password)

    c1 := c.MongoDatabase.C("users")

    user := models.User{id, username, firstname, lastname, email, twitteruid, facebookuid, password, ""}

    err := c1.Insert(&user)
    if err != nil {
        log.Fatal(err)
    }

    return c.RenderJson(user)
}


func (c Users) UpdateAvatar(id string, fid string) revel.Result {
    _, err := c.CheckToken();
    if err != nil {
        c.internalError()
    }

    c1 := c.MongoDatabase.C("users")

    err = c1.Update(bson.M{"_id": bson.ObjectIdHex(id)}, bson.M{"$set": bson.M{"avatarid": fid}})
    if err != nil {
        log.Fatal(err)
    }

    // Update Token with the avatar id
    result := models.User{}
    err = c1.Find(bson.M{"_id": bson.ObjectIdHex(id)}).One(&result)
    if err != nil {
        log.Fatal(err)
    }

    signingKey, _ := revel.Config.String("app.signingKey")
    tokenString := jwt.GenerateToken(result, signingKey)
    result2 := models.PublicUser{User: &result, Token: tokenString}
    c.Session["Token"] = string(tokenString)

    //return c.RenderJson("OK")
    return c.RenderJson(result2)
}

func (c Users) Update(id string) revel.Result {
    user, err := c.parseUserItem()
    if err != nil {
        return c.RenderText("Unable to parse the UserItem from JSON.")
    }

    c1 := c.MongoDatabase.C("users")

    err = c1.Update(bson.M{"_id": bson.ObjectIdHex(id)}, &user)
    if err != nil {
        log.Fatal(err)
    }
    return c.RenderJson(user)
}

func (c Users) Delete(id string) revel.Result {
    fmt.Println("id:", id)
    c1 := c.MongoDatabase.C("users")

    //user := models.User{id, username, firstname, lastname, email, twitteruid, facebookuid, password}

    err := c1.Remove(bson.M{"_id": bson.ObjectIdHex(id)})
    //err := c1.Remove(bson.M{"id": id})
    if err != nil {
        log.Fatal(err)
    }

    return c.RenderJson(bson.M{"id": id})
}


func (c Users) CreateUsers() revel.Result {
	fmt.Println("CreateUsers()")
	user1 := models.User{bson.NewObjectId(), "jdoo","John","Doo","john@doo","0","0","password",""}
	user2 := models.User{bson.NewObjectId(), "mluis","Maria","Luis","maria@luis","0","0","password",""}
	user3 := models.User{bson.NewObjectId(), "test1","firstn","lastn","test1@test1","0","0","password",""}

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
