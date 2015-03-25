package utils

import (
        "github.com/revel/revel"
 		//"fmt"
        "gopkg.in/mgo.v2"
        "sync"
        //"gopkg.in/mgo.v2/bson"

)

// Extension du controlleur
type Mongo struct {
    *revel.Controller
    MongoSession *mgo.Session
    MongoDatabase *mgo.Database
}

// Stockage globale de la session dont la visibilite est restreinte au package
var session *mgo.Session

// Singleton
var dial sync.Once

// Renvoie la session mgo en cours, si aucune n'existe, elle est créée
func GetSession() *mgo.Session {
    host, _ := revel.Config.String("mongo.host")

    // Appelé une seule fois grace au sync et de maniere synchrone
    dial.Do(func() {
        var err error
        session, err = mgo.Dial(host)
        if err != nil {
            panic(err)
        }
    })
    return session
}

// Alimente les proprietes affectées au controlleur en clonant la session mongo
func (c *Mongo) Bind() revel.Result {
    // Oublie pas de mettre mongo.database dans  app.conf, genre "localhost"
    databaseName, _ := revel.Config.String("mongo.database")
    c.MongoSession = GetSession().Clone()
    c.MongoDatabase = c.MongoSession.DB(databaseName)
    return nil
}

// Ferme un clone
func (c *Mongo) Close() revel.Result {
    if c.MongoSession != nil {
        c.MongoSession.Close()
    }
    return nil
}

// Fonction appelée au chargement de l'application
// Elle effectue un appel a notre fonction bind avant
// Chaque execution du controlleur
func init() {
    revel.InterceptMethod((*Mongo).Bind, revel.BEFORE)
    revel.InterceptMethod((*Mongo).Close, revel.AFTER)
    // On veut aussi fermer le clone si le controlleur plante
    revel.InterceptMethod((*Mongo).Close, revel.PANIC)
}








