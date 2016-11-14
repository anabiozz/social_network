package api

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"

	Password "social_network/api/password"
	Documentation "social_network/api/specification"
)

func UserCreate(w http.ResponseWriter, r *http.Request) {

	for _, domain := range PermittedDomains {
		fmt.Println("allowing", domain)
		w.Header().Set("Access-Control-Allow-Origin", domain)
	}

	NewUser := User{}
	NewUser.Name = r.FormValue("user")
	NewUser.Email = r.FormValue("email")
	NewUser.First = r.FormValue("first")
	NewUser.Last = r.FormValue("last")
	NewUser.Password = r.FormValue("password")
	salt, hash := Password.ReturnPassword(NewUser.Password)
	fmt.Println(salt, hash)
	output, err := json.Marshal(NewUser)
	fmt.Println(string(output))
	if err != nil {
		fmt.Println("Something went wrong!")
	}

	Response := CreateResponse{}
	sql := "INSERT INTO users set user_nickname='" + NewUser.Name + "', user_first='" + NewUser.First + "', user_last='" + NewUser.Last + "', user_email='" + NewUser.Email + "'" + ", user_password='" + hash + "', user_salt='" + salt + "'"
	q, err := Database.Exec(sql)
	if err != nil {
		errorMessage, errorCode := dbErrorParse(err.Error())
		fmt.Println(errorMessage)
		error, httpCode, msg := ErrorMessages(errorCode)
		Response.Error = msg
		Response.ErrorCode = error
		http.Error(w, "Conflict", httpCode)
	}

	fmt.Println(q)
}

func UsersUpdate(w http.ResponseWriter, r *http.Request) {
	Response := UpdateResponse{}
	params := mux.Vars(r)
	uid := params["id"]
	email := r.FormValue("email")

	var userCount int

	err := Database.QueryRow("SELECT count(user_id) from users where user_id=?", uid).Scan(&userCount)
	if userCount == 0 {
		error, httpCode, msg := ErrorMessages(404)
		log.Println(error)
		log.Println(w, msg, httpCode)
		Response.Error = msg
		Response.ErrorCode = httpCode
		http.Error(w, msg, httpCode)

	} else if err != nil {

	} else {

		_, uperr := Database.Exec("UPDATE users set user_email=? where user_id=?", email, uid)
		if uperr != nil {
			_, errorCode := dbErrorParse(uperr.Error())
			_, httpCode, msg := ErrorMessages(errorCode)

			Response.Error = msg
			Response.ErrorCode = httpCode
			http.Error(w, msg, httpCode)
		} else {
			Response.Error = "success"
			Response.ErrorCode = 0
			output := SetFormat(Response)
			fmt.Fprintln(w, string(output))
		}
	}

}

func UsersRetrieve(w http.ResponseWriter, r *http.Request) {
	log.Println("Starting retrieval")

	accessToken := r.FormValue("access_token")
	if accessToken == "" || CheckToken(accessToken) == false {

	}

	GetFormat(r)
	start := 0
	limit := 10

	next := start + limit

	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Link", "<http://localhost:8080/api/users?start="+string(next)+"; rel=\"next\"")

	rows, _ := Database.Query("select user_id, user_nickname, user_first, user_last, user_email from users LIMIT 10")
	Response := Users{}

	for rows.Next() {

		user := User{}
		rows.Scan(&user.ID, &user.Name, &user.First, &user.Last, &user.Email)
		fmt.Println(user)
		Response.Users = append(Response.Users, user)
	}

	output := SetFormat(Response)
	fmt.Fprintln(w, string(output))
}

func UsersInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Allow", "DELETE,GET,HEAD,OPTIONS,POST,PUT")

	UserDocumentation := []DocMethod{}
	UserDocumentation = append(UserDocumentation, Documentation.UserPOST)
	UserDocumentation = append(UserDocumentation, Documentation.UserOPTIONS)

	output := SetFormat(UserDocumentation)
	fmt.Fprintln(w, string(output))
}
