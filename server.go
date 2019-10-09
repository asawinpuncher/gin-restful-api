package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client
var collection *mongo.Collection

// Trainer represent struct
type Trainer struct {
	Name      string
	Age       int
	Workplace string
}

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
		v1.GET("/:name", findTrainer)
		v1.GET("/", findTrainers)
		v1.PUT("/:name/:age", updateTrainer)
	}
	r.Run()

}

func connectMongo() *mongo.Client {
	// Set client options
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	// Connect to MongoDB
	cli, err := mongo.Connect(context.TODO(), clientOptions)
	client = cli

	collection = client.Database("TrainerDB").Collection("Trainers")

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

func findTrainer(c *gin.Context) {
	name := c.Param("name")

	filter := bson.D{{"name", name}}

	var result Trainer
	err := collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Get Success",
		"result":  result,
	})
}

func findTrainers(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	cur, err := collection.Find(ctx, bson.D{{}})
	if err != nil {
		log.Fatal(err)
	}

	var allTrainers []*Trainer

	for cur.Next(context.TODO()) {
		var trainerP Trainer
		err := cur.Decode(&trainerP)
		if err != nil {
			log.Fatal(err)
		}
		allTrainers = append(allTrainers, &trainerP)
	}

	defer cur.Close(context.TODO())

	c.JSON(http.StatusOK, gin.H{
		"message": "Get Success",
		"result":  allTrainers,
	})
}

func updateTrainer(c *gin.Context) {
	name := c.Param("name")
	age := c.Param("age")
	// if len(strings.TrimSpace(name)) == 0 {
	// 	c.JSON(http.StatusOK, gin.H{
	// 		"message": "Wrong Name Input",
	// 	})
	// }

	filter := bson.D{{"name", name}}

	newAge := bson.D{{"$set", bson.D{{"age", age}}}}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	res, err := collection.UpdateOne(ctx, filter, newAge)
	if err != nil {
		log.Fatal(err)
	}

	c.JSON(http.StatusOK, gin.H{
		"message":        "Get Success",
		"match count":    res.MatchedCount,
		"modified count": res.ModifiedCount,
	})

}
