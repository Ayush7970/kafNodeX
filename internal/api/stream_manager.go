package api

import (
    // "context"
    //"log"
    "sync"
    "github.com/segmentio/kafka-go"
)

type StreamManager struct {
    producers map[string]*kafka.Writer
    consumers map[string]*kafka.Reader
    mu        sync.Mutex
}

func NewStreamManager() *StreamManager {
    return &StreamManager{
        producers: make(map[string]*kafka.Writer),
        consumers: make(map[string]*kafka.Reader),
    }
}

func (sm *StreamManager) CreateProducer(brokers []string, streamID string) *kafka.Writer {
    sm.mu.Lock()
    defer sm.mu.Unlock()

	
		// mostly error handling and extra logging statements
    if streamID == "" {
        log.Println("Error: Received empty streamID for producer creation")
        return nil
    }

    if producer, exists := sm.producers[streamID]; exists {
        log.Printf("Returning existing producer for streamID: %s", streamID)
        return producer
    }

	

    log.Printf("Creating new producer for streamID: %s", streamID)
    producer := KafkaWriter(brokers, streamID)
    if producer == nil {
        log.Printf("Failed to create Kafka producer for streamID: %s", streamID)
        return nil
    }
    sm.producers[streamID] = producer
    return producer
}


// CreateConsumer initializes a new Kafka consumer for a given stream
func (sm *StreamManager) CreateConsumer(brokers []string, streamID, groupID string) *kafka.Reader {
    sm.mu.Lock()
    defer sm.mu.Unlock()

    if consumer, exists := sm.consumers[streamID]; exists {
        return consumer
    }

    consumer := KafkaReader(brokers, streamID, groupID)
    sm.consumers[streamID] = consumer
    return consumer
}

// CloseStream gracefully shuts down the producer and consumer for a stream
func (sm *StreamManager) CloseStream(streamID string) {
    sm.mu.Lock()
    defer sm.mu.Unlock()

    if producer, exists := sm.producers[streamID]; exists {
        producer.Close()
        delete(sm.producers, streamID)
    }
    if consumer, exists := sm.consumers[streamID]; exists {
        consumer.Close()
        delete(sm.consumers, streamID)
    }
}
