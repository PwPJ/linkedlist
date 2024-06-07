package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Node represents an element in the linked list
type Node struct {
	Value int
	Next  *Node
}

// L represents a linked list
type L struct {
	head   *Node
	length uint
}

// new creates and returns a new instance of a linked list
func New() *L {
	return &L{}
}

type InsertRequest struct {
	Index uint `json:"index"`
	Value int  `json:"value"`
}

type GetRequest struct {
	Index uint `json:"index"`
}

type RemoveRequest struct {
	Index uint `json:"index"`
}

type FindRequest struct {
	Value int `json:"value"`
}

func (l *L) Find(n int) (index uint, found bool) {
	current := l.head
	index = 0
	for current != nil {
		if current.Value == n {
			return index, true
		}
		current = current.Next
		index++
	}
	return 0, false
}

func (l *L) Get(index uint) (int, bool) {
	current := l.head

	for i := uint(0); i < index; i++ {
		if current == nil {
			return 0, false // return 0 and false if index is out of range
		}
		current = current.Next // move to the next node
	}

	if current == nil {
		return 0, false // return 0 and false if index is exactly one past the last node
	}

	return current.Value, true // return the value at the node and true indicating success
}

func (l *L) Insert(index uint, val int) bool {
	if index > l.length {
		return false
	}

	if index == 0 {
		l.head = &Node{Value: val, Next: l.head}
		l.length++
		return true
	}

	current := l.head
	for i := uint(0); i < index-1; i++ {
		if current == nil {
			return false
		}
		current = current.Next
	}

	current.Next = &Node{Value: val, Next: current.Next}
	l.length++
	return true //successful insertion
}

func (l *L) Remove(index uint) bool {
	if index >= l.length { // Check if index is out of range
		return false
	}

	if index == 0 { // Remove the first element
		l.head = l.head.Next
		l.length--
		return true
	}

	current := l.head
	for i := uint(0); i < index-1; i++ {
		if current.Next == nil {
			return false // If there's no next element, fail
		}
		current = current.Next
	}

	if current.Next == nil {
		return false // Safety check for the last element
	}

	current.Next = current.Next.Next // Adjust the pointer to skip over the removed node
	l.length--
	return true
}

func main() {
	list := New()
	http.HandleFunc("/insert", func(w http.ResponseWriter, r *http.Request) {
		var req InsertRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}
		success := list.Insert(req.Index, req.Value)
		if !success {
			http.Error(w, "Insert failed", http.StatusInternalServerError)
			return
		}
		fmt.Fprintln(w, "Insert successful")
	})

	http.HandleFunc("/get", func(w http.ResponseWriter, r *http.Request) {
		var req GetRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}
		value, found := list.Get(req.Index)
		if !found {
			http.Error(w, "Index out of range", http.StatusNotFound)
			return
		}
		fmt.Fprintf(w, "Value at index %d: %d\n", req.Index, value)
	})

	http.HandleFunc("/remove", func(w http.ResponseWriter, r *http.Request) {
		var req RemoveRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}
		success := list.Remove(req.Index)
		if !success {
			http.Error(w, "Remove failed", http.StatusInternalServerError)
			return
		}
		fmt.Fprintln(w, "Remove successful")
	})

	http.HandleFunc("/find", func(w http.ResponseWriter, r *http.Request) {
		var req FindRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}
		index, found := list.Find(req.Value)
		if !found {
			http.Error(w, "Value not found", http.StatusNotFound)
			return
		}
		fmt.Fprintf(w, "Value %d found at index %d\n", req.Value, index)
	})

	http.HandleFunc("/list", func(w http.ResponseWriter, r *http.Request) {
		current := list.head
		var values []int
		for current != nil {
			values = append(values, current.Value)
			current = current.Next
		}
		json.NewEncoder(w).Encode(values)
	})

	http.ListenAndServe(":8080", nil)
}
