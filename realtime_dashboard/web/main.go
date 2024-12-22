package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

// Global state for store data
var (
	storeData = make(map[int]string)
	storeMux  sync.RWMutex
)

// SSE infrastructure
type sseClient chan []byte

var (
	register      = make(chan sseClient)
	unregister    = make(chan sseClient)
	broadcast     = make(chan []byte)
	activeClients = make(map[sseClient]bool)
)

// Run the SSE broker
func runSSEBroker() {
	for {
		select {
		case client := <-register:
			activeClients[client] = true
		case client := <-unregister:
			delete(activeClients, client)
			close(client)
		case msg := <-broadcast:
			// Send to all connected clients
			for client := range activeClients {
				select {
				case client <- msg:
				default:
					// If blocked, remove it
					delete(activeClients, client)
					close(client)
				}
			}
		}
	}
}

// SSE handler
func sseHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	clientChan := make(chan []byte)
	register <- clientChan
	defer func() { unregister <- clientChan }()

	notify := r.Context().Done()

	for {
		select {
		case <-notify:
			return // client closed connection
		case msg := <-clientChan:
			if len(msg) == 0 {
				log.Println("Skipping empty SSE message.")
				continue
			}
			fmt.Fprintf(w, "data: %s\n\n", msg)
			flusher.Flush()
		}
	}
}

// Snapshot endpoint
func snapshotHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	storeMux.RLock()
	snap := make(map[int]string, len(storeData))
	for k, v := range storeData {
		snap[k] = v
	}
	storeMux.RUnlock()

	json.NewEncoder(w).Encode(snap)
}

// Kafka consumer
func runKafkaConsumer(ready chan<- bool) {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "kafka:29092",
		"group.id":          "go-sse-group",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		log.Fatalf("Failed to create consumer: %v", err)
	}
	defer consumer.Close()

	if err := consumer.Subscribe("aggregated_store_orders", nil); err != nil {
		log.Fatalf("Failed to subscribe to topic: %v", err)
	}
	log.Println("[KafkaConsumer] Subscribed to aggregated_store_orders")

	for {
		msg, err := consumer.ReadMessage(-1)
		if err != nil {
			log.Printf("Consumer error: %v (%v)\n", err, msg)

			// Retry logic for errors from ReadMessage
			for i := 1; i <= 10; i++ {
				log.Printf("Retrying to read message (attempt %d/%d)...", i, 10)
				time.Sleep(10 * time.Second)

				msg, err = consumer.ReadMessage(-1)
				if err == nil {
					log.Println("Successfully connected to broker.")
					// Signal that Kafka is ready
					ready <- true
					break
				}

				log.Printf("Retry %d failed: %v", i, err)
				if i == 10 {
					log.Fatalf("Exceeded maximum retries (10) while reading messages from topic aggregated_store_orders")
				}
			}
		}

		// Validate and process the message
		if len(msg.Value) == 0 {
			continue
		}

		var partial struct {
			StoreID int `json:"storeId"`
		}
		if err := json.Unmarshal(msg.Value, &partial); err != nil {
			log.Printf("JSON parse error: %v", err)
			continue
		}

		storeID := partial.StoreID
		rawJSON := string(msg.Value)

		// Update storeData
		storeMux.Lock()
		storeData[storeID] = rawJSON
		storeMux.Unlock()

		// Broadcast the update to SSE clients
		broadcast <- []byte(rawJSON)
	}
}

// Serve React build
func ServeReactHandler(buildDir string) http.Handler {
	fs := http.FileServer(http.Dir(buildDir))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := filepath.Join(buildDir, r.URL.Path)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			http.ServeFile(w, r, filepath.Join(buildDir, "index.html"))
			return
		}
		fs.ServeHTTP(w, r)
	})
}

// Main function
func main() {
	// Start the SSE broker
	go runSSEBroker()

	// Signal to indicate Kafka readiness
	kafkaReady := make(chan bool)

	// Start Kafka consumer
	go runKafkaConsumer(kafkaReady)

	// Wait for Kafka consumer to signal readiness
	<-kafkaReady
	log.Println("[Main] Kafka consumer is ready. Starting server.")

	// API endpoints
	http.HandleFunc("/snapshot", snapshotHandler)
	http.HandleFunc("/events", sseHandler)

	// Serve React build
	reactBuildDir := "./frontend/build"
	http.Handle("/", ServeReactHandler(reactBuildDir))

	// Start HTTP server
	srv := &http.Server{
		Addr: ":8081",
	}
	log.Println("Serving on :8081...")
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
