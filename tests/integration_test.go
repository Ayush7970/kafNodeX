// tests/integration_test.go
package tests

import (
    "bytes"
    "encoding/json"
    "my-golang-api/internal/api"
    "net/http"
    "net/http/httptest"
    "testing"
	"github.com/gorilla/websocket"
	"github.com/gorilla/mux"
	"time"
)

// TestStartStreamEndpoint tests the /stream/start endpoint for initiating a new stream.
func TestStartStreamEndpoint(t *testing.T) {
	api.RegisterMetrics()
    req := httptest.NewRequest("POST", "/stream/start", nil)
    w := httptest.NewRecorder()
    handler := http.HandlerFunc(api.StartStream)

    handler.ServeHTTP(w, req)

    if w.Code != http.StatusOK {
        t.Errorf("Expected status OK, got %v", w.Code)
    }

    var response map[string]string
    err := json.Unmarshal(w.Body.Bytes(), &response)
    if err != nil {
        t.Fatalf("Failed to parse JSON response: %v", err)
    }

    if _, ok := response["stream_id"]; !ok {
        t.Errorf("Expected stream_id in response, got %v", response)
    }
}




func TestSendDataEndpoint(t *testing.T) {
	api.RegisterMetrics()

    // Step 1: Start a new stream
    startReq := httptest.NewRequest("POST", "/stream/start", nil)
    startW := httptest.NewRecorder()
    startHandler := http.HandlerFunc(api.StartStream)

    startHandler.ServeHTTP(startW, startReq)

    var startResponse map[string]string
    err := json.Unmarshal(startW.Body.Bytes(), &startResponse)
    if err != nil {
        t.Fatalf("Failed to parse JSON response: %v", err)
    }

    streamID, exists := startResponse["stream_id"]
    if !exists || streamID == "" {
        t.Fatalf("Failed to retrieve valid stream_id from start stream response: %v", startResponse)
    }

    t.Logf("Successfully retrieved streamID: %s", streamID)

    // Sending data to the created stream
    sendPayload := map[string]string{"data": "test_data"}
    jsonPayload, _ := json.Marshal(sendPayload)
    sendReq := httptest.NewRequest("POST", "/stream/"+streamID+"/send", bytes.NewBuffer(jsonPayload))
    sendReq.Header.Set("Content-Type", "application/json")
    sendW := httptest.NewRecorder()

    // Call SendData with streamID directly
    api.SendData(sendW, sendReq, streamID)

    // Validating the response
    if sendW.Code != http.StatusOK {
        t.Errorf("Expected status OK for send data, got %v", sendW.Code)
    }

    if !bytes.Contains(sendW.Body.Bytes(), []byte("Data sent and processed")) {
        t.Errorf("Expected response to contain 'Data sent and processed', got %s", sendW.Body.String())
    }
}




// TestGetResultsEndpoint tests the /stream/{stream_id}/results endpoint for retrieving results via WebSocket.
func TestGetResultsEndpoint(t *testing.T) {
    // Starting the server to test WebSocket connections
    router := mux.NewRouter()
    router.HandleFunc("/stream/start", api.StartStream).Methods("POST")
    router.HandleFunc("/stream/{stream_id}/results", api.GetResults).Methods("GET")

    server := httptest.NewServer(router)
    defer server.Close()

    // Starting a new stream and retrieve the stream_id
    startResp, err := http.Post(server.URL+"/stream/start", "application/json", nil)
    if err != nil {
        t.Fatalf("Failed to start stream: %v", err)
    }
    defer startResp.Body.Close()

    var startResponse map[string]string
    if err := json.NewDecoder(startResp.Body).Decode(&startResponse); err != nil {
        t.Fatalf("Failed to parse start stream response: %v", err)
    }

    streamID, exists := startResponse["stream_id"]
    if !exists || streamID == "" {
        t.Fatalf("Failed to retrieve valid stream_id from start stream response")
    }

    t.Logf("Successfully retrieved streamID: %s", streamID)

    // Establishing a WebSocket connection to the results endpoint
    wsURL := "ws" + server.URL[len("http"):] + "/stream/" + streamID + "/results"
    dialer := websocket.Dialer{
        HandshakeTimeout: 5 * time.Second,
    }

    conn, _, err := dialer.Dial(wsURL, nil)
    if err != nil {
        t.Fatalf("Failed to establish WebSocket connection: %v", err)
    }
    defer conn.Close()

    // Verifying WebSocket message reception (Optional, depending on server implementation)
    conn.SetReadDeadline(time.Now().Add(5 * time.Second))
    _, message, err := conn.ReadMessage()
    if err != nil {
        t.Fatalf("Failed to read WebSocket message: %v", err)
    }

    t.Logf("Received WebSocket message: %s", message)
}


