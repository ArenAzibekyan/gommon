package configo

import (
	"encoding/json"
	"os"
)

func DecodeJSONFile(name string, v interface{}) error {
	f, err := os.Open(name)
	if err != nil {
		return err
	}
	defer f.Close()
	dec := json.NewDecoder(f)
	return dec.Decode(v)
}
