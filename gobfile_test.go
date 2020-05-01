package httpauth_test

import (
	"os"
	"testing"

	"github.com/lootch/httpauth"
)

// Establish new gobfile for testing due to issues with busy process from previous test.
var gobfile = "gobfile_test.gob"

func TestInitGobFileAuthBackend(t *testing.T) {
	err := os.Remove(gobfile)
	b, err := httpauth.NewGobFileAuthBackend(gobfile)
	if err != httpauth.ErrMissingBackend {
		t.Fatal(err.Error())
	}

	_, err = os.Create(gobfile)
	if err != nil {
		t.Fatal(err.Error())
	}
	b, err = httpauth.NewGobFileAuthBackend(gobfile)
	if err != nil {
		t.Fatal(err.Error())
	}
	if b.Path() != gobfile {
		t.Fatal("File path not saved.")
	}
	if b.Nums() != 0 {
		t.Fatal("Users initialized with items.")
	}

	testBackend(t, b)
}

func TestGobReopen(t *testing.T) {
	b, err := httpauth.NewGobFileAuthBackend(gobfile)
	if err != nil {
		t.Fatal(err.Error())
	}
	b.Close()
	b, err = httpauth.NewGobFileAuthBackend(gobfile)
	if err != nil {
		t.Fatal(err.Error())
	}

	testBackend2(t, b)
}
