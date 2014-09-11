package main

import (
	"github.com/gin-gonic/gin"
	"labix.org/v2/mgo"
)

func main() {

	// ===============================
	// # 1 = Routing and handler funcs
	// ===============================
	r := gin.Default()
	r.GET("/hello", func(c *gin.Context) {
		c.String(200, "hello world")
	})
	// !!!! remember r.Run(":8080")

	// ==============================
	// # 2 = Deserializing POST Pasta
	// ==============================
	type PastaPing struct {
		Name        string `json:"name"`
		CookingTime int    `json:"cookingTime"`
	}
	r.POST("/pasta-pong", func(c *gin.Context) {
		postedPasta := PastaPing{}
		c.Bind(&postedPasta)
		c.JSON(200, postedPasta)
	})
	// SHOW:
	// POST /pasta-pong
	// {
	//    "name": "fusilli",
	//    "cookingTime": 8
	// }

	// ========================
	// # 3 = Persisting a Pasta
	// ========================
	const (
		databaseAddress string = "localhost"
		databaseName    string = "test"
		collectionName  string = "pastas"
	)
	type Pasta struct {
		Name        string `json:"name" bson:"name"`
		CookingTime int    `json:"cookingTime" bson:"cookingTime,omitempty"`
	}
	mgoSession, _ := mgo.Dial(databaseAddress)
	defer mgoSession.Close()
	r.POST("/pasta", func(c *gin.Context) {
		dbSession := mgoSession.Copy()
		defer dbSession.Close()

		newPasta := Pasta{}
		if c.Bind(&newPasta) {
			dbSession.DB(databaseName).C(collectionName).Insert(newPasta)
			c.JSON(200, gin.H{"status": "ok"})
		}
	})
	// SHOW:
	// POST /pasta
	// {
	//    "name": "orecchiette",
	//    "cookingTime": 12
	// }
	//
	// $ mongo
	// > show collections
	// > db.pastas.find()
	// > db.pastas.drop()

	// ====================
	// # 4 = getting Params
	// ====================
	r.GET("/hello/:name", func(c *gin.Context) {
		name := c.Params.ByName("name")
		c.String(200, "hi there, "+name+"!")
	})

	// ===============================
	// # 5 = retriving a Pasta by name
	// ===============================
	r.GET("/pasta/:name", func(c *gin.Context) {
		name := c.Params.ByName("name")

		dbSession := mgoSession.Copy()
		defer dbSession.Close()
		coll := dbSession.DB(databaseName).C(collectionName)

		query := Pasta{
			Name: name,
		}
		foundPasta := Pasta{}
		err := coll.Find(query).One(&foundPasta)
		if err != nil {
			c.JSON(404, gin.H{"status": err.Error()})
		} else {
			c.JSON(200, foundPasta)
		}
	})

	// ===========================
	// # 6 = retrieving all Pastas
	// ===========================
	r.GET("/pasta", func(c *gin.Context) {
		dbSession := mgoSession.Copy()
		defer dbSession.Close()

		coll := dbSession.DB(databaseName).C(collectionName)

		var allPastas []Pasta
		err := coll.Find(nil).All(&allPastas)
		if err != nil {
			panic(err)
		} else {
			c.JSON(200, allPastas)
		}
	})

	// ===========
	// running the server....
	r.Run(":8080")

}
