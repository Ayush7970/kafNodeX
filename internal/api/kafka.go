package api

import (
    "context"
    "github.com/segmentio/kafka-go"
    "github.com/sirupsen/logrus"
    "time"
    "io"
	"fmt"
)
var log = logrus.New()

func init() {
    log.Formatter = &logrus.JSONFormatter{}
    log.Level = logrus.InfoLevel
}




func KafkaWriter(brokers []string, topic string) *kafka.Writer {
    if topic == "" {
        log.Println("Error: KafkaWriter received empty topic") // Check for empty topic
        return nil
    }

    log.Printf("Attempting to dial Kafka broker at: %s", brokers[0])     // Connect to Kafka broker
    conn, err := kafka.Dial("tcp", brokers[0])
    if err != nil {
        log.WithFields(logrus.Fields{
            "broker": brokers[0],
            "error":  err.Error(),
        }).Error("Failed to connect to Kafka broker")    // Log error if broker connection fails
        return nil
    }
    defer conn.Close()
    log.Println("Kafka broker connection successful")

    log.Printf("Attempting to create topic: %s", topic)       // Try to create the topic
    err = conn.CreateTopics(kafka.TopicConfig{
        Topic:             topic,
        NumPartitions:     1,
        ReplicationFactor: 1,
    })
    if err != nil {
        log.WithFields(logrus.Fields{
            "topic": topic,
            "error": err.Error(),
        }).Warn("Could not create Kafka topic (it may already exist)")
    } else {       
        log.WithField("topic", topic).Info("Successfully created Kafka topic")         // Confirm successful topic creation
    }

    log.Println("Initializing Kafka writer")
    writer := kafka.NewWriter(kafka.WriterConfig{
        Brokers:  brokers,
        Topic:    topic,
        Balancer: &kafka.LeastBytes{},
    })

    if writer == nil {
        log.Println("Error: Kafka writer initialization failed, returning nil")
    } else {
        log.WithField("topic", topic).Info("Kafka writer initialized with topic")
    }

    return writer
}



// KafkaReader sets up a new Kafka consumer
func KafkaReader(brokers []string, topic string, groupID string) *kafka.Reader {
    reader := kafka.NewReader(kafka.ReaderConfig{
        Brokers:     brokers,
        Topic:       topic,
        GroupID:     groupID,
        MinBytes:    10e3,  // 10KB
        MaxBytes:    10e6,  // 10MB
        StartOffset: kafka.FirstOffset,
    })

	log.WithFields(logrus.Fields{
        "topic":   topic,
        "groupID": groupID,
    }).Info("Kafka reader initialized")
    return reader
}

// ProduceMessage sends messages to the Kafka topic with error handling
func ProduceMessage(w *kafka.Writer, key, message []byte) error {
    if w == nil {
        err := fmt.Errorf("Kafka writer is not initialized")
        log.WithField("error", err.Error()).Error("Failed to produce message")
        return err
    }

    attempt := 1
    maxRetries := 3                          // Set the maximum number of retry attempts
    for attempt <= maxRetries {
        err := w.WriteMessages(context.Background(),
            kafka.Message{
                Key:   key,
                Value: message,
            },
        )
        if err != nil {
            log.WithFields(logrus.Fields{
                "attempt": attempt,
                "key":     string(key),
                "error":   err.Error(),
            }).Error("Failed to write message to Kafka")							// Log error if message production fails

            if attempt == maxRetries {
                return err
            }
            attempt++
            time.Sleep(2 * time.Second) // Exponential backoff
            continue
        }

        log.WithFields(logrus.Fields{
            "key":   string(key),
            "topic": w.Topic,
        }).Info("Successfully produced message")               // Log success if message is produced
        return nil
    }

    return fmt.Errorf("failed to produce message after %d attempts", maxRetries)
}

// ConsumeMessages reads messages from the Kafka topic with error handling
func ConsumeMessages(r *kafka.Reader) error {
    if r == nil {
		err := fmt.Errorf("Kafka reader is not initialized")
        log.WithField("error", err.Error()).Error("Failed to consume messages")
        return err    
	}

    for {
        m, err := r.ReadMessage(context.Background())
        if err != nil {
            if err == io.EOF {
                log.Info("Reached end of stream; waiting for new messages")    	// End of stream reached, wait for new messages
                time.Sleep(1 * time.Second)
                continue
            }
            log.WithField("error", err.Error()).Error("Failed to read message from Kafka")
            return err
        }
        
		log.WithFields(logrus.Fields{
            "offset": m.Offset,
            "key":    string(m.Key),
            "value":  string(m.Value),
        }).Info("Consumed message from Kafka")								

    }
}			// Log successfully consumed message
