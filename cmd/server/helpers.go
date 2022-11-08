package main

import (
	"encoding/json"
	"net/http"
)

type Envelope map[string]interface{}

func (server *Server) JSON(
	w http.ResponseWriter,
	status int,
	data Envelope,
	headers http.Header,
) error {
	json, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	json = append(json, '\n')

	for key, value := range headers {
		w.Header().Set(key, value[0])
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	w.Write(json)

	return nil
}