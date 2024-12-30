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
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type AggregatedStoreData struct {
	StoreId          string  `json:"store_id"`
	OrderCount       int     `json:"order_count"`
	TotalOrderAmount float32 `json:"total_order_amount"`
	Lat              float32 `json:"lat"`
	Lng              float32 `json:"lng"`
	Received         int64
}

type AggregatedAllStoreData struct {
	AllStoreTotalOrderCount  int     `json:"all_store_total_order_count"`
	AllStoreTotalOrderAmount float32 `json:"all_store_total_order_amount"`
}

var storeData = make(map[string]AggregatedStoreData)
var storeMux sync.RWMutex
var broadcast = make(chan []byte)
var allStoreData = AggregatedAllStoreData{0, 0}

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

			var message AggregatedStoreData
			if err := json.Unmarshal(msg.Value, &message); err != nil {
				log.Printf("Failed to unmarshal message: %v", err)
				continue
			}

			message.Received = time.Now().UnixMilli()
			storeID := message.StoreId
			storeMux.Lock()
			storeData[storeID] = message
			storeMux.Unlock()

			select {
			case broadcast <- msg.Value:
			default:
			}
		}
	}
}

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

func snapshotHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	storeMux.RLock()
	defer storeMux.RUnlock()

	if err := json.NewEncoder(w).Encode(storeData); err != nil {
		http.Error(w, "Failed to encode snapshot", http.StatusInternalServerError)
	}
}

func allStoreHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	storeMux.RLock()

	allStoreData.AllStoreTotalOrderCount = 0
	allStoreData.AllStoreTotalOrderAmount = 0.0

	for _, store := range storeData {
		allStoreData.AllStoreTotalOrderCount += store.OrderCount
		allStoreData.AllStoreTotalOrderAmount += store.TotalOrderAmount
	}

	storeMux.RUnlock()

	if err := json.NewEncoder(w).Encode(allStoreData); err != nil {
		http.Error(w, "Failed to encode snapshot", http.StatusInternalServerError)
	}

}

func reactHandler(w http.ResponseWriter, r *http.Request) {
	currentDir, _ := os.Getwd()
	buildPath := filepath.Join(currentDir, "frontend", "build")
	fs := http.FileServer(http.Dir(buildPath))
	wrappedHandler := http.StripPrefix("/", fs)
	wrappedHandler.ServeHTTP(w, r)
}

func staleDataCheck() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			storeMux.Lock()
			for storeId, data := range storeData {
				if time.Now().UnixMilli()-data.Received > 10000 {
					data.OrderCount = 0
					data.TotalOrderAmount = 0
					msg, _ := json.Marshal(data)
					broadcast <- msg
					delete(storeData, storeId)
				}
			}
			storeMux.Unlock()
		}
	}
}

func main() {
	broker := "kafka:29092"
	topic := "aggregated_store_orders"
	groupID := "store-aggregation-group"

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go runKafkaConsumer(ctx, broker, topic, groupID)

	go staleDataCheck()

	http.HandleFunc("/snapshot", snapshotHandler)
	http.HandleFunc("/events", sseHandler)
	http.HandleFunc("/aggregates", allStoreHandler)
	http.HandleFunc("/", reactHandler)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("HTTP server error: %v", err)
	}
}
