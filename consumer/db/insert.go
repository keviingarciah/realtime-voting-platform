package db

import (
	"consumer/models"
	"encoding/json"
	"log"
	"time"
)

func InsertData(event []byte) {
	// Convert the Album to JSON
	var data models.Data
	err := json.Unmarshal(event, &data)
	if err != nil {
		log.Fatal("failed to unmarshal data:", err)
	}

	log := models.Vote{
		Data:         data,
		CreationDate: time.Now().String(),
	}

	InsertMongo(log)
	go InsertRedis(log)
}
