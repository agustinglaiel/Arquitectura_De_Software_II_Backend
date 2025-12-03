package db

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoDb *mongo.Database
var client *mongo.Client

func InitDB() error {
	/*uri := os.Getenv("localhost:27017")
	if uri == "" {
		uri = "mongodb://root:root@mongo:27017"
	}*/

	clientOpts := options.Client().ApplyURI("mongodb://root:root@localhost:27017")
	cli, err := mongo.Connect(context.TODO(), clientOpts)
	client = cli

	if err != nil {
		clientOpts22 := options.Client().ApplyURI("mongodb://root:root@localhost:27017")
		cli2, err2 := mongo.Connect(context.TODO(), clientOpts22)
		client=cli2		
		if err2!=nil{
			return err2
		}
	}


	MongoDb = client.Database("arqui_de_software_2")

	dbNames, err3 := client.ListDatabaseNames(context.TODO(), bson.M{})
	if err3 != nil {
		log.Println(err3)
		return err
	}

	fmt.Println("Available databases:")
	fmt.Println(dbNames)

	return nil
}

func DisconnectDB() {
	if client != nil {
		client.Disconnect(context.TODO())
	}
}
