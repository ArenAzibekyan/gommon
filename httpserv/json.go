package httpserv

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

func DecodeJSONReq(r *http.Request, v interface{}) error {
	if r.Header.Get("Content-Type") != "application/json" {
		return errors.New("'Content-Type' header must be 'application/json'")
	}
	data, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}

func EncodeJSONResp(w http.ResponseWriter, v interface{}) error {
	w.Header().Add("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}
