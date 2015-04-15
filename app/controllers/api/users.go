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
        //"io"
        //"errors"
        //"encoding/json"
        "gopkg.in/mgo.v2"
        "gopkg.in/mgo.v2/bson"
        "crypto/rand"
        "bytes"
        "image"
        "image/jpeg"
        _ "image/png"
        "github.com/nfnt/resize"
)

type Users struct {
	*revel.Controller
    mongo.Mongo
    jwt.Security
}

func randString(n int) string {
    const alphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
    var bytes = make([]byte, n)
    rand.Read(bytes)
    for i, b := range bytes {
        bytes[i] = alphanum[b % byte(len(alphanum))]
    }
    return string(bytes)
}

/*Encode to base64*/
/*func encodeBase64Token(hexVal string) string {
    token := base64.URLEncoding.EncodeToString([]byte(hexVal))
    return token
}
*/
/*Decode from base64*/
/*func decodeBase64Token(token string) string {
    hexVal, err := base64.URLEncoding.DecodeString(token)
    if err != nil {
        return ""
    }
    return string(hexVal)
}
*/
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

func (c Users) SaveImage() revel.Result {
    // Verification du Token si invalide, retourne un 401
    //
    _, err := c.CheckToken();
    if err != nil {
        c.internalError()
    }


/*    unique_filename := "myfilename" + randString(10)
    fmt.Println("unique_filename ", unique_filename)
  */  
    // 1 fichier file-0
    file, handler, err := c.Request.FormFile("file-0")

    if err != nil {
        fmt.Println("request.FormFile", err)

        return c.RenderJson(map[string]string{
            "state": "ERROR", 
        })

    }
    defer file.Close();

    // Get the original filename
    filename := handler.Filename

    // Read the file into memory
    data, err := ioutil.ReadAll(file)
    // ... check err value for nil

    // Specify the Mongodb database
    my_db := c.MongoDatabase

    // Create the file in the Mongodb Gridfs instance
    my_file, err := my_db.GridFS("fs").Create(filename)
    // ... check err value for nil

    // Set the Meta Data
    //my_file.SetMeta(bson.M{"username": (*res).Username, "email": (*res).Email, "id": (*res).Id})


    // Decoding gives you an Image.
    // If you have an io.Reader already, you can give that to Decode 
    // without reading it into a []byte.
    original_image, _, err := image.Decode(bytes.NewReader(data))
    fmt.Println("Etape1 --------------")
    if err != nil {
        fmt.Println(err)
        return c.RenderJson(map[string]string{"state": "ERROR",})      
    }
    
    newImage := resize.Resize(160, 0, original_image, resize.Lanczos3)
    fmt.Println("Etape2 --------------")

    // Encode uses a Writer, use a Buffer if you need the raw []byte
    err = jpeg.Encode(my_file, newImage, nil)
    fmt.Println("Etape3 --------------")

    fmt.Println(err)
    // check err

    // Write the file to the Mongodb Gridfs instance
    //n, err := my_file.Write(data)
    // ... check err value for nil

    //encode file id and serve
    fileId := c.EncodeBase64Token(my_file.Id().(bson.ObjectId).Hex())

    // Close the file
    err = my_file.Close()
    // ... check err value for nil

    // Write a log type message
    //fmt.Println("%d bytes written to the Mongodb instance\n", n)

    //   
    //return c.RenderJson(fileId)
    return c.RenderJson(map[string]string{
            "fid": fileId, 
            //"size" : string(n), 
            "state": "SUCCESS", 
        })
}

func (c Users) GetImage(fid string) revel.Result {
    fmt.Println("fid: ", fid)    

    file_id := c.DecodeBase64Token(fid)

    // Verification du Token si invalide, retourne un 401
    //
    _, err := c.CheckToken();
    if err != nil {
        c.internalError()
    }

    // Specify the Mongodb database
    my_db := c.MongoDatabase

    //my_file, _ := my_db.GridFS("fs").Open(name)
    my_file, _ := my_db.GridFS("fs").OpenId(bson.ObjectIdHex(file_id))

    b := make([]byte, my_file.Size())
    my_file.Read(b)
    //fmt.Println(string(b))
    _ = my_file.Close()

    c.Response.Status = http.StatusOK
    c.Response.ContentType = "image/png"

    return c.RenderText(string(b))
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
        log.Fatal(err)
    }

    fmt.Println("User:", result)

    signingKey, _ := revel.Config.String("app.signingKey")

    //tokenString := jwt.GenerateToken(username, signingKey)
    tokenString := jwt.GenerateToken(result, signingKey)

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

    err = c1.Update(bson.M{"id": id}, bson.M{"$set": bson.M{"avatar": fid}})
    if err != nil {
        log.Fatal(err)
    }

    // Update Token with the avatar id
    result := models.User{}
    err = c1.Find(bson.M{"id": id}).One(&result)
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

    err = c1.Update(bson.M{"id": id}, &user)
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
	user1 := models.User{"1", "jdoo","John","Doo","john@doo","0","0","password",""}
	user2 := models.User{"2", "mluis","Maria","Luis","maria@luis","0","0","password",""}
	user3 := models.User{"3", "test1","firstn","lastn","test1@test1","0","0","password",""}

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
