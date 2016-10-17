package main

import (
	"encoding/xml"
	"fmt"
	"net/http"
)

type User struct {
	Name  string "xml:name"
	Email string "xml:email"
	ID    int    "xml:int"
}

func userRouter(w http.ResponseWriter, r *http.Request) {
	outUser := User{}
	outUser.Name = "Bill Smith"
	outUser.Email = "bill@example.com"
	outUser.ID = 100

	output, _ := xml.Marshal(&outUser)
	fmt.Fprintln(w, string(output))
}

func main() {
	fmt.Println("Starting JSON server")
	http.HandleFunc("/user", userRouter)
	http.ListenAndServe(":8080", nil)
}
