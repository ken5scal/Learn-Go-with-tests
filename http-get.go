package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func main() {
	//server := &PlayServer{
	//	store: NewInMemoryPlayerStore(),
	//	router: http.NewServeMux(),
	//}

	server := NewPlayerServer(NewInMemoryPlayerStore())

	if err := http.ListenAndServe(":5000", server); err != nil {
		log.Fatalf("could not listen on port 5000 %v", err)
	}
}

type PlayStore interface {
	GetPlayerScore(name string) int
	RecordWin(name string)
	GetLeague() []Player
}

type PlayServer struct {
	store PlayStore
	http.Handler // embedding ServeHTTP() method
}

func NewPlayerServer(store PlayStore) *PlayServer {
	p := new(PlayServer)

	p.store = store

	router := http.NewServeMux()
	router.HandleFunc("/league", p.leagueHandler)
	router.HandleFunc("/players/", p.playersHandler)

	p.Handler = router

	return p
}

func (p *PlayServer) showScore(w http.ResponseWriter, player string) {
	score := p.store.GetPlayerScore(player)

	if score == 0 {
		w.WriteHeader(http.StatusNotFound)
	}

	fmt.Fprint(w, score)
}

func (p *PlayServer) processWin(w http.ResponseWriter, player string) {
	p.store.RecordWin(player)
	w.WriteHeader(http.StatusAccepted)
}

func (p *PlayServer) leagueHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(p.store.GetLeague())
	w.WriteHeader(http.StatusOK)
}

func (p *PlayServer) getLeagueTable() []Player{
	return []Player{
		{"Chris", 20},
	}
}

func (p *PlayServer) playersHandler(w http.ResponseWriter, r *http.Request) {
	player := r.URL.Path[len("/players/"):]
	fmt.Println(player)

	switch r.Method {
	case http.MethodPost:
		p.processWin(w, player)
	case http.MethodGet:
		p.showScore(w, player)
	}
}

type InMemoryPlayerStore struct {
	store map[string]int
}

func (i *InMemoryPlayerStore)GetPlayerScore(name string) int {
	return i.store[name]
}

func (i *InMemoryPlayerStore) RecordWin(name string) {
	i.store[name]++
}

func (i *InMemoryPlayerStore) GetLeague() []Player{
	var league  []Player
	for name, wins := range i.store {
		league = append(league, Player{name, wins})
	}
	return league
}

func NewInMemoryPlayerStore() *InMemoryPlayerStore {
	return &InMemoryPlayerStore{map[string]int{}}
}

type Player struct {
	Name string
	Wins int
}