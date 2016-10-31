package api

import (
	"database/sql"
	"encoding/json"
	"encoding/xml"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"

	Documentation "social_network/api/specification"
)

var Database *sql.DB

var Routes *mux.Router
var Format string

type Count struct {
	DBCount int
}

const (
	serverName   = "localhost"
	SSLport      = ":443"
	HTTPport     = ":8080"
	SSLprotocol  = "https://"
	HTTPprotocol = "http://"
)

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

type Users struct {
	Users []User `json:"users"`
}

type User struct {
	ID    int    "json:id"
	Name  string "json:username"
	Email string "json:email"
	First string "json:first"
	Last  string "json:last"
}

var PermittedDomains []string

type DocMethod interface {
}

func Init(allowedDomains []string) {
	/*for _, domain := range allowedDomains {
		PermittedDomains = append(PermittedDomains, domain)
	}*/
	Routes = mux.NewRouter()
	Routes.HandleFunc("/api/users", UserCreate).Methods("POST")
	Routes.HandleFunc("/api/users", UsersRetrieve).Methods("GET")
	Routes.HandleFunc("/api/users/{id:[0-9]+}", UsersUpdate).Methods("PUT")
	Routes.HandleFunc("/api/users", UsersInfo).Methods("OPTIONS")
}

func UsersInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Allow", "DELETE,GET,HEAD,OPTIONS,POST,PUT")

	UserDocumentation := []DocMethod{}
	UserDocumentation = append(UserDocumentation, Documentation.UserPOST)
	UserDocumentation = append(UserDocumentation, Documentation.UserOPTIONS)

	output := SetFormat(UserDocumentation)
	fmt.Fprintln(w, string(output))
}

type CreateResponse struct {
	Error     string "json:error"
	ErrorCode int    "json:code"
}

type UpdateResponse struct {
	Error     string "json:error"
	ErrorCode int    "json:code"
}

func secureRequest(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "You have arrived at port 443, but you are not yet secure")
}

func redirectNonSecure(w http.ResponseWriter, r *http.Request) {
	log.Println("Non-secure request initiated, redirecting.")
	redirectURL := SSLprotocol + serverName + r.RequestURI
	http.Redirect(w, r, redirectURL, http.StatusMovedPermanently)
}

func dbErrorParse(err string) (string, int64) {
	Parts := strings.Split(err, ":")
	errorMessage := Parts[1]
	Code := strings.Split(Parts[0], "Error ")
	errorCode, _ := strconv.ParseInt(Code[1], 10, 32)
	return errorMessage, errorCode
}

func UserCreate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:9000")
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

	Response := CreateResponse{}
	// Note: This represents a SQL injection vulnerability ... keep reading!
	sql := "INSERT INTO users set user_nickname='" + NewUser.Name + "', user_first='" + NewUser.First + "', user_last='" + NewUser.Last + "', user_email='" + NewUser.Email + "'"
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
	createOutput, _ := json.Marshal(Response)
	fmt.Fprintln(w, string(createOutput))
}

func UsersRetrieve(w http.ResponseWriter, r *http.Request) {
	log.Println("Starting retrieval")
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

func StartServer() {
	/*OauthServices.InitServices()
	fmt.Println(Password.GenerateSalt(22))
	fmt.Println(Password.GenerateSalt(41))*/

	db, err := sql.Open("mysql", "root:Almera103@/social_network")
	if err != nil {

	}
	Database = db

	wg := sync.WaitGroup{}

	log.Println("Starting redirection server, try to access @ http:")

	wg.Add(1)
	go func() {
		http.ListenAndServe(HTTPport, http.HandlerFunc(redirectNonSecure))
		wg.Done()
	}()
	wg.Add(1)
	go func() {
		http.ListenAndServeTLS(SSLport, "cert.pem", "key.pem",
			http.HandlerFunc(secureRequest))
		//http.ListenAndServe(SSLport,http.HandlerFunc(secureRequest))
		wg.Done()
	}()

	wg.Wait()
}

func main() {
	/*http.Handle("/", Routes)
	http.ListenAndServe(":8080", nil)*/
}
