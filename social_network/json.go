package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type User struct {
	Name  string "json:name"
	Email string "json:email"
	ID    int    "json:int"
}

func userRouter(w http.ResponseWriter, r *http.Request) {
	outUser := User{}
	outUser.Name = "Bill Smith"
	outUser.Email = "bill@example.com"
	outUser.ID = 100

	output, _ := json.Marshal(&outUser)
	fmt.Fprintln(w, string(output))
}

func GetUser(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Pragma", "no-cache")

	urlParams := mux.Vars(r)
	id := urlParams["id"]
	ReadUser := User{}
	err := database.QueryRow("select * from users where user_id=?", id).Scan(&ReadUser.ID, &ReadUser.Name, &ReadUser.First, &ReadUser.Last, &ReadUser.Email)
	switch {
	case err == sql.ErrNoRows:
		fmt.Fprintf(w, "No such user")
	case err != nil:
		log.Fatal(err)
	default:
		output, _ := json.Marshal(ReadUser)
		fmt.Fprintf(w, string(output))
	}
}

func main() {
	fmt.Println("Starting JSON server")
	http.HandleFunc("/user", userRouter)
	http.ListenAndServe(":8080", nil)
}
