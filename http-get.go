package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
)

const dbFileName = "game.db.json"


func main() {
	//server := &PlayServer{
	//	store: NewInMemoryPlayerStore(),
	//	router: http.NewServeMux(),
	//}

	db, err := os.OpenFile(dbFileName, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatalf("problem opening %s %v", dbFileName, err)
	}

	store, err := NewFileSystemPlayerStore(db)
	if err != nil {
		log.Fatalf("problem creating file system player store, %v ", err)
	}
	server := NewPlayerServer(store)

	if err := http.ListenAndServe(":5000", server); err != nil {
		log.Fatalf("could not listen on port 5000 %v", err)
	}
}

type PlayStore interface {
	GetPlayScore(name string) int
	RecordWin(name string)
	GetLeague() League
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
	score := p.store.GetPlayScore(player)

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

type FileSystemPlayerStore struct {
	database *json.Encoder
	league League
}

func NewFileSystemPlayerStore(file *os.File) (*FileSystemPlayerStore, error) {
	if err := initialisePlayerDBFile(file); err != nil {
		return nil, fmt.Errorf("problem initialising player db file, %v", err)
	}

	league, err := NewLeague(file)

	if err != nil {
		return nil, fmt.Errorf("problem loading player store from file %s, %v", file.Name(), err)
	}

	return &FileSystemPlayerStore{
		database: json.NewEncoder(&tape{file: file}),
		league:league,
	}, nil
}

func initialisePlayerDBFile(file *os.File) error {
	file.Seek(0, 0)

	info, err := file.Stat()

	if err != nil {
		return fmt.Errorf("problem getting file info from file %s, %v", file.Name(), err)
	}

	if info.Size()==0 {
		file.Write([]byte("[]"))
		file.Seek(0, 0)
	}

	return nil
}

func (f *FileSystemPlayerStore) GetLeague() League {
	sort.Slice(f.league, func(i, j int) bool {
		return f.league[i].Wins > f.league[j].Wins
	})
	return f.league
}

func (f *FileSystemPlayerStore) GetPlayScore(name string) int {
	player := f.league.Find(name)

	if player != nil {
		return player.Wins
	}

	return 0
}

func (f *FileSystemPlayerStore) RecordWin(name string) {
	player := f.league.Find(name)

	if player != nil {
		player.Wins++
	} else {
		f.league = append(f.league, Player{name, 1})
	}

	f.database.Encode(f.league)
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