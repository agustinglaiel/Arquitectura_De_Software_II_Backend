package db

import (
	"context"
	"fmt"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func InitDB() (*mongo.Database, error) {
  uri := os.Getenv("MONGODB_URI")
  if uri == "" {
    uri = "mongodb://root:root@localhost:27017"
  }

  clientOpts := options.Client().ApplyURI(uri)
  cli, err := mongo.Connect(context.TODO(), clientOpts)
  client = cli
  if err != nil {
    return nil, err
  }

  db := client.Database("arqui_de_software_2")

  dbNames, err := client.ListDatabaseNames(context.TODO(), bson.M{})
  if err != nil {
    return nil, err
  }

  fmt.Println("Available databases:")
  fmt.Println(dbNames)

  return db, nil
}

func DisconnectDB() {
  if client != nil {
    client.Disconnect(context.TODO())
  }
}
