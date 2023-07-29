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

type testPayload struct {
	Id   int     `json:"id"`
	Name string  `json:"name"`
	Sub  SubType `json:"sub"`
}

type SubType struct {
	Bol bool      `json:"bolean"`
	T   time.Time `json:"time"`
}

func TestCRUDSingle(t *testing.T) {
	tcs := []struct {
		name   string
		flags  []DbFlag
		dbName string
	}{
		{
			name:   "in memory only",
			flags:  []DbFlag{},
			dbName: "", // empty string is in memory
		},
		{
			name:   "with write to disk",
			flags:  []DbFlag{},
			dbName: strings.ReplaceAll(t.Name(), "/", "_") + ".json",
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			dbFile := tc.dbName
			if tc.dbName != "" { // empty string is in memory

				dir, err := os.MkdirTemp("", strings.ReplaceAll(t.Name(), "/", "_"))
				if err != nil {
					t.Fatalf("unable to create dir: %s", err)
				}
				defer os.RemoveAll(dir)

				dbFile = filepath.Join(dir, tc.dbName)
			}

			db, err := New(dbFile, HumanReadableJson)
			if err != nil {
				t.Fatalf("error connecton to file: %v", err)
			}

			col := db.Collection("test")

			td := testPayload{
				Id:   1,
				Name: "1",
				Sub: SubType{
					Bol: false,
					T:   time.Time{},
				},
			}

			// Create
			err = col.Write(td)
			if err != nil {
				t.Errorf("unable to write data: %v", err)
			}

			// Read
			td2 := testPayload{}
			err = col.Read(&td2)
			if err != nil {
				t.Errorf("unable to read data: %v", err)
			}

			if diff := cmp.Diff(td2, td); diff != "" {
				t.Errorf("unexpected value (-got +want)\n%s", diff)
			}

			// Update
			td.Id = 3
			td.Name = "3"
			err = col.Write(td)
			if err != nil {
				t.Errorf("unable to write data: %v", err)
			}

			newPayload := testPayload{}
			err = col.Read(&newPayload)
			if err != nil {
				t.Errorf("unable to read data: %v", err)
			}
			if diff := cmp.Diff(newPayload, td); diff != "" {
				t.Errorf("unexpected value (-got +want)\n%s", diff)
			}

			// Delete
			err = db.DelCollection("test")
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			td3 := testPayload{}
			err = col.Read(&td3)
			if err == nil {
				t.Error("Expecting error but none got")
			} else {
				expect := "collection does not exists"
				if err.Error() != expect {
					t.Errorf("expecting error: %s, but got %v", expect, err)
				}
			}
		})
	}
}

func TestCRUDMultiple(t *testing.T) {
	tcs := []struct {
		name   string
		dbName string
	}{
		{
			name:   "in memory only",
			dbName: "", // empty string is in memory
		},
		{
			name:   "with write to disk",
			dbName: strings.ReplaceAll(t.Name(), "/", "_") + ".json",
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			dbFile := tc.dbName
			if tc.dbName != "" { // empty string is in memory

				dir, err := os.MkdirTemp("", strings.ReplaceAll(t.Name(), "/", "_"))
				if err != nil {
					t.Fatalf("unable to create dir: %s", err)
				}
				defer os.RemoveAll(dir)

				dbFile = filepath.Join(dir, tc.dbName)
			}

			db, err := New(dbFile, HumanReadableJson)
			if err != nil {
				t.Fatalf("error connecton to file: %v", err)
			}

			col := db.Collection("test2")

			td := []testPayload{
				{Name: "1"},
				{Name: "2"},
			}

			// Create
			err = col.Write(td)
			if err != nil {
				t.Errorf("unable to write data: %v", err)
			}

			//Read
			td2 := []testPayload{}
			err = col.Read(&td2)
			if err != nil {
				t.Errorf("unable to read data: %v", err)
			}

			if diff := cmp.Diff(td2, td); diff != "" {
				t.Errorf("unexpected value (-got +want)\n%s", diff)
			}

			//Update
			td[0].Id = 3
			td[0].Name = "3"
			err = col.Write(td)
			if err != nil {
				t.Errorf("unable to write data: %v", err)
			}

			newPayload := []testPayload{}
			err = col.Read(&newPayload)
			if err != nil {
				t.Errorf("unable to read data: %v", err)
			}
			if diff := cmp.Diff(newPayload, td); diff != "" {
				t.Errorf("unexpected value (-got +want)\n%s", diff)
			}

			// Delete
			err = db.DelCollection("test2")
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			td3 := []testPayload{}
			err = col.Read(&td3)
			if err == nil {
				t.Error("Expecting error but none got")
			} else {
				expect := "collection does not exists"
				if err.Error() != expect {
					t.Errorf("expecting error: %s, but got %v", expect, err)
				}
			}
		})
	}
}

func TestIsSingle(t *testing.T) {
	item := map[string]any{}

	type a struct {
		a string
		b int
	}
	item["a"] = a{}
	got := isSingle(item["a"])
	if got != true {
		t.Fatalf("expet identify single item")
	}

	item["b"] = []a{
		{
			a: "1",
			b: 0,
		},
		{
			a: "2",
			b: 2,
		},
	}

	got = isSingle(item["b"])
	_ = got
	if got != false {
		t.Fatalf("expet identify  item as multiple ")
	}

}
