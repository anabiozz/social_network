package api

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"

	Pseudoauth "social_network/api/pseudoauth"
)

func CheckCredentials(w http.ResponseWriter, r *http.Request) {
	var Credentials string
	Response := CreateResponse{}
	consumerKey := r.FormValue("consumer_key")
	fmt.Println(consumerKey)
	timestamp := r.FormValue("timestamp")
	signature := r.FormValue("signature")
	nonce := r.FormValue("nonce")
	err := Database.QueryRow("SELECT consumer_secret from api_credentials where consumer_key=?", consumerKey).Scan(&Credentials)
	if err != nil {
		error, httpCode, msg := ErrorMessages(404)
		log.Println(error)
		log.Println(w, msg, httpCode)
		Response.Error = msg
		Response.ErrorCode = httpCode
		http.Error(w, msg, httpCode)
		return
	}

	token, err := Pseudoauth.ValidateSignature(consumerKey, Credentials, timestamp, nonce, signature, 0)
	if err != nil {
		error, httpCode, msg := ErrorMessages(401)
		log.Println(error)
		log.Println(w, msg, httpCode)
		Response.Error = msg
		Response.ErrorCode = httpCode
		http.Error(w, msg, httpCode)
		return
	}
	fmt.Println(token)
	AccessRequest := OauthAccessResponse{}
	AccessRequest.AccessToken = token.AccessToken
	output := SetFormat(AccessRequest)
	fmt.Fprintln(w, string(output))
}
