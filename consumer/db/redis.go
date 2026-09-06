package db

import (
	"log"

	models "consumer/models"

	"github.com/go-redis/redis"
)

var clientRedis *redis.Client

func ConnectRedis() *redis.Client {
	if clientRedis == nil {
		clientRedis = redis.NewClient(&redis.Options{
			Addr:     "redis-service:6379", // Actualiza con tu dirección y puerto
			Password: "",                   // No password set
			DB:       0,                    // use default DB
		})
	}

	_, err := clientRedis.Ping().Result()
	if err != nil {
		log.Fatal(err)
	}

	return clientRedis
}

func InsertVote(client *redis.Client, data models.Vote) {
	// Incrementa el conteo del álbum
	err := client.HIncrBy("votes", data.Data.Album, 1).Err()
	if err != nil {
		log.Println(err)
	} else {
		log.Println("Voto registrado para el álbum:", data.Data.Album)
	}
}

func InsertRedis(data models.Vote) {
	// Conecta al cliente de Redis una vez al inicio del programa
	client := ConnectRedis()

	InsertVote(client, data)
}
