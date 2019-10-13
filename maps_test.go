package main

import (
	"testing"
)

func TestSearch(t *testing.T) {
	d := Dictionary{"test": "this is just a test"}

	t.Run("known owrd", func(t *testing.T) {
		got, _ := d.Search( "test")
		want := "this is just a test"

		assertStrings(t, got, want)
	})

	t.Run("unknown owrd", func(t *testing.T) {
		_, err := d.Search( "unknown")
		assertError(t, err, ErrNotFound)
	})
}

func TestAdd(t *testing.T) {
	d := Dictionary{}
	d.Add("test", "this is just a test")
	want := "this is just a test"
	if got , err := d.Search("test"); err != nil {
		t.Fatal("should find added word:", err)
	} else if want != got {
		t.Errorf("got %q want %q", got ,want)
	}

	t.Run("new word", func(t *testing.T) {
		d:= Dictionary{}
		w := "test"
		def := "this is just a test"
		err := d.Add(w, def)
		assertError2(t, err, nil)
		assertDefinition(t, d, w, def)
	})

	t.Run("existing word", func(t *testing.T) {
		w := "test"
		def := "this is just a test"
		d:= Dictionary{w: def}
		err := d.Add(w, "new test")
		assertError2(t, err, ErrWordExists)
		assertDefinition(t, d, w, def)
	})
}

func TestUpdate(t *testing.T) {
	t.Run("existing word", func(t *testing.T) {
		w := "test"
		def := "this is just a test"
		newDef := "new def"
		d:= Dictionary{w: def}
		err := d.Update(w, newDef)
		assertError2(t, err, nil)
		assertDefinition(t, d, w, newDef)
	})

	t.Run("new word", func(t *testing.T) {
		w := "test"
		def := "this is just a test"
		d:= Dictionary{}
		err := d.Update(w, def)
		assertError2(t, err, ErrWordDoesNotExist)
	})
}

func assertStrings(t *testing.T, got string, want string) {
	t.Helper()

	if got != want {
		t.Errorf("got %q want %q give, %q", got ,want, "test")
	}
}

func assertError2(t *testing.T, got error, want error) {
	t.Helper()

	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func assertDefinition(t *testing.T, d Dictionary, word, definigion string) {
	t.Helper()
	if got , err := d.Search(word); err != nil {
		t.Fatal("should find added word:", err)
	} else if definigion != got {
		t.Errorf("got %q want %q", got , definigion)
	}

}