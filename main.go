package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/gorilla/mux"
)

func main() {
	fmt.Print("Hello World")

	// make this cute by making it train themed

}

type App struct {
	Router *mux.Router
	DB     *gorm.DB
}

type AccountManager struct {
	DB *gorm.DB
}

type User struct {
	userName     string
	emailAddress string
}

func (app *App) initializeRoutes() {
	// app.Router.HandleFunc("/products", app.getUsers).Methods("GET")
	app.Router.HandleFunc("/product", app.createUser).Methods("POST")
}

func (am *AccountManager) addUser(user User) (err error) {
	log.Printf("Adding user %s", user.userName)

	// gorm database here
	return
}

func (app *App) createUser(w http.ResponseWriter, r *http.Request) {
	var am AccountManager
	var user User
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&user); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	if err := am.addUser(user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response, _ := json.Marshal(user)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)

}

func (app *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(":8010", app.Router))
}

func (app *App) Initialize() {
	// connectionString :=
	// 	fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", user, password, dbname)

	// var err error
	// a.DB, err = sql.Open("postgres", connectionString)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// https://github.com/go-gorm/postgres
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  "user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Shanghai",
		PreferSimpleProtocol: true, // disables implicit prepared statement usage
	}), &gorm.Config{})

	if err != nil {
		log.Fatal("Unable to open DB: ", err.Error())
	}

	app.DB = db
	app.Router = mux.NewRouter()

	app.initializeRoutes()
}

// func handleRequests() {
// 	http.HandleFunc("/bar", addUser).Methods()
// }

func addUser(w http.ResponseWriter, r *http.Request) {

	var user string
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

}
