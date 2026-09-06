package main

import (
	"context"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Data struct {
	Name  string
	Album string
	Year  string
	Rank  string
}

type Vote struct {
	Data         Data
	CreationDate string
}

var collection *mongo.Collection

func main() {
	// Conexión a MongoDB
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://35.193.107.74:27017"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	cancel()
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.Background())

	// Acceder a la colección "votes" de la base de datos "so1_proyecto2"
	collection = client.Database("so1_proyecto2").Collection("votes")

	// Crear una nueva aplicación Fiber
	app := fiber.New()

	// Agregar el middleware CORS
	app.Use(cors.New())

	// Definir la ruta /logs
	app.Get("/logs", getLogs)

	// Iniciar el servidor
	log.Fatal(app.Listen(":5000"))
}

func getLogs(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Obtener los últimos 10 documentos de la colección
	opts := options.Find()
	opts.SetLimit(20)
	opts.SetSort(bson.D{{Key: "$natural", Value: -1}})
	cursor, err := collection.Find(ctx, bson.D{}, opts)
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(ctx)

	var results []*Vote
	if err = cursor.All(ctx, &results); err != nil {
		log.Fatal(err)
	}

	// Devolver los resultados como JSON
	return c.JSON(results)
}
