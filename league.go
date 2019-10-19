package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type League []Player

func (l League) Find(name string) *Player {
	for i, p := range l {
		if p.Name==name {
			return &l[i]
		}
	}
	return nil
}

func NewLeague(rdr io.Reader) ([]Player, error) {
	var league []Player
	err := json.NewDecoder(rdr).Decode(&league)
	if err != nil {
		var league []Player
		err := json.NewDecoder(rdr).Decode(&league)
		if err != nil {
			err = fmt.Errorf("problem parsing league, %v", err)
		}
	}

	return league, err
}

type tape struct {
	file *os.File
}

func (t *tape) Write(p []byte) (n int, err error) {
	t.file.Truncate(0)
	t.file.Seek(0, 0)
	return t.file.Write(p)
}
