package database

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"myapp/config"
	"time"
)

type Connection struct {
	*mongo.Client
	EmailCollection *mongo.Collection
	*mongo.Database
}

func ConnectToDB() *Connection {
	u := config.Config("DB_URL")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(u))

	if err != nil {
		panic(err)
	}

	// create database
	db := client.Database("frekwent-emailer")

	// create collection
	emailCollection := db.Collection("emails")

	dbConnection := &Connection{
		client,
		emailCollection,
		db,
	}

	return dbConnection
}
