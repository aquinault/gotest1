package api

import (
		"github.com/revel/revel"
        "gotest1/app/models"
        "gotest1/app/modules/mongo"
        "gotest1/app/modules/jwt"
 		//"fmt"
 		"log"
        "net/http"
        "encoding/json"
        "io/ioutil"
        //"io"
        //"errors"
        //"encoding/json"
        //"gopkg.in/mgo.v2"
        "gopkg.in/mgo.v2/bson"
)

type Albums struct {
	*revel.Controller
    mongo.Mongo
    jwt.Security
}

func (c Albums) internalError() revel.Result {
    c.Response.Status = http.StatusUnauthorized       
    return c.RenderError(&revel.Error{
        Title:       "Not authorized",
        Description: "Token not valid for url " + string(c.Request.RequestURI ),
    })
}

func (c Albums) parseUserItem() (models.User, error) {
    body, err := ioutil.ReadAll(c.Request.Body)
    if err != nil {
        log.Fatal(err)
    }
    useritem := models.User{}
    err = json.Unmarshal([]byte(body), &useritem)
    return useritem, err
}

func (c Albums) SaveAlbum() revel.Result {
    // Verification du Token si invalide, retourne un 401
    //
    user, err := c.CheckToken();
    if err != nil {
        c.internalError()
    }

	c1 := c.MongoDatabase.C("albums")

	album := models.Album{"album1", user.Username}

	err = c1.Insert(&album)
	if err != nil {
		log.Fatal(err)
	}

	return c.RenderJson(album)
}

func (c Albums) GetAlbums() revel.Result {
    // Verification du Token si invalide, retourne un 401
    //
    _, err := c.CheckToken();
    if err != nil {
        c.internalError()
    }

	//
	//
	c1 := c.MongoDatabase.C("albums")
	results := []models.Album{}
	err = c1.Find(bson.M{}).All(&results)
	if err != nil {
		log.Fatal(err)
	}

	results2 := make([]models.PublicAlbum, len(results))
	for i := 0; i < len(results); i++ {
		//results2[i] = models.PublicUser{User: &results[i], Token: "tokenString"}
		results2[i] = models.PublicAlbum{Album: &results[i]}
	}

	return c.RenderJson(results2)
}
