package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
)

const (
	maxSize = 10 * (1 << 20) // 10MB
)

var (
	db     = make(map[string][]byte)
	dbLock sync.RWMutex
)

func handleSet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	defer r.Body.Close()
	rdr := io.LimitReader(r.Body, maxSize)
	data, err := io.ReadAll(rdr)
	if err != nil {
		log.Printf("Reader Error: %s", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	dbLock.Lock()         // acquiring read-write lock for db operation
	defer dbLock.Unlock() // ensuring read-write lock is released properly
	db[key] = data        // Updating in-memory database with data (request payload) for the given key

	resp := map[string]interface{}{
		"key":  key,
		"size": len(data),
	}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Error sending: %s", err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

}

func handleGet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	dbLock.RLock()         // acquiring read lock for db operation
	defer dbLock.RUnlock() // ensuring read lock is released properly

	data, ok := db[key]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func handleList(w http.ResponseWriter, r *http.Request) {
	dbLock.RLock()         // acquiring read lock for db operation
	defer dbLock.RUnlock() // ensuring read lock is released properly

	keys := make([]string, 0, len(db))
	for key := range db {
		keys = append(keys, key)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(keys); err != nil {
		log.Printf("error sending: %s", err)
	} else {
		w.WriteHeader(http.StatusOK)
	}

}

func handleDelete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	dbLock.RLock()         // acquiring read lock for db operation
	defer dbLock.RUnlock() // ensuring read lock is released properly
	data := db[key]
	delete(db, key)
	resp := map[string]interface{}{
		"data": data,
		"size": len(data),
	}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Error sending: %s", err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func main() {
	r := mux.NewRouter()
	// handle routing
	r.HandleFunc("/kv/{key}", handleSet).Methods("POST")
	r.HandleFunc("/kv/{key}", handleGet).Methods("GET")
	r.HandleFunc("/kv/{key}", handleDelete).Methods("DELETE")
	r.HandleFunc("/kv", handleList).Methods("GET")
	http.Handle("/", r)

	addr := ":8080"
	log.Printf("server ready on %s", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}
