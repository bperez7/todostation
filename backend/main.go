package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/gorilla/mux"
)

func main() {
	fmt.Print("Hello World")
	// router := mux.NewRouter()

	// db, err := DBConnection()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	var app App
	app.Initialize()

	// app.Router = router
	// app.DB = db

	app.initializeRoutes()

	if err := http.ListenAndServe(":8080", app.Router); err != nil {
		log.Fatal(err)
	}

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
	UserName     string `json:"userName"`
	EmailAddress string `json:"emailAddress"`
	UserID       int    `json:"userID"`
}

type TaskRequest struct {
	UserID          int    `json:"userID"`
	TaskName        string `json:"taskName"`
	TaskDescription string `json:"taskDescription"`
	DueDate         string `json:"dueDate"`
	ExpirationDate  string `json:"expirationDate"`
}

type Task struct {
	UserID          int    `json:"userID"`
	TaskName        string `json:"taskName"`
	TaskDescription string `json:"taskDescription"`
	DueDate         string `json:"dueDate"`
	ExpirationDate  string `json:"expirationDate"`
}

func DBConnection() (db *gorm.DB, err error) {
	db, err = gorm.Open(postgres.Open("test.db"), &gorm.Config{})
	if err != nil {
		return db, err
	}
	err = db.AutoMigrate(&User{})
	err = db.AutoMigrate(&Task{})
	err = db.AutoMigrate(&TaskRequest{})

	if err != nil {
		return db, err
	}
	return db, nil
}

func (app *App) getHome(w http.ResponseWriter, r *http.Request) {
	response, _ := json.Marshal("Status OK")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func (app *App) initializeRoutes() {
	log.Print("initializing routes")
	app.Router.HandleFunc("/users", app.getUsers).Methods("GET")
	app.Router.HandleFunc("/users", app.createUser).Methods("POST")
	app.Router.HandleFunc("/{userId}/tasks", app.createTask).Methods("POST")
	app.Router.HandleFunc("/{userId}/tasks", app.getTasksForUser).Methods("GET")
	// app.Router.HandleFunc("/", app.getHome).Methods("GET")
	// http.Handle("/", app.Router)

}

func (am *AccountManager) getDB() (DB *gorm.DB) {
	return am.DB
}

func (am *AccountManager) addUser(user User) (err error) {

	// TODO: find way to generate unique id

	log.Printf("Adding user %s", user.UserName)

	db := am.getDB()
	result := db.Create(&user) // pass pointer of data to Create
	log.Print("Result: ", result)
	// gorm database here
	return
}

func (am *AccountManager) getUsers() (users []User, err error) {
	log.Print("Getting users")
	db := am.getDB()
	log.Print("got db")
	db.Find(&users)

	log.Print("Users:  ", users)

	return

}

func (am *AccountManager) addTask(taskRequest TaskRequest) (err error) {
	log.Printf("Adding task %s", taskRequest.TaskName)
	db := am.getDB()
	task := Task(taskRequest)
	result := db.Create(&task) // pass pointer of data to Create
	log.Print("Result: ", result)

	return

}

func (am *AccountManager) getTasks(userID int) (tasks []Task, err error) {
	log.Printf("Getting tasks from user %d", userID)
	db := am.getDB()
	db.Where("user_id = ?", userID).Find(&tasks) // pass pointer of data to Create

	log.Print("Tasks:  ", tasks)

	return

}

func (app *App) createUser(w http.ResponseWriter, r *http.Request) {
	var am AccountManager
	am.DB = app.DB
	var user User
	log.Print(r.Body)
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&user); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	log.Print("create user: ", user)

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

func (app *App) getUsers(w http.ResponseWriter, r *http.Request) {
	var am AccountManager
	am.DB = app.DB
	users, err := am.getUsers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response, _ := json.Marshal(users)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)

}

func (app *App) createTask(w http.ResponseWriter, r *http.Request) {
	var am AccountManager
	am.DB = app.DB
	var taskRequest TaskRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&taskRequest); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	if err := am.addTask(taskRequest); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response, _ := json.Marshal(taskRequest)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func (app *App) getTasksForUser(w http.ResponseWriter, r *http.Request) {
	var am AccountManager
	am.DB = app.DB
	routeVars := mux.Vars(r)

	userID := routeVars["userId"]

	defer r.Body.Close()
	log.Print("Getting tasks for user ", userID)
	userIDInteger, err := strconv.Atoi(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tasks, err := am.getTasks(userIDInteger)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response, _ := json.Marshal(tasks)

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
		DSN:                  "user=brandonperez password=gorm dbname=postgres port=5432 sslmode=disable TimeZone=Asia/Shanghai",
		PreferSimpleProtocol: true, // disables implicit prepared statement usage
	}), &gorm.Config{})

	err = db.AutoMigrate(&User{})
	err = db.AutoMigrate(&Task{})
	err = db.AutoMigrate(&TaskRequest{})

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
