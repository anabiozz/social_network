package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

var database *sql.DB

type Users struct {
	Users []User "json:users"
}

type User struct {
	ID    int    "json:id"
	Name  string "json:username"
	Email string "json:email"
	First string "json:first"
	Last  string "json:last"
}

type CreateResponce struct {
	Error string "json:error"
}

func UserCreate(w http.ResponseWriter, r *http.Request) {

	NewUser := User{}
	NewUser.Name = r.FormValue("user")
	NewUser.Email = r.FormValue("email")
	NewUser.First = r.FormValue("first")
	NewUser.Last = r.FormValue("last")
	output, err := json.Marshal(NewUser)
	fmt.Println(string(output))
	if err != nil {
		fmt.Println("Something went wrong!")
	}

	Responce := CreateResponce{}
	sql := "INSERT INTO users set user_nickname='" + NewUser.Name + "', user_first='" + NewUser.First + "', user_last='" + NewUser.Last + "', user_email='" + NewUser.Email + "'"
	q, err := database.Exec(sql)
	if err != nil {
		Responce.Error = err.Error()
	}
	fmt.Println(q)
	createOutput, _ := json.Marshal(Responce)
	fmt.Fprintln(w, string(createOutput))
}

func UsersRetrive(w http.ResponseWriter, r *http.Request) {
	/*urlParams := mux.Vars(r)
	key := vars["key"]
	*/
	log.Println("starting retrieval")
	start := 0
	limit := 10

	next := start + limit

	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Link", "<http://localhost:8080/api/users?start="+string(next)+"; rel=\"next\"")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:9000")

	rows, _ := database.Query("select * from users LIMIT 10")
	Responce := Users{}

	for rows.Next() {
		user := User{}
		rows.Scan(&user.ID, &user.Name, &user.First, &user.Last, &user.Email)
		Responce.Users = append(Responce.Users, user)
	}
	output, _ := json.Marshal(Responce)
	fmt.Fprintln(w, string(output))
}

//{key:[A-Za-z0-9\-]
func main() {
	db, err := sql.Open("mysql", "root:Almera103@/social_network")
	if err != nil {

	}
	database = db
	routes := mux.NewRouter()
	routes.HandleFunc("/api/users", UserCreate).Methods("POST")
	routes.HandleFunc("/api/users", UsersRetrive).Methods("GET")
	http.Handle("/", routes)
	http.ListenAndServe(":8080", nil)
}
