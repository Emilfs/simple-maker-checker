package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	"github.com/yuin/goldmark"
)

type Item struct {
	ID          string `json:"id"`
	Message     string `json:"message"`
	RecipientID string `json:"recipient_id"`
	Status      string `json:"status"` // "pending", "approved", "rejected"
}

var (
	items = make(map[string]*Item)
	mu    sync.Mutex
	port  = "8080"

	APPROVED = "approved"
	REJECTED = "rejected"
	PENDING  = "pending"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/items", createItem).Methods("POST")
	r.HandleFunc("/items/{id}/approve", approveItem).Methods("PUT")
	r.HandleFunc("/items/{id}/reject", rejectItem).Methods("PUT")
	r.HandleFunc("/items", listItems).Methods("GET")
	r.HandleFunc("/", rootHandler).Methods("GET")

	log.Printf("Serving at Port %s\n", port)
	http.ListenAndServe(":"+port, r)
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fileContent, err := ioutil.ReadFile("README.md")
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	var buf bytes.Buffer
	if err := goldmark.Convert(fileContent, &buf); err != nil {
		http.Error(w, "Failed to convert markdown", http.StatusInternalServerError)
		return
	}
	w.Write(buf.Bytes())
}

func createItem(w http.ResponseWriter, r *http.Request) {
	var item Item
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	mu.Lock()
	defer mu.Unlock()
	if _, exists := items[item.ID]; exists {
		http.Error(w, "Item already exists", http.StatusConflict)
		return
	}
	item.Status = PENDING
	items[item.ID] = &item
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(item)
}
func approveItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := mux.Vars(r)["id"]
	mu.Lock()
	// assumption: rejected messages can still be approved and sent
	item, exists := items[id]
	mu.Unlock()
	if !exists {
		http.Error(w, "Item not found", http.StatusNotFound)
		return
	}
	item.Status = APPROVED
	go sendMessage(item) // Send message in a non-blocking way
	json.NewEncoder(w).Encode(item)
}
func rejectItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := mux.Vars(r)["id"]
	mu.Lock()
	defer mu.Unlock()
	item, exists := items[id]
	if item.Status == APPROVED {
		http.Error(w, "Message has been approved before", http.StatusBadRequest)
		return
	}
	if !exists {
		http.Error(w, "Item not found", http.StatusNotFound)
		return
	}
	item.Status = REJECTED
	json.NewEncoder(w).Encode(item)
}
func listItems(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	mu.Lock()
	defer mu.Unlock()
	var itemList []Item
	for _, item := range items {
		itemList = append(itemList, *item)
	}
	json.NewEncoder(w).Encode(itemList)
}
func sendMessage(item *Item) {
	// Simulate sending the message.
	log.Printf("Sending message to Recipient ID %s : %s\n", item.RecipientID, item.Message)
	mu.Lock()
	defer mu.Unlock()
}
