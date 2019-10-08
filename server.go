package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client
var error *mongo.Client

type Trainer struct {
	Name      string
	Age       int
	Workplace string
}

// func main() {
// 	r := gin.Default()

// 	// Set client options
// 	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

// 	// Connect to MongoDB
// 	client, err := mongo.Connect(context.TODO(), clientOptions)

// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// Check the connection
// 	err = client.Ping(context.TODO(), nil)

// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	ash := Trainer{"Ash", 10, "Pallet Town"}
// 	//misty := Trainer{"Misty", 10, "Cerulean City"}

// 	collection := client.Database("TrainerDB").Collection("Trainers")

// 	insertResult, err := collection.InsertOne(context.TODO(), ash)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	r.GET("/ping", func(c *gin.Context) {
// 		c.JSON(http.StatusOK, gin.H{
// 			"message": insertResult.InsertedID,
// 		})

// 	})

// 	r.Run()
// }

func main() {
	r := gin.Default()

	client = connectMongo()
	// // Check the connection
	// err = client.Ping(context.TODO(), nil)

	// if err != nil {
	// 	log.Fatal(err)
	// }

	//misty := Trainer{"Misty", 10, "Cerulean City"}

	v1 := r.Group("/api/v1/trainers")
	{
		v1.POST("/", insertTrainer)
		v1.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "Test Success",
			})
		})
	}
	r.Run()

}

func connectMongo() *mongo.Client {
	// Set client options
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	// Connect to MongoDB
	cli, err := mongo.Connect(context.TODO(), clientOptions)
	client = cli

	if err != nil {
		log.Fatal(err)
	}
	return client
}

func insertTrainer(c *gin.Context) {
	name := c.PostForm("name")

	workplace := c.PostForm("workplace")

	fmt.Println(name)

	fmt.Println(workplace)
	age, err := strconv.Atoi(c.PostForm("age"))
	fmt.Println(age)
	fmt.Println(err)

	dataTrainer := Trainer{
		c.PostForm("name"), age, c.PostForm("workplace")}

	collection := client.Database("TrainerDB").Collection("Trainers")
	fmt.Println(collection)

	insertResult, err := collection.InsertOne(context.TODO(), dataTrainer)
	if err != nil {
		log.Fatal(err)
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Inserted Success",
		"resourceId": insertResult.InsertedID,
	})

}
