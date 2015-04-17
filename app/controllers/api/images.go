package api

import (
		"github.com/revel/revel"
        "gotest1/app/modules/mongo"
        "gotest1/app/modules/jwt"
        "gotest1/app/models"
 		"fmt"
        "net/http"
        "io/ioutil"
        "gopkg.in/mgo.v2/bson"
        //"crypto/rand"
        "bytes"
        "image"
        "image/jpeg"
        _ "image/png"
        "github.com/nfnt/resize"
)

type Images struct {
	*revel.Controller
    mongo.Mongo
    jwt.Security
}

func (c Images) RenderPublicErrorJson(statusType string, statusMsg string) revel.Result {    
    result2 := models.PublicError{}
    result2.Code.Type = statusType
    result2.Code.Msg = statusMsg
    return c.RenderJson(result2)
}

func (c Images) internalError() revel.Result {
    c.Response.Status = http.StatusUnauthorized       
    return c.RenderError(&revel.Error{
        Title:       "Not authorized",
        Description: "Token not valid for url " + string(c.Request.RequestURI ),
    })
}

func (c Images) SaveImage() revel.Result {
    // Verification du Token si invalide, retourne un 401
    //
    _, err := c.CheckToken();
    if err != nil {
        c.internalError()
    }

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
    my_file1, err := my_db.GridFS("fs").Create(filename)
    //my_file2, err := my_db.GridFS("fs").Create(filename)
    //my_file3, err := my_db.GridFS("fs").Create(filename)
    // ... check err value for nil

    // Set the Meta Data
    //my_file.SetMeta(bson.M{"username": (*res).Username, "email": (*res).Email, "id": (*res).Id})


    // Decoding gives you an Image.
    // If you have an io.Reader already, you can give that to Decode 
    // without reading it into a []byte.
    original_image, _, err := image.Decode(bytes.NewReader(data))
    if err != nil {
        return c.RenderPublicErrorJson("KO", "Save Image: " + err.Error()) 
    }
    
    newImage1 := resize.Resize(100, 0, original_image, resize.Lanczos3)
    //newImage2 := resize.Resize(200, 0, original_image, resize.Lanczos3)
    //newImage3 := resize.Resize(500, 0, original_image, resize.Lanczos3)

    // Encode uses a Writer, use a Buffer if you need the raw []byte
    err = jpeg.Encode(my_file1, newImage1, nil)
    if err != nil {
        return c.RenderPublicErrorJson("KO", "Save Image: " + err.Error()) 
    }
    //err = jpeg.Encode(my_file2, newImage2, nil)
    //err = jpeg.Encode(my_file3, newImage3, nil)

//    fmt.Println(err)
    // check err

    // Write the file to the Mongodb Gridfs instance
    //n, err := my_file.Write(data)
    // ... check err value for nil

    //encode file id and serve
    file1Id := c.EncodeBase64Token(my_file1.Id().(bson.ObjectId).Hex())
    //file2Id := c.EncodeBase64Token(my_file2.Id().(bson.ObjectId).Hex())
    //file3Id := c.EncodeBase64Token(my_file3.Id().(bson.ObjectId).Hex())

    // Close the file
    err = my_file1.Close()
    //err = my_file2.Close()
    //err = my_file3.Close()
    // ... check err value for nil

    // Write a log type message
    //fmt.Println("%d bytes written to the Mongodb instance\n", n)

    //   
    //return c.RenderJson(fileId)
    return c.RenderJson(map[string]string{
            "fid1": file1Id, 
            //"fid2": file2Id, 
            //"fid3": file3Id, 
            //"size" : string(n), 
            "state": "SUCCESS", 
        })
}

func (c Images) DeleteImage(fid string) revel.Result {
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
    _ = my_db.GridFS("fs").RemoveId(bson.ObjectIdHex(file_id))

    return c.RenderText("Remove OK")
}


func (c Images) GetImage(fid string) revel.Result {
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