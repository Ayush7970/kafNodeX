package main

import (
    "log"
    "my-golang-api/internal/api"
    "net/http"
    "github.com/gorilla/mux"
)


func sendDataWrapper(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    streamID, exists := vars["stream_id"]
    if !exists || streamID == "" {
        http.Error(w, "Stream ID is required", http.StatusBadRequest)
        return
    }
    
    // Pass `streamID` to `SendData`
    api.SendData(w, r, streamID)
}



func main() {
	api.RegisterMetrics()
    router := mux.NewRouter()

    router.Use(apiKeyAuthMiddleware)

	

    // Use the global rate limiter middleware
    router.Use(api.RateLimiterMiddleware())

	
  

    // Update handlers to use api package
    router.HandleFunc("/stream/start", api.StartStream).Methods("POST")
    router.HandleFunc("/stream/{stream_id}/send", sendDataWrapper).Methods("POST")
    router.HandleFunc("/stream/{stream_id}/results", api.GetResults).Methods("GET")

    router.Handle("/metrics", api.MetricsHandler())


    log.Println("Server is running on http://localhost:8080")
    if err := http.ListenAndServe(":8080", router); err != nil {
        log.Fatalf("Failed to start server: %s", err)
    }
}
