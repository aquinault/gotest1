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

func (c Albums) parseAlbumItem() (models.Album, error) {
    body, err := ioutil.ReadAll(c.Request.Body)
    if err != nil {
        log.Fatal(err)
    }
    albumitem := models.Album{}
    err = json.Unmarshal([]byte(body), &albumitem)
    return albumitem, err
}

func (c Albums) RenderPublicAlbumJson(album *models.Album, statusType string, statusMsg string) revel.Result {
	result2 := models.PublicAlbum{Album: album}
	result2.StatusType = statusType
	result2.StatusMsg = statusMsg
	return c.RenderJson(result2)
}

func (c Albums) DeleteAlbumImage(id string, fid string) revel.Result {
	// Verification du Token si invalide, retourne un 401
	//
	_, err := c.CheckToken();
	if err != nil {
		c.internalError()
	}

	c1 := c.MongoDatabase.C("albums")

    // Find the album Before Delete an image
	result := models.Album{}
	err = c1.Find(bson.M{"_id": bson.ObjectIdHex(id)}).One(&result)
	if err != nil {
		return c.RenderPublicAlbumJson(&result,"KO", "Delete Album Image: " + err.Error())
	}

	err = c1.Update(bson.M{"_id": bson.ObjectIdHex(id)}, bson.M{"$pull": bson.M{"images": fid}})
	if err != nil {
		//log.Fatal(err)
		return c.RenderPublicAlbumJson(&result,"KO", "Delete Album Image: " + err.Error())
	}

    // Find the album After Delete an image
	//result := models.Album{}
	err = c1.Find(bson.M{"_id": bson.ObjectIdHex(id)}).One(&result)
	if err != nil {
		return c.RenderPublicAlbumJson(&result,"KO", "Delete Album Image: " + err.Error())
	}

//	return c.RenderJson(user)
	return c.RenderPublicAlbumJson(&result,"OK", "Delete Album Image")

}

func (c Albums) UpdateAlbum(id string) revel.Result {
	// Verification du Token si invalide, retourne un 401
	//
	_, err := c.CheckToken();
	if err != nil {
		c.internalError()
	}

    album, err := c.parseAlbumItem()
    if err != nil {
        //return c.RenderText("Unable to parse the UserItem from JSON.")
        return c.RenderPublicAlbumJson(&album,"KO", "Update Album: " + err.Error())
	}
   

	c1 := c.MongoDatabase.C("albums")

    err = c1.Update(bson.M{"_id": bson.ObjectIdHex(id)}, &album)
    if err != nil {
        //log.Fatal(err)
        return c.RenderPublicAlbumJson(&album,"KO", "Update Album: " + err.Error())
    }

	//return c.RenderJson(album)
	return c.RenderPublicAlbumJson(&album,"OK", "Update Album")
}

func (c Albums) UpdateAlbumImage(id string, fid string) revel.Result {
	// Verification du Token si invalide, retourne un 401
	//
	_, err := c.CheckToken();
	if err != nil {
		c.internalError()
	}

	c1 := c.MongoDatabase.C("albums")

    // Find the album Before Update
	result := models.Album{}
	err = c1.Find(bson.M{"_id": bson.ObjectIdHex(id)}).One(&result)
	if err != nil {
		return c.RenderPublicAlbumJson(&result,"KO", "Update Album Image: " + err.Error())
	}


	err = c1.Update(bson.M{"_id": bson.ObjectIdHex(id)}, bson.M{"$push": bson.M{"images": fid}})
	if err != nil {
		log.Fatal(err)
	}

    // Find the album After Update
	//result = models.Album{}
	err = c1.Find(bson.M{"_id": bson.ObjectIdHex(id)}).One(&result)
	if err != nil {
		return c.RenderPublicAlbumJson(&result,"KO", "Update Album Image: " + err.Error())
	}


	//return c.RenderJson(user)
	return c.RenderPublicAlbumJson(&result,"OK", "Update Album Images")
}

func (c Albums) GetAlbumImages(id string) revel.Result {
    // Verification du Token si invalide, retourne un 401
    _, err := c.CheckToken();
    if err != nil {
        c.internalError()
    }

	c1 := c.MongoDatabase.C("albums")
	result := models.Album{}
	err = c1.Find(bson.M{"_id": bson.ObjectIdHex(id)}).One(&result)
	if err != nil {
		return c.RenderPublicAlbumJson(&result,"KO", "Get Album Images: " + err.Error())
	}

	return c.RenderPublicAlbumJson(&result,"OK", "Get Album Images")
}

func (c Albums) DeleteAlbum(id string) revel.Result {
    // Verification du Token si invalide, retourne un 401
    _, err := c.CheckToken();
    if err != nil {
        c.internalError()
    }

    c1 := c.MongoDatabase.C("albums")

    // Find the album
	result := models.Album{}
	err = c1.Find(bson.M{"_id": bson.ObjectIdHex(id)}).One(&result)
	if err != nil {
		return c.RenderPublicAlbumJson(&result,"KO", "Delete Album: " + err.Error())
	}

    // Remove the album
    err = c1.Remove(bson.M{"_id": bson.ObjectIdHex(id)})
    if err != nil {
        return c.RenderPublicAlbumJson(&result,"KO", "Delete Album: " + err.Error())
    }

    return c.RenderPublicAlbumJson(&result,"OK", "Delete Album")
}


func (c Albums) SaveAlbum() revel.Result {
    // Verification du Token si invalide, retourne un 401
    user, err := c.CheckToken();
    if err != nil {
        c.internalError()
    }

	c1 := c.MongoDatabase.C("albums")

	album := models.Album{bson.NewObjectId(), "album1", user.Id.Hex(), []string{}}

	err = c1.Insert(&album)
	if err != nil {
		//log.Fatal(err)
		return c.RenderPublicAlbumJson(&album,"KO", "Save Album : " + err.Error())
	}

	/*
	err = c1.Update(bson.M{"name": "album1"}, bson.M{"$push": bson.M{"chapters": "A"}})
	if err != nil {
		log.Fatal(err)
	}
	*/

//	return c.RenderJson(album)
    return c.RenderPublicAlbumJson(&album,"OK", "Save Album")
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
