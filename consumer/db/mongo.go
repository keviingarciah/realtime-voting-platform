package db

import (
	"context"
	"log"

	models "consumer/models"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var clientMongo *mongo.Client

func ConnectMongo() (*mongo.Client, error) {
	// Establece las opciones del cliente
	clientOptions := options.Client().ApplyURI("mongodb://mongodb-service:27017")

	// Conecta al cliente de MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Comprueba la conexión
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	return client, nil
}

func InsertLog(client *mongo.Client, data models.Vote) {
	collection := client.Database("so1_proyecto2").Collection("votes")
	_, err := collection.InsertOne(context.TODO(), data)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Log insertado en mongoDB!")
}

func InsertMongo(data models.Vote) {
	if clientMongo == nil {
		var err error
		clientMongo, err = ConnectMongo()
		if err != nil {
			log.Fatal(err)
		}
	}

	InsertLog(clientMongo, data)
}
