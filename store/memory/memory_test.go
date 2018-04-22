package memory

import (
	"fmt"
	"reflect"
	"testing"
)

func TestStore(t *testing.T) {
	c := NewStore()

	defer c.Close()

	values := []interface{}{
		"go",
		true,
		[]string{"abc"},
		1,
		1.2,
		[]int{1, 2, 3},
		uint64(3),
		map[string]interface{}{"name": "go"},
	}

	for _, v := range values {
		if err := c.Set("value", v, 0); err != nil {
			t.Fatal(err)
		}

		var ok bool

		r, _ := c.Get("value")

		switch v.(type) {
		case map[string]string:
			a := r.(map[string]interface{})
			b := v.(map[string]interface{})
			ok = true
			for k := range a {
				if !reflect.DeepEqual(a[k], b[k]) {
					ok = false
					break
				}
			}
		default:
			ok = reflect.DeepEqual(r, v)
		}

		if !ok {
			t.Fatal(fmt.Errorf("%v does not match the expected value: %v", r, v))
		}

		if err := c.Remove("value"); err != nil {
			t.Fatal(err)
		}
	}
}

func TestStoreExpired(t *testing.T) {
	c := NewStore()

	defer c.Close()

	if err := c.Set("test", "test", 1); err != nil {
		t.Fatal(err)
	}

	if _, err := c.Get("test"); err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestStoreStruct(t *testing.T) {
	c := NewStore()

	defer c.Close()

	type User struct {
		Name string `json:"name"`
	}

	v := &User{Name: "go"}
	var o *User

	if err := c.Set("struct", v, 0); err != nil {
		t.Fatal(err)
	}

	if err := c.Result("struct", &o); err == nil && o.Name != "go" {
		t.Fatal(fmt.Errorf("User name does not match the expected value: %v", o.Name))
	}

	if err := c.Remove("struct"); err != nil {
		t.Fatal(err)
	}
}
