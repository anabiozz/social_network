package api

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"html/template"
	"net/http"

	Password "social_network/api/password"
	Pseudoauth "social_network/api/pseudoauth"
	//OauthServices "social_network/api/services"
)

func ApplicationAuthorize(w http.ResponseWriter, r *http.Request) {

	CheckLogin(w, r)

	username := r.FormValue("username")
	password := r.FormValue("password")
	allow := r.FormValue("authorize")
	authType := r.FormValue("auth_type")
	redirect := r.FormValue("redirect")

	var dbPassword string
	var dbSalt string
	var dbUID string

	uerr := Database.QueryRow("SELECT user_password, user_salt, user_id from users where user_nickname=?", username).Scan(&dbPassword, &dbSalt, &dbUID)
	if uerr != nil {

	}

	consumerKey := r.FormValue("consumer_key")

	var CallbackURL string
	var appUID string
	if authType == "client" {
		err := Database.QueryRow("SELECT user_id,callback_url from api_credentials where consumer_key=?", consumerKey).Scan(&appUID, &CallbackURL)
		if err != nil {
			fmt.Println("SELECT user_id,callback_url from api_credentials where consumer_key=?", consumerKey)
			fmt.Println(err.Error())
			return
		}
	}

	expectedPassword := Password.GenerateHash(dbSalt, password)
	fmt.Println("allow:", allow)
	fmt.Println("authtype:", authType)
	fmt.Println(dbPassword, "=", expectedPassword)
	if (dbPassword == expectedPassword) && (allow == "1") && (authType == "consumer") {
		fmt.Println("Yes!")
		requestToken := Pseudoauth.GenerateToken()

		authorizeSQL := "INSERT INTO api_tokens set application_user_id=?, user_id=?, api_token_key=?"

		q, connectErr := Database.Exec(authorizeSQL, appUID, dbUID, requestToken)
		if connectErr != nil {
			fmt.Println(connectErr.Error())
		} else {
			fmt.Println(q)
		}
		redirectURL := CallbackURL + "?request_token=" + requestToken
		fmt.Println(redirectURL)
		http.Redirect(w, r, redirectURL, http.StatusAccepted)

	} else if (dbPassword == expectedPassword) && authType == "user" {

		_, err := Database.Exec("insert into sessions set session_id=?,user_id=?", Session.ID, dbUID)
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Println("redirecting")
		http.Redirect(w, r, redirect, http.StatusOK)
	} else {
		fmt.Println(authType)
		fmt.Println(dbPassword, expectedPassword)
		http.Redirect(w, r, "/authorize", http.StatusUnauthorized)
	}

}

func ApplicationAuthenticate(w http.ResponseWriter, r *http.Request) {

	Authorize := Page{}
	Authorize.Authenticate = true
	Authorize.Title = "Login"
	Authorize.Application = ""
	Authorize.Action = "/authorize"
	if len(r.URL.Query()["consumer_key"]) > 0 {
		Authorize.ConsumerKey = r.URL.Query()["consumer_key"][0]
	} else {
		Authorize.ConsumerKey = ""
	}
	if len(r.URL.Query()["redirect"]) > 0 {
		Authorize.Redirect = r.URL.Query()["redirect"][0]
	} else {
		Authorize.Redirect = ""
	}

	if Authorize.ConsumerKey == "" && Authorize.Redirect != "" {
		Authorize.PageType = "user"
	} else {
		Authorize.PageType = "consumer"
	}

	tpl := template.Must(template.New("main").ParseFiles("authorize.html"))
	tpl.ExecuteTemplate(w, "authorize.html", Authorize)
}

/*func ServiceAuthorize(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	service := params["service"]

	loggedIn := CheckLogin(w, r)
	if loggedIn == false {
		Cookie := http.Cookie{Name: "sessionid", Value: Session.ID, Expires: Session.Expire}
		fmt.Println("Setting cookie!")
		http.SetCookie(w, &Cookie)
		redirect := url.QueryEscape("/authorize/" + service)
		http.Redirect(w, r, "/authorize?redirect="+redirect, http.StatusUnauthorized)
		return
	}

	redURL := OauthServices.GetAccessTokenURL(service, "")
	http.Redirect(w, r, redURL, http.StatusFound)

}*/
