package store

import (
	"bytes"
	"encoding/gob"
	"reflect"

	"github.com/pquerna/ffjson/ffjson"
)

func isJSON(s string) bool {
	return (s[0] == '{' && s[len(s)-1] == '}') || (s[0] == '[' || s[len(s)-1] == ']')
}

// Marshal value to bytes using gob or ffjson.
func Marshal(value interface{}) ([]byte, error) {
	switch reflect.ValueOf(value).Kind() {
	case reflect.Ptr, reflect.Struct, reflect.Map:
		return ffjson.Marshal(&Item{
			Object: value,
		})
	default:
		b := &bytes.Buffer{}
		err := gob.NewEncoder(b).Encode(&Item{
			Object: value,
		})

		if err != nil {
			return nil, err
		}

		return b.Bytes(), nil
	}
}

// Unmarshal bytes with gob or ffjson.
func Unmarshal(buf []byte, value interface{}) error {
	if isJSON(string(buf)) {
		return ffjson.Unmarshal(buf, value)
	}

	return gob.NewDecoder(bytes.NewBuffer(buf)).Decode(value)
}