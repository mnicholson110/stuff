package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// Structs for incoming Kafka messages
type Message struct {
	OrderId    uint   `json:"order_id"`
	Data       string `json:"data"` // Handle as raw JSON string
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
	SourceTsMs uint64 `json:"__source_ts_ms"`
}

type MessageData struct {
	Order OrderData `json:"order"`
	Store StoreData `json:"store"`
}

type OrderData struct {
	OrderAmount   float32 `json:"order_amount"`
	OrderStatus   string  `json:"order_status"`
	OrderStatusId int     `json:"order_status_id"`
}

type StoreData struct {
	StoreId   int    `json:"store_id"`
	StoreLoc  Loc    `json:"store_loc"`
	StoreAddr string `json:"store_addr"`
}

type Loc struct {
	Lat  float32 `json:"store_lat"`
	Long float32 `json:"store_long"`
}

// Struct for aggregated store data
type AggregatedStoreData struct {
	StoreId          int     `json:"store_id"`
	OrderCount       int     `json:"order_count"`
	TotalOrderAmount float32 `json:"total_order_amount"`
	StoreAddr        string  `json:"store_addr"`
	Lat              float32 `json:"lat"`
	Long             float32 `json:"long"`
}

// State management
var storeData = make(map[int]*AggregatedStoreData)
var storeMux sync.RWMutex
var broadcast = make(chan []byte)

// Kafka consumer
func runKafkaConsumer(ctx context.Context, broker string, topic string, groupID string) {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": broker,
		"group.id":          groupID,
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		log.Fatalf("Failed to create Kafka consumer: %v", err)
	}
	defer consumer.Close()

	err = consumer.Subscribe(topic, nil)
	if err != nil {
		log.Fatalf("Failed to subscribe to topic: %v", err)
	}
	log.Printf("Subscribed to topic: %s", topic)

	for {
		select {
		case <-ctx.Done():
			return
		default:
			msg, err := consumer.ReadMessage(-1)
			if err != nil {
				log.Printf("Consumer error: %v", err)
				continue
			}

			var message Message
			// Unmarshal the raw message first
			if err := json.Unmarshal(msg.Value, &message); err != nil {
				log.Printf("Failed to unmarshal message: %v", err)
				continue
			}

			// Unmarshal the stringified "data" field
			var messageData MessageData
			if err := json.Unmarshal([]byte(message.Data), &messageData); err != nil {
				log.Printf("Failed to unmarshal data field: %v", err)
				continue
			}

			// Filter for "Delivered" orders
			if messageData.Order.OrderStatus != "Delivered" {
				continue
			}

			// Aggregate state
			storeID := messageData.Store.StoreId
			storeMux.Lock()
			if _, exists := storeData[storeID]; !exists {
				storeData[storeID] = &AggregatedStoreData{
					StoreId:          storeID,
					StoreAddr:        messageData.Store.StoreAddr,
					Lat:              messageData.Store.StoreLoc.Lat,
					Long:             messageData.Store.StoreLoc.Long,
					OrderCount:       0,
					TotalOrderAmount: 0,
				}
			}
			storeData[storeID].OrderCount++
			storeData[storeID].TotalOrderAmount += messageData.Order.OrderAmount
			storeMux.Unlock()

			// Broadcast update over SSE
			updatedData, err := json.Marshal(storeData[storeID])
			if err != nil {
				log.Printf("Failed to marshal aggregated data: %v", err)
				continue
			}
			select {
			case broadcast <- updatedData:
			default: // No active connections; drop the message
			}
		}
	}
}

// SSE Handler
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
	go func() {
		<-r.Context().Done()
		close(clientChan)
	}()

	for {
		select {
		case msg := <-broadcast:
			_, _ = fmt.Fprintf(w, "data: %s\n\n", msg)
			flusher.Flush()
		case <-r.Context().Done():
			return
		}
	}
}

// Snapshot Handler
func snapshotHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	storeMux.RLock()
	defer storeMux.RUnlock()

	if err := json.NewEncoder(w).Encode(storeData); err != nil {
		http.Error(w, "Failed to encode snapshot", http.StatusInternalServerError)
	}
}

// Serve React App
func reactHandler(w http.ResponseWriter, r *http.Request) {
	// Static file server for React build folder
	currentDir, _ := os.Getwd()
	buildPath := filepath.Join(currentDir, "frontend", "build")
	fs := http.FileServer(http.Dir(buildPath))
	wrappedHandler := http.StripPrefix("/", fs)
	wrappedHandler.ServeHTTP(w, r)
}

func main() {
	// Kafka configuration
	broker := "kafka:29092"
	topic := "order_db.order_schema.order"
	groupID := "store-aggregation-group"

	// Context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Run Kafka consumer
	go runKafkaConsumer(ctx, broker, topic, groupID)

	// HTTP Handlers
	http.HandleFunc("/snapshot", snapshotHandler)
	http.HandleFunc("/events", sseHandler)
	http.HandleFunc("/", reactHandler)

	// Start HTTP server
	log.Println("Starting HTTP server on :8081")
	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatalf("HTTP server error: %v", err)
	}
}
