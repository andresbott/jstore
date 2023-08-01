package jstore

//
//import (
//	"os"
//	"path/filepath"
//	"strings"
//	"testing"
//)
//
//func TestKV(t *testing.T) {
//
//	tcs := []struct {
//		name   string
//		flags  []DbFlag
//		dbName string
//	}{
//		{
//			name:   "in memory only",
//			flags:  []DbFlag{},
//			dbName: "", // empty string is in memory
//		},
//		{
//			name:   "with write to disk",
//			flags:  []DbFlag{},
//			dbName: strings.ReplaceAll(t.Name(), "/", "_") + ".json",
//		},
//	}
//
//	for _, tc := range tcs {
//		t.Run(tc.name, func(t *testing.T) {
//			dbFile := tc.dbName
//			if tc.dbName != "" { // empty string is in memory
//
//				dir, err := os.MkdirTemp("", strings.ReplaceAll(t.Name(), "/", "_"))
//				if err != nil {
//					t.Fatalf("unable to create dir: %s", err)
//				}
//				defer os.RemoveAll(dir)
//
//				dbFile = filepath.Join(dir, tc.dbName)
//			}
//
//			db, err := New(dbFile, HumanReadableJson)
//			if err != nil {
//				t.Fatalf("error connecton to file: %v", err)
//			}
//
//			kv := db.Kv("")
//			itemKey := "item1"
//
//			err = kv.Set(itemKey, "banana")
//			if err != nil {
//				t.Fatalf("unexpected error: %v", err)
//			}
//
//			got, err := kv.Get(itemKey)
//			if err != nil {
//				t.Fatalf("unexpected error: %v", err)
//			}
//			if got != "banana" {
//				t.Errorf("got unexpected value \"%s\"", got)
//			}
//
//			// overwrite the value
//			err = kv.Set(itemKey, "apple")
//			if err != nil {
//				t.Fatalf("unexpected error: %v", err)
//			}
//			got, err = kv.Get(itemKey)
//			if err != nil {
//				t.Fatalf("unexpected error: %v", err)
//			}
//			if got != "apple" {
//				t.Errorf("got unexpected value \"%s\"", got)
//			}
//
//			// unset the item
//			kv.Set(itemKey, "")
//			got, err = kv.Get(itemKey)
//			if err != nil {
//				t.Fatalf("unexpected error: %v", err)
//			}
//			if got != "" {
//				t.Errorf("got unexpected value %s, expoecting \"\"", got)
//			}
//
//			// get a value that does not exits
//			got, err = kv.Get("not Exist")
//			if err != nil {
//				t.Fatalf("unexpected error: %v", err)
//			}
//			if got != "" {
//				t.Errorf("got unexpected value %s, expoecting \"\"", got)
//			}
//
//		})
//	}
//
//}
