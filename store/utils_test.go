package store

import (
	"fmt"
	"reflect"
	"testing"
)

type User struct {
	Name string `json:"name"`
}

func TestMarshalUnmarshal(t *testing.T) {
	values := []interface{}{
		"go",
		true,
		[]string{"abc"},
		1,
		1.2,
		[]int{1, 2, 3},
		uint64(3),
	}

	for _, v := range values {
		b, err := Marshal(v)

		if err != nil {
			t.Fatal(err)
		}

		var i *Item

		if err := Unmarshal(b, &i); err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(i.Object, v) {
			t.Fatal(fmt.Errorf("%v does not match the expected value: %v", i.Object, v))
		}
	}
}

func TestMarshalUnmarshalStruct(t *testing.T) {
	type User struct {
		Name string `json:"name"`
	}

	v := &User{Name: "go"}

	b, err := Marshal(v)
	if err != nil {
		t.Fatal(err)
	}

	var o *User
	if err := Unmarshal(b, &o); err != nil {
		t.Fatal(err)
	}

	if o.Name != "go" {
		t.Fatal(fmt.Errorf("User name does not match the expected value: %v", o.Name))
	}
}
