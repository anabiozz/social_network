package api

import (
	"database/sql"
	"encoding/json"
	"encoding/xml"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/kabukky/httpscerts"

	"log"
	"net/http"
	//"net/url"

	"strconv"
	"strings"

	Password "social_network/api/password"

	"sync"
	"time"
)

var Database *sql.DB
var Routes *mux.Router
var Format string

type UserSession struct {
	ID              string
	GorillaSesssion *sessions.Session
	UID             int
	Expire          time.Time
}

var Session UserSession

func (us *UserSession) Create() {
	us.ID = Password.GenerateSessionID(32)
}

const serverName = "127.0.0.1"
const SSLport = ":8081"
const HTTPport = ":8080"
const SSLprotocol = "https://"
const HTTPprotocol = "http://"

var PermittedDomains []string

type Count struct {
	DBCount int
}

type UpdateResponse struct {
	Error     string "json:error"
	ErrorCode int    "json:code"
}

type CreateResponse struct {
	Error     string "json:error"
	ErrorCode int    "json:code"
}

type Users struct {
	Users []User `json:"users"`
}

type User struct {
	ID       int    "json:id"
	Name     string "json:username"
	Email    string "json:email"
	First    string "json:first"
	Last     string "json:last"
	Password string "json:password"
	Salt     string "json:salt"
	Hash     string "json:hash"
}

type UserDocumentation struct {
}

type OauthAccessResponse struct {
	AccessToken string `json:"access_key"`
}

type Page struct {
	Title        string
	Authorize    bool
	Authenticate bool
	Application  string
	Action       string
	ConsumerKey  string
	Redirect     string
	PageType     string
}

func Init(allowedDomains []string) {
	for _, domain := range allowedDomains {
		PermittedDomains = append(PermittedDomains, domain)
	}
	Routes = mux.NewRouter()
	Routes.HandleFunc("/interface", APIInterface).Methods("GET", "POST", "PUT", "UPDATE")
	Routes.HandleFunc("/api/users", UserCreate).Methods("POST")
	Routes.HandleFunc("/api/users", UsersRetrieve).Methods("GET")
	Routes.HandleFunc("/api/users/{id:[0-9]+}", UsersUpdate).Methods("PUT")
	Routes.HandleFunc("/api/users", UsersInfo).Methods("OPTIONS")
	Routes.HandleFunc("/api/statuses", StatusCreate).Methods("POST")
	Routes.HandleFunc("/api/statuses", StatusRetrieve).Methods("GET")
	Routes.HandleFunc("/api/statuses/{id:[0-9]+}", StatusUpdate).Methods("PUT")
	Routes.HandleFunc("/api/statuses/{id:[0-9]+}", StatusDelete).Methods("DELETE")
	Routes.HandleFunc("/authorize", ApplicationAuthorize).Methods("POST")
	Routes.HandleFunc("/authorize", ApplicationAuthenticate).Methods("GET")
	//Routes.HandleFunc("/authorize/{service:[a-z]+}", ServiceAuthorize).Methods("GET")
	Routes.HandleFunc("/connect/{service:[a-z]+}", ServiceConnect).Methods("GET")
	Routes.HandleFunc("/oauth/token", CheckCredentials).Methods("POST")
}

func CheckLogin(w http.ResponseWriter, r *http.Request) bool {
	cookieSession, err := r.Cookie("sessionid")
	if err != nil {
		fmt.Println("no such cookie")
		Session.Create()
		fmt.Println(Session.ID)
		currTime := time.Now()
		Session.Expire = currTime.Local()
		Session.Expire.Add(time.Hour)

		return false
	} else {
		fmt.Println("found cookki")
		tmpSession := UserSession{UID: 0}
		loggedIn := Database.QueryRow("select user_id from sessions where session_id=?", cookieSession).Scan(&tmpSession.UID)
		if loggedIn != nil {
			return false
		} else {
			if tmpSession.UID == 0 {
				return false
			} else {

				return true
			}
		}
	}
	return false
}

func ServiceConnect(w http.ResponseWriter, r *http.Request) {
	code := r.FormValue("code")
	fmt.Println(code)
}

func redirectNonSecure(w http.ResponseWriter, r *http.Request) {
	log.Println("Non-secure request initiated, redirecting.")
	redirectURL := SSLprotocol + serverName + SSLport + r.RequestURI
	http.Redirect(w, r, redirectURL, http.StatusMovedPermanently)
}

func ValidateUserRequest(cKey string, cToken string) string {
	var UID string
	var aUID string
	var appUID string
	Database.QueryRow("SELECT at.user_id,at.application_user_id,ac.user_id as appuser from api_tokens at left join api_credentials ac on ac.user_id=at.application_user_id where api_token_key=?", cToken).Scan(&UID, &aUID, &appUID)

	return appUID
}

func ErrorMessages(err int64) (int, int, string) {
	errorMessage := ""
	statusCode := 200
	errorCode := 0
	switch err {
	case 1062:
		errorMessage = http.StatusText(409)
		errorCode = 10
		statusCode = http.StatusConflict
	default:
		errorMessage = http.StatusText(int(err))
		errorCode = 0
		statusCode = int(err)
	}

	return errorCode, statusCode, errorMessage

}

func GetFormat(r *http.Request) {

	if len(r.URL.Query()["format"]) > 0 {
		Format = r.URL.Query()["format"][0]
	} else {
		Format = "json"
	}
}

func SetFormat(data interface{}) []byte {

	var apiOutput []byte
	if Format == "json" {
		output, _ := json.Marshal(data)
		apiOutput = output
	} else if Format == "xml" {
		output, _ := xml.Marshal(data)
		apiOutput = output
	} else {
		output, _ := json.Marshal(data)
		apiOutput = output
	}
	return apiOutput
}

func dbErrorParse(err string) (string, int64) {
	Parts := strings.Split(err, ":")
	errorMessage := Parts[1]
	Code := strings.Split(Parts[0], "Error ")
	errorCode, _ := strconv.ParseInt(Code[1], 10, 32)
	return errorMessage, errorCode
}

type DocMethod interface {
}

func CheckToken(token string) bool {
	return true
}

func StartServer() {
	//OauthServices.InitServices()
	fmt.Println(Password.GenerateSalt(22))
	fmt.Println(Password.GenerateSalt(41))

	db, err := sql.Open("mysql", "root:Almera103@/social_network")
	if err != nil {

	}
	Database = db

	wg := sync.WaitGroup{}

	// Check if the cert files are available.
	hs_err := httpscerts.Check("cert.pem", "key.pem")
	// If they are not available, generate new ones.
	if hs_err != nil {
		hs_err = httpscerts.Generate("cert.pem", "key.pem", "127.0.0.1:8081")
		if hs_err != nil {
			log.Fatal("Error: Couldn't create https certs.")
		}
	}

	log.Println("Starting redirection server, try to access @ http:")

	wg.Add(1)
	go func() {
		http.ListenAndServe(HTTPport, http.HandlerFunc(redirectNonSecure))
		wg.Done()
	}()
	wg.Add(1)
	go func() {
		http.ListenAndServeTLS(SSLport, "cert.pem", "key.pem", Routes)
		//http.ListenAndServe(SSLport,http.HandlerFunc(secureRequest))
		wg.Done()
	}()

	wg.Wait()
}
