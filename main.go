package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type message struct {
	Title string            `json:"title"`
	Data  map[string]string `json:"data"`
}

type testResponse struct {
	Message string            `json:"message"`
	Data    map[string]string `json:"data"`
}

func main() {
	wantSecret := os.Getenv("WEBHOOK_SECRET")

	http.HandleFunc("/hook", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		gotSecret := r.Header.Get("X-hook-secret")
		if gotSecret != wantSecret {
			http.Error(w, "Invalid Secret", http.StatusForbidden)
			return
		}

		isTest := r.Header.Get("X-test") != ""
		if isTest {
			var buffer []byte
			var err error

			defer r.Body.Close()
			if buffer, err = ioutil.ReadAll(r.Body); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			var msg message
			if err := json.Unmarshal(buffer, &msg); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			resp := testResponse{
				Message: "Received submission for form " + msg.Title,
				Data:    msg.Data,
			}

			if buffer, err = json.Marshal(resp); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("content-type", "application/json")
			w.Write(buffer)
			return
		}

		// Actual handling is not yet implemented
		http.Error(w, "Not Implemented", http.StatusNotImplemented)
	})

	http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("PORT")), nil)
}
