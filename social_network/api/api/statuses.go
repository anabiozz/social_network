package api

import (
	"fmt"
	"net/http"
)

func StatusCreate(w http.ResponseWriter, r *http.Request) {

	Response := CreateResponse{}
	UserID := r.FormValue("user")
	Status := r.FormValue("status")
	Token := r.FormValue("token")
	ConsumerKey := r.FormValue("consumer_key")

	vUID := ValidateUserRequest(ConsumerKey, Token)
	if vUID != UserID {
		Response.Error = "Invalid user"
		http.Error(w, Response.Error, 401)
		//fmt.Println(w, SetFormat(Response))
	} else {
		_, inErr := Database.Exec("INSERT INTO users_status set user_status_text=?, user_id=?", Status, UserID)
		if inErr != nil {
			fmt.Println(inErr.Error())
			Response.Error = "Error creating status"
			http.Error(w, Response.Error, 500)
			fmt.Fprintln(w, Response)
		} else {
			Response.Error = "Status created"
			fmt.Fprintln(w, Response)
		}
	}

}

func StatusRetrieve(w http.ResponseWriter, r *http.Request) {

}

func StatusDelete(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Nothing to see here")
}

func StatusUpdate(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Coming soon to an API near you!")
}
