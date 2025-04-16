package api

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "sync"
    "io"
	

    "github.com/google/uuid"
    "github.com/gorilla/mux"
    "github.com/gorilla/websocket"
    "github.com/segmentio/kafka-go"
    "time"

)



var (
    streamManager   = NewStreamManager()
    wsConnections   = make(map[string]*websocket.Conn) // Store WebSocket connections per stream
    wsMutex         = sync.Mutex{}
    upgrader        = websocket.Upgrader{
        ReadBufferSize:  1024,
        WriteBufferSize: 1024,
        CheckOrigin: func(r *http.Request) bool { return true }, // Allow connections from any origin
    }
)



// Handler for starting a new data stream
func StartStream(w http.ResponseWriter, r *http.Request) {
	streamID := uuid.New().String()

    producer := streamManager.CreateProducer([]string{"localhost:9092"}, streamID)
	if producer == nil {
        http.Error(w, "Failed to initialize Kafka producer", http.StatusInternalServerError)
        return
    }

	

	response := map[string]string{"message": "New stream started", "stream_id": streamID}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}


func SendData(w http.ResponseWriter, r *http.Request, streamID string) {
    start := time.Now() // Start timing the request

    // Log to confirm the `streamID` value
    if streamID == "" {
        http.Error(w, "Invalid stream ID", http.StatusBadRequest)
        log.Println("Error: Received empty streamID in SendData")
        httpRequestsTotal.WithLabelValues("400", "POST").Inc() // Log failed request
        return
    }
    log.Printf("Processing SendData for streamID: %s", streamID)

    

    //  Creating a producer for the stream
    producer := streamManager.CreateProducer([]string{"localhost:9092"}, streamID)
    if producer == nil {
        http.Error(w, "Failed to initialize Kafka producer", http.StatusInternalServerError)
        log.Printf("Error: Failed to create Kafka producer for streamID %s", streamID)
        httpRequestsTotal.WithLabelValues("500", "POST").Inc()
        return
    }
	log.Println("Kafka producer initialized successfully")

	// Decode JSON request body
    var requestBody map[string]string
    err := json.NewDecoder(r.Body).Decode(&requestBody)
    if err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        httpRequestsTotal.WithLabelValues("400", "POST").Inc()
        return
    }
    data, exists := requestBody["data"]
    if !exists {
        http.Error(w, "Missing 'data' field in request body", http.StatusBadRequest)
        httpRequestsTotal.WithLabelValues("400", "POST").Inc()
        return
    }

    // Sending the raw data to Kafka
    err = producer.WriteMessages(context.Background(),
        kafka.Message{
            Key:   []byte("key"),
            Value: []byte(data),
        },
    )
    if err != nil {
        http.Error(w, "Failed to send data to Kafka: "+err.Error(), http.StatusInternalServerError)
        httpRequestsTotal.WithLabelValues("500", "POST").Inc()
        return
    }

	log.Println("Successfully wrote message to Kafka")

    // Increment Kafka messages produced counter
    kafkaMessagesProduced.Inc()

    // Processing the data in real-time
    processedData := ProcessData(data)

    // Pushing the processed data back to the client via WebSocket if a connection exists
    wsMutex.Lock()
    conn, exists := wsConnections[streamID]
    wsMutex.Unlock()

    if exists && conn != nil {
        err = conn.WriteMessage(websocket.TextMessage, []byte(processedData))
        if err != nil {
            log.Printf("Failed to send data to WebSocket for stream %s: %v", streamID, err)
        }
    } else {
        log.Printf("No WebSocket connection found for stream %s", streamID)
    }

    // Sending a response indicating the data was processed and sent
    response := fmt.Sprintf("Data sent and processed for stream %s", streamID)
    w.Header().Set("Content-Type", "text/plain; charset=utf-8")
    w.WriteHeader(http.StatusOK)
    w.Write([]byte(response))

    // Record the total time taken for the request
    duration := time.Since(start).Seconds()
    httpRequestDuration.WithLabelValues("POST").Observe(duration)
    httpRequestsTotal.WithLabelValues("200", "POST").Inc()
}








// Handler for establishing WebSocket connection for real-time results
func GetResults(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    streamID := vars["stream_id"]

    // Upgrade the HTTP connection to a WebSocket connection
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Printf("Failed to upgrade to WebSocket: %v", err)
        http.Error(w, "Failed to open WebSocket connection", http.StatusInternalServerError)
        return
    }

    // Store the WebSocket connection for this stream ID
    wsMutex.Lock()
    wsConnections[streamID] = conn
    wsMutex.Unlock()

    // Ensure connection cleanup
    defer func() {
        wsMutex.Lock()
        delete(wsConnections, streamID)
        wsMutex.Unlock()
        conn.Close()
        log.Printf("WebSocket connection for stream %s closed", streamID)
    }()

    // Create or get a Kafka consumer for the stream ID topic
    consumer := streamManager.CreateConsumer([]string{"localhost:9092"}, streamID, "group-"+streamID)
    defer consumer.Close() // Ensure the consumer is closed when done

    // Context to handle cancellation
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel() // Ensure context cancellation

    // Inform the WebSocket client that consumption has started
    conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Started consuming messages for stream %s", streamID)))

    // Launch a goroutine to read messages from Kafka and send to WebSocket
    go func() {
        for {
            // Read message from Kafka
            m, err := consumer.ReadMessage(ctx)
            if err != nil {
                if err == io.EOF {
                    log.Println("No new messages available.")
                    conn.WriteMessage(websocket.TextMessage, []byte("No new messages available."))
                    return
                }
                log.Printf("Error reading messages for stream %s: %v", streamID, err)
                conn.WriteMessage(websocket.TextMessage, []byte("Error reading messages: "+err.Error()))
                return
            }

            // Process the message
            processedMessage := ProcessData(string(m.Value))

            // Send the processed message to WebSocket
            err = conn.WriteMessage(websocket.TextMessage, []byte(processedMessage))
            if err != nil {
                log.Printf("Failed to send message to WebSocket for stream %s: %v", streamID, err)
                cancel() // Cancel Kafka reading on WebSocket error
                return
            }

            log.Printf("Sent message to WebSocket for stream %s: %s", streamID, processedMessage)
        }
    }()

    // Keep the WebSocket connection open until the client disconnects
    for {
        _, _, err := conn.ReadMessage()
        if err != nil {
            log.Printf("WebSocket connection for stream %s closed by client: %v", streamID, err)
            break
        }
    }
}
