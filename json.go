package main

import (
	"log"
	"encoding/json"
	"net/http"
)

func respondWithError(w http.ResponseWriter, code int, msg string, err error) {
	if err != nil {
		log.Println(err)
	}

	type errStruct struct {
		err  string `json:"error"` 
	}
	
	e := errStruct{err : msg}

	respondWithJson(w, code, e)
}

func respondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	data, err := json.Marshal(payload)
	if err != nil {
		log.Println("Failed to marshal data: ", err)
		w.WriteHeader(500)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}
