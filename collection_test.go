package jstore

import (
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestCreateAndUpdateCollection(t *testing.T) {

	t.Run("string", func(t *testing.T) {

		db, err := New(InMemoryDb, MinimizedJson)
		failOnErr(err, t)
		col := db.Use(t.Name())

		input := "a string"
		err = col.Set(input)
		failOnErr(err, t)

		updated := "another string"
		err = col.Set(updated)
		failOnErr(err, t)

		got := ""
		err = col.Get(&got)
		failOnErr(err, t)

		if diff := cmp.Diff(got, updated); diff != "" {
			t.Errorf("unexpected value (-got +want)\n%s", diff)
		}

		expectedJson := `{"TestCreateAndUpdateCollection/string":"another string"}`
		if diff := cmp.Diff(string(db.Json()), expectedJson); diff != "" {
			t.Errorf("unexpected value (-got +want)\n%s", diff)
		}
	})

	t.Run("integer", func(t *testing.T) {

		db, err := New(InMemoryDb, MinimizedJson)
		failOnErr(err, t)
		col := db.Use(t.Name())

		input := 100
		err = col.Set(input)
		failOnErr(err, t)

		updated := 200
		err = col.Set(updated)
		failOnErr(err, t)

		got := 0
		err = col.Get(&got)
		failOnErr(err, t)

		if diff := cmp.Diff(got, updated); diff != "" {
			t.Errorf("unexpected value (-got +want)\n%s", diff)
		}

		expectedJson := `{"TestCreateAndUpdateCollection/integer":200}`
		if diff := cmp.Diff(string(db.Json()), expectedJson); diff != "" {
			t.Errorf("unexpected value (-got +want)\n%s", diff)
		}
	})

	t.Run("struct", func(t *testing.T) {

		db, err := New(InMemoryDb, MinimizedJson)
		failOnErr(err, t)
		col := db.Use(t.Name())

		input := testPayload{Id: 2}
		err = col.Set(input)
		failOnErr(err, t)

		updated := testPayload{Id: 3}
		err = col.Set(updated)
		failOnErr(err, t)

		got := testPayload{}
		err = col.Get(&got)
		failOnErr(err, t)

		if diff := cmp.Diff(got, updated); diff != "" {
			t.Errorf("unexpected value (-got +want)\n%s", diff)
		}

		expectedJson := `{"TestCreateAndUpdateCollection/struct":{"id":3,"name":"","sub":{"bolean":false,"time":"0001-01-01T00:00:00Z"}}}`
		if diff := cmp.Diff(string(db.Json()), expectedJson); diff != "" {
			t.Errorf("unexpected value (-got +want)\n%s", diff)
		}
	})

	t.Run("slice", func(t *testing.T) {

		db, err := New(InMemoryDb, MinimizedJson)
		failOnErr(err, t)
		col := db.Use(t.Name())

		input := []testPayload{{Id: 1}}
		err = col.Set(input)
		failOnErr(err, t)

		updated := []testPayload{{Id: 3}}
		err = col.Set(updated)
		failOnErr(err, t)

		got := []testPayload{}
		err = col.Get(&got)
		failOnErr(err, t)

		if diff := cmp.Diff(got, updated); diff != "" {
			t.Errorf("unexpected value (-got +want)\n%s", diff)
		}

		expectedJson := `{"TestCreateAndUpdateCollection/slice":[{"id":3,"name":"","sub":{"bolean":false,"time":"0001-01-01T00:00:00Z"}}]}`
		if diff := cmp.Diff(string(db.Json()), expectedJson); diff != "" {
			t.Errorf("unexpected value (-got +want)\n%s", diff)
		}
	})

	t.Run("expect error", func(t *testing.T) {
		db, err := New(InMemoryDb)
		failOnErr(err, t)
		col := db.Use(t.Name())

		input := "a string"
		err = col.Set(input)
		failOnErr(err, t)

		var got int
		err = col.Get(&got)
		if err == nil {
			t.Error("Expecting error but none got")
		} else {
			expect := "json: cannot unmarshal string into Go value of type int"
			if err.Error() != expect {
				t.Errorf("expecting error: %s, but got %v", expect, err)
			}
		}
	})
}
func TestDeleteCollection(t *testing.T) {
	db, err := New(InMemoryDb)
	failOnErr(err, t)
	col := db.Use(t.Name())

	input := "a string"
	err = col.Set(input)
	failOnErr(err, t)

	// Delete
	err = db.DelCollection(t.Name())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	td3 := testPayload{}
	err = col.Get(&td3)
	if err == nil {
		t.Error("Expecting error but none got")
	} else {
		expect := "collection does not exists"
		if err.Error() != expect {
			t.Errorf("expecting error: %s, but got %v", expect, err)
		}
	}
}
