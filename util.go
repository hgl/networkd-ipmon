package main

import (
	"encoding/json"
	"os"
)

func ReadJSON(name string, v any) error {
	f, err := os.Open(name)
	if err != nil {
		return err
	}
	defer f.Close()
	dec := json.NewDecoder(f)
	return dec.Decode(&v)
}
