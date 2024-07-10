package v1

import (
	"encoding/json"
	"linkedlist/linkedlist"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

type SafeLinkedList struct {
	list  *linkedlist.LinkedList
	mutex sync.Mutex
}

func NewSafeLinkedList() *SafeLinkedList {
	return &SafeLinkedList{list: linkedlist.NewLinkedList()}
}

func (s *SafeLinkedList) Find(n int) (index uint, found bool) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.list.Find(n)
}

func (s *SafeLinkedList) Get(index uint) (int, bool) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.list.Get(index)
}

func (s *SafeLinkedList) Insert(index uint, val int) bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.list.Insert(index, val)
}

func (s *SafeLinkedList) Remove(index uint) bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.list.Remove(index)
}

func handleInsert(w http.ResponseWriter, r *http.Request, list *SafeLinkedList) {
	var req struct {
		Index uint `json:"index"`
		Value int  `json:"value"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	success := list.Insert(req.Index, req.Value)
	if !success {
		http.Error(w, "Index out of range", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Insert successful"})
}

func handleGet(w http.ResponseWriter, r *http.Request, list *SafeLinkedList) {
	indexStr := strings.TrimPrefix(r.URL.Path, "/get/")
	index, err := strconv.Atoi(indexStr)
	if err != nil {
		http.Error(w, "Invalid index", http.StatusBadRequest)
		return
	}

	value, found := list.Get(uint(index))
	if !found {
		http.Error(w, "Index out of range", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{"value": value})
}

func handleRemove(w http.ResponseWriter, r *http.Request, list *SafeLinkedList) {
	indexStr := strings.TrimPrefix(r.URL.Path, "/remove/")
	index, err := strconv.Atoi(indexStr)
	if err != nil {
		http.Error(w, "Invalid index", http.StatusBadRequest)
		return
	}

	success := list.Remove(uint(index))
	if !success {
		http.Error(w, "Index out of range", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Remove successful"})
}

func handleFind(w http.ResponseWriter, r *http.Request, list *SafeLinkedList) {
	valueStr := strings.TrimPrefix(r.URL.Path, "/find/")
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		http.Error(w, "Invalid value", http.StatusBadRequest)
		return
	}

	index, found := list.Find(value)
	if !found {
		http.Error(w, "Value not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]uint{"index": index})
}

func handleList(w http.ResponseWriter, _ *http.Request, list *SafeLinkedList) {
	list.mutex.Lock()
	defer list.mutex.Unlock()

	values := list.list.HandleList()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(values)
}

func V1() http.Handler {
	list := NewSafeLinkedList()

	h := http.NewServeMux()

	h.HandleFunc("POST /insert", func(w http.ResponseWriter, r *http.Request) {
		handleInsert(w, r, list)
	})
	h.HandleFunc("GET /get/{index}", func(w http.ResponseWriter, r *http.Request) {
		handleGet(w, r, list)
	})
	h.HandleFunc("DELETE /remove/{index}", func(w http.ResponseWriter, r *http.Request) {
		handleRemove(w, r, list)
	})
	h.HandleFunc("GET /find/{value}", func(w http.ResponseWriter, r *http.Request) {
		handleFind(w, r, list)
	})
	h.HandleFunc("GET /list", func(w http.ResponseWriter, r *http.Request) {
		handleList(w, r, list)
	})

	return h
}
