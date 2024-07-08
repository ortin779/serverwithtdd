package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

func TestRecordingWinsAndRetrivingThem(t *testing.T) {
	store := NewInMemoryStore()
	server := PlayerServer{store}
	player := "Pepper"

	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))
	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))
	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))

	response := httptest.NewRecorder()

	server.ServeHTTP(response, newGetScoreRequest(player))

	assertStatus(t, response.Code, http.StatusOK)

	assertResponseBody(t, response.Body.String(), "3")
}

func TestRecordingWinsAndRetrivingThemConcurrently(t *testing.T) {
	store := NewInMemoryStore()
	server := PlayerServer{store}
	player := "Pepper"
	incrCount := 1000

	wg := sync.WaitGroup{}
	wg.Add(incrCount)

	for i := 0; i < incrCount; i++ {
		go func() {
			server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))
			wg.Done()
		}()
	}

	wg.Wait()

	response := httptest.NewRecorder()

	server.ServeHTTP(response, newGetScoreRequest(player))

	assertStatus(t, response.Code, http.StatusOK)

	assertResponseBody(t, response.Body.String(), "1000")
}

func newPostWinRequest(player string) *http.Request {
	request, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/players/%s", player), nil)
	return request
}
