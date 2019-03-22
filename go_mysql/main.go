package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

const defaultContentType = "application/json; charset=utf-8"

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func handler(w http.ResponseWriter, r *http.Request) {
	conn, err := Init()
	if err != nil {
		fmt.Fprintf(w, "db init error: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	var users []*User

	rows, err := conn.Query(`
      SELECT 
        id,
        name,
        age 
       from users`)

	if err != nil {
		fmt.Fprintf(w, "get users error: %v\n", err)
		os.Exit(1)
	}

	for rows.Next() {
		var user User
		err := rows.Scan(
			&(user.ID),
			&(user.Name),
			&(user.Age))

		if err != nil {
			fmt.Fprintf(w, "scan users error: %v\n", err)
			os.Exit(1)
		}

		users = append(users, &user)
	}

	w.Header().Set("content-type", defaultContentType)
	e := json.NewEncoder(w)
	e.SetIndent("", "    ")
	if err := e.Encode(users); err != nil {
		fmt.Fprintf(w, "json encode err: %v\n", err)
		os.Exit(1)
	}
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
