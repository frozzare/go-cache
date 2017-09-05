package redis

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
)

func TestRedis(t *testing.T) {
	c := NewStore(nil)

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

func TestRedisIncrementDecrement(t *testing.T) {
	c := NewStore(nil)

	if err := c.Remove("num"); err != nil {
		t.Fatal(err)
	}

	if _, err := c.Increment("num"); err != nil {
		t.Fatal(err)
	}

	if r, err := c.Number("num"); r != 1 {
		fmt.Println(err)
		t.Fatal(fmt.Errorf("%v does not match the expected value: %v", r, 1))
	}

	if _, err := c.Decrement("num"); err != nil {
		t.Fatal(err)
	}

	if r, _ := c.Number("num"); r != 0 {
		t.Fatal(errors.New("Num does not match expected number"))
	}

	if err := c.Remove("num"); err != nil {
		t.Fatal(err)
	}
}

func TestRedisStruct(t *testing.T) {
	c := NewStore(nil)

	type User struct {
		Name string `json:"name"`
	}

	v := &User{Name: "go"}
	var o *User

	if err := c.Set("struct", v, 0); err != nil {
		t.Fatal(err)
	}

	if err := c.Result("struct", &o); err == nil && o.Name == "go" {
		t.Fatal(fmt.Errorf("User name does not match the expected value: %v", o.Name))
	}

	if err := c.Remove("struct"); err != nil {
		t.Fatal(err)
	}
}
