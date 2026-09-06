package kafka

import (
	"context"
	"encoding/json"
	"log"
	models "server/models"
	"time"

	"github.com/segmentio/kafka-go"
)

func SendData(data models.Vote) {
	// to produce messages
	topic := "so1-proyecto2"
	partition := 0

	conn, err := kafka.DialLeader(context.Background(), "tcp", "my-cluster-kafka-bootstrap:9092", topic, partition)
	if err != nil {
		log.Fatal("failed to dial leader:", err)
	}

	conn.SetWriteDeadline(time.Now().Add(10 * time.Second))

	// Convert the Album to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Fatal("failed to marshal data:", err)
	}

	_, err = conn.WriteMessages(
		kafka.Message{Value: jsonData},
	)
	if err != nil {
		log.Fatal("failed to write messages:", err)
	}

	if err := conn.Close(); err != nil {
		log.Fatal("failed to close writer:", err)
	}

	log.Println("Message sent successfully")
}
