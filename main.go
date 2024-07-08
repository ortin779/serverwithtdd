package main

import (
	"log"
	"net/http"
	"sync"
)

func NewInMemoryStore() *InMemoryPlayerStore {
	return &InMemoryPlayerStore{
		scores: map[string]int{},
	}
}

type InMemoryPlayerStore struct {
	mu     sync.Mutex
	scores map[string]int
}

func (i *InMemoryPlayerStore) GetPlayerScore(name string) int {
	return i.scores[name]
}

func (i *InMemoryPlayerStore) RecordWin(name string) {
	i.mu.Lock()
	defer i.mu.Unlock()

	i.scores[name] += 1
}

func main() {
	log.Fatal(http.ListenAndServe(":5000", &PlayerServer{store: NewInMemoryStore()}))
}
