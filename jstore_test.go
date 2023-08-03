package jstore

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/google/go-cmp/cmp"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestCreateDb(t *testing.T) {
	_ = spew.Config // keep the spew package
	dir, err := os.MkdirTemp("", t.Name())
	if err != nil {
		t.Fatalf("unable to create dir: %s", err)
	}
	defer os.RemoveAll(dir)

	t.Run("file in dir", func(t *testing.T) {

		dbFile := filepath.Join(dir, strings.ReplaceAll(t.Name(), "/", "_")+".json")
		_, e := New(dbFile)
		if e != nil {
			t.Errorf("error connecton to file: %v", e)
		}
	})

	t.Run("dir", func(t *testing.T) {
		_, e := New(dir)
		expect := fmt.Sprintf("open %s: is a directory", dir)
		if e == nil {
			t.Errorf("expecting an error, but none got")
		} else {
			if e.Error() != expect {
				t.Errorf("expecting error: %s, but got: %s", expect, e)
			}
		}
	})
}
func failOnErr(err error, t *testing.T, tmpl ...string) {
	template := "unexpected err: %v"
	if len(tmpl) > 0 && tmpl[0] != "" {
		template = tmpl[0]
	}

	if err != nil {
		t.Fatalf(template, err)
	}
}

type testPayload struct {
	Id   int     `json:"id"`
	Name string  `json:"name"`
	Sub  SubType `json:"sub"`
}

type SubType struct {
	Bol bool      `json:"bolean"`
	T   time.Time `json:"time"`
}

func TestStoreBackends(t *testing.T) {

	payload := testPayload{
		Id:   2,
		Name: "2",
		Sub: SubType{
			Bol: true,
			T:   time.Time{},
		}}

	t.Run("in memory", func(t *testing.T) {

		db, err := New(InMemoryDb, MinimizedJson)
		failOnErr(err, t)

		col := db.Use("test")

		err = col.Set(payload)
		failOnErr(err, t)

		got := testPayload{}
		err = col.Get(&got)
		failOnErr(err, t)

		if diff := cmp.Diff(got, payload); diff != "" {
			t.Errorf("unexpected value (-got +want)\n%s", diff)
		}

	})

	t.Run("with file backend", func(t *testing.T) {
		dir, err := os.MkdirTemp("", strings.ReplaceAll(t.Name(), "/", "_"))
		failOnErr(err, t)
		defer os.RemoveAll(dir)
		dbFile := filepath.Join(dir, strings.ReplaceAll(t.Name(), "/", "_")+".json")

		db, err := New(dbFile, MinimizedJson)
		failOnErr(err, t)

		col := db.Use("test")

		err = col.Set(payload)
		failOnErr(err, t)

		got := testPayload{}
		err = col.Get(&got)
		failOnErr(err, t)

		if diff := cmp.Diff(got, payload); diff != "" {
			t.Errorf("unexpected value (-got +want)\n%s", diff)
		}
	})
}

// todo: write tests around manual flush
