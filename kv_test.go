package jstore

import (
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestCreateAndUpdateKV(t *testing.T) {

	key := "key"
	t.Run("string", func(t *testing.T) {

		db, err := New(InMemoryDb, MinimizedJson)
		failOnErr(err, t)
		kv := db.Kv()

		input := "a string"
		err = kv.Set(key, input)
		failOnErr(err, t, "unexpected error setting value: %v")

		input = "second value"
		err = kv.Set(key+"-2", input)
		failOnErr(err, t, "unexpected error setting value: %v")

		updated := "another string"
		err = kv.Set(key, updated)
		failOnErr(err, t, "unexpected error setting value: %v")

		got := ""
		err = kv.Get(key, &got)
		failOnErr(err, t, "unexpected error reading value: %v")

		if diff := cmp.Diff(got, updated); diff != "" {
			t.Errorf("unexpected value (-got +want)\n%s", diff)
		}

		expectedJson := `{"kv":{"key":"another string","key-2":"second value"}}`
		if diff := cmp.Diff(string(db.Json()), expectedJson); diff != "" {
			t.Errorf("unexpected value (-got +want)\n%s", diff)
		}

		if exists := kv.Exists(key); exists != true {
			t.Errorf("expecting key to exits")
		}

		// verify delete
		err = kv.Del(key)
		failOnErr(err, t, "unexpected error deleting key: %v")

		if exists := kv.Exists(key); exists != false {
			t.Errorf("expecting key to NOT exits")
		}

		gotDeleted := ""
		err = kv.Get(key, &gotDeleted)
		failOnErr(err, t, "unexpected error reading value: %v")

		if diff := cmp.Diff(gotDeleted, ""); diff != "" {
			t.Errorf("unexpected value (-got +want)\n%s", diff)
		}
	})

	t.Run("integer", func(t *testing.T) {

		db, err := New(InMemoryDb, MinimizedJson)
		failOnErr(err, t)
		kv := db.Kv("ints")

		input := 100
		err = kv.Set(key, input)
		failOnErr(err, t)

		updated := 200
		err = kv.Set(key, updated)
		failOnErr(err, t)

		got := 0
		err = kv.Get(key, &got)
		failOnErr(err, t)

		if diff := cmp.Diff(got, updated); diff != "" {
			t.Errorf("unexpected value (-got +want)\n%s", diff)
		}

		expectedJson := `{"ints":{"key":200}}`
		if diff := cmp.Diff(string(db.Json()), expectedJson); diff != "" {
			t.Errorf("unexpected value (-got +want)\n%s", diff)
		}
	})

	t.Run("struct", func(t *testing.T) {

		db, err := New(InMemoryDb, MinimizedJson)
		failOnErr(err, t)
		kv := db.Kv("struct")

		input := testPayload{Id: 2}
		err = kv.Set(key, input)
		failOnErr(err, t)

		updated := testPayload{Id: 3}
		err = kv.Set(key, updated)
		failOnErr(err, t)

		got := testPayload{}
		err = kv.Get(key, &got)
		failOnErr(err, t)

		if diff := cmp.Diff(got, updated); diff != "" {
			t.Errorf("unexpected value (-got +want)\n%s", diff)
		}

		expectedJson := `{"struct":{"key":{"id":3,"name":"","sub":{"bolean":false,"time":"0001-01-01T00:00:00Z"}}}}`
		if diff := cmp.Diff(string(db.Json()), expectedJson); diff != "" {
			t.Errorf("unexpected value (-got +want)\n%s", diff)
		}

		if exists := kv.Exists(key); exists != true {
			t.Errorf("expecting key to exits")
		}

		// verify delete
		err = kv.Del(key)
		failOnErr(err, t, "unexpected error deleting key: %v")

		if exists := kv.Exists(key); exists != false {
			t.Errorf("expecting key to NOT exits")
		}

		gotDeleted := testPayload{}
		err = kv.Get(key, &gotDeleted)
		failOnErr(err, t, "unexpected error reading value: %v")

		if diff := cmp.Diff(gotDeleted, testPayload{}); diff != "" {
			t.Errorf("unexpected value (-got +want)\n%s", diff)
		}
	})

	t.Run("slice", func(t *testing.T) {

		db, err := New(InMemoryDb, MinimizedJson)
		failOnErr(err, t)
		kv := db.Kv("slice")

		input := []testPayload{{Id: 1}}
		err = kv.Set(key, input)
		failOnErr(err, t)

		updated := []testPayload{{Id: 3}, {Id: 4}}
		err = kv.Set(key, updated)
		failOnErr(err, t)

		got := []testPayload{}
		err = kv.Get(key, &got)
		failOnErr(err, t)

		if diff := cmp.Diff(got, updated); diff != "" {
			t.Errorf("unexpected value (-got +want)\n%s", diff)
		}

		expectedJson := `{"slice":{"key":[{"id":3,"name":"","sub":{"bolean":false,"time":"0001-01-01T00:00:00Z"}},{"id":4,"name":"","sub":{"bolean":false,"time":"0001-01-01T00:00:00Z"}}]}}`
		if diff := cmp.Diff(string(db.Json()), expectedJson); diff != "" {
			t.Errorf("unexpected value (-got +want)\n%s", diff)
		}
	})

	t.Run("expect error", func(t *testing.T) {
		db, err := New(InMemoryDb)
		failOnErr(err, t)
		kv := db.Kv(t.Name())

		input := "a string"
		err = kv.Set(key, input)
		failOnErr(err, t)

		var got int
		err = kv.Get(key, &got)
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
func TestDeleteKV(t *testing.T) {
	key := "the key"

	db, err := New(InMemoryDb)
	failOnErr(err, t)
	kv := db.Kv(t.Name())

	input := "a string"
	err = kv.Set(key, input)
	failOnErr(err, t)

	// Delete
	err = db.DelCollection(t.Name())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	td3 := testPayload{}
	err = kv.Get(key, &td3)
	if err == nil {
		t.Error("Expecting error but none got")
	} else {
		expect := "collection does not exists"
		if err.Error() != expect {
			t.Errorf("expecting error: %s, but got %v", expect, err)
		}
	}
}
