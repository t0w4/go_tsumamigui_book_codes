package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

const defaultContentType = "application/json; charset=utf-8"

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

// UserService finds users.
type UserService struct {
	db *sql.DB
}

func (us *UserService) handler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	conn, err := us.db.Conn(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "db init error: %v\n", err)
		return
	}
	defer conn.Close()

	var users []*User

	rows, err := conn.QueryContext(ctx, `
      SELECT 
        id,
        name,
        age 
       from users`)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "get users error: %v\n", err)
		return
	}

	for rows.Next() {
		var user User
		err := rows.Scan(
			&(user.ID),
			&(user.Name),
			&(user.Age))

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "scan users error: %v\n", err)
			return
		}

		users = append(users, &user)
	}

	w.Header().Set("content-type", defaultContentType)
	e := json.NewEncoder(w)
	e.SetIndent("", "    ")
	if err := e.Encode(users); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "json encode err: %v\n", err)
	}
}

func main() {
	db, err := Init()
	if err != nil {
		log.Fatalf("fail: %+v", err)
	}
	us := &UserService{db}
	http.HandleFunc("/", us.handler)
	http.ListenAndServe(":8080", nil)
}
