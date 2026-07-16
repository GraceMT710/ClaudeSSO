package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"sync"
)

// Todo merepresentasikan satu item task.
type Todo struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

// store adalah in-memory storage dengan mutex untuk thread-safety.
// Untuk MVP ini cukup — kalau nanti butuh persistence, tinggal ganti
// implementasi store ini dengan Firestore/Postgres tanpa mengubah handler.
type store struct {
	mu     sync.Mutex
	nextID int
	todos  map[int]*Todo
}

func newStore() *store {
	return &store{nextID: 1, todos: make(map[int]*Todo)}
}

func (s *store) list() []*Todo {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]*Todo, 0, len(s.todos))
	for _, t := range s.todos {
		out = append(out, t)
	}
	return out
}

func (s *store) create(title string) *Todo {
	s.mu.Lock()
	defer s.mu.Unlock()
	t := &Todo{ID: s.nextID, Title: title}
	s.todos[t.ID] = t
	s.nextID++
	return t
}

func (s *store) toggle(id int) (*Todo, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	t, ok := s.todos[id]
	if !ok {
		return nil, false
	}
	t.Done = !t.Done
	return t, true
}

func (s *store) delete(id int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.todos[id]; !ok {
		return false
	}
	delete(s.todos, id)
	return true
}

// withCORS membungkus handler supaya frontend (beda origin saat dev, misal
// localhost:5173) bisa memanggil API di localhost:8080.
func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func main() {
	s := newStore()
	mux := http.NewServeMux()

	// Go 1.22+ mendukung method-based routing langsung di ServeMux.
	mux.HandleFunc("GET /api/todos", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, s.list())
	})

	mux.HandleFunc("POST /api/todos", func(w http.ResponseWriter, r *http.Request) {
		var body struct{ Title string `json:"title"` }
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Title == "" {
			http.Error(w, "title is required", http.StatusBadRequest)
			return
		}
		writeJSON(w, http.StatusCreated, s.create(body.Title))
	})

	mux.HandleFunc("PATCH /api/todos/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}
		t, ok := s.toggle(id)
		if !ok {
			http.NotFound(w, r)
			return
		}
		writeJSON(w, http.StatusOK, t)
	})

	mux.HandleFunc("DELETE /api/todos/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}
		if !s.delete(id) {
			http.NotFound(w, r)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	})

	log.Println("backend running on :8080")
	log.Fatal(http.ListenAndServe(":8080", withCORS(mux)))
}
