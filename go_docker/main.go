package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

const defaultContentType = "application/json; charset=utf-8"

type User struct {
	Name    string  `json:"name"`
	Age     int     `json:"age"`
}

func handler(w http.ResponseWriter, r *http.Request) {
	user := User{
		Name:    "taro",
		Age:     18,
	}

	w.Header().Set("content-type", defaultContentType)
	e := json.NewEncoder(w)
	e.SetIndent("", "    ")
	if err := e.Encode(user); err != nil {
		fmt.Fprintf(w, "json encode err: %v\n", err)
		os.Exit(1)
	}
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080",nil)
}
