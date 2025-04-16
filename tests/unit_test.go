// tests/unit_test.go
package tests

import (
    "my-golang-api/internal/api"
    "testing"
    "strings"
	"context"
	"github.com/segmentio/kafka-go"

)
// TestProcessData checks if the ProcessData function correctly processes input data

func TestProcessData(t *testing.T) {
    input := "test_data"
    result := api.ProcessData(input)

	    // Check that the result contains both "Processed" and the input data

    if !strings.Contains(result, "Processed") || !strings.Contains(result, input) {
        t.Errorf("Expected processed result to contain 'Processed' and input data, got %s", result)
    }
}


// TestKafkaWriter tests the KafkaWriter function by initializing a writer and sending a test message

func TestKafkaWriter(t *testing.T) {
    brokers := []string{"localhost:9092"}    // Set the Kafka broker address
    topic := "test_topic"

    writer := api.KafkaWriter(brokers, topic)
    if writer == nil {
        t.Fatal("Expected a Kafka writer, got nil")    // Fail if the writer is nil
    }
    defer writer.Close()

    err := writer.WriteMessages(context.Background(), kafka.Message{Key: []byte("test_key"), Value: []byte("test_value")})
    if err != nil {
        t.Errorf("Failed to write message to Kafka: %v", err)
    }
}



