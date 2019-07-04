package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
)

type Response struct {
	Message string `json: "message"`
	Status  int    `json: "status_404"`
	IsValid bool   `json: "is_valid"`
}

var Users = struct {
	m map[string]User
	sync.RWMutex
}{m: make(map[string]User)}

type User struct {
	userName string `json: "name"`
}

func hola(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hola Mundo!!!!"))
}
func holaJson(w http.ResponseWriter, r *http.Request) {
	response := CreateResponse("Hola desde Json", 200, true)
	json.NewEncoder(w).Encode(response)
}

func loadHtml(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./Front/index.html")
}

func userExist(user_name string) bool {
	Users.RLock()
	defer Users.RUnlock()

	if _, ok := Users.m[user_name]; ok {
		return true
	}
	return false
}

func validate(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	userName := r.FormValue("user_name")
	fmt.Println(userName)
	response := Response{}

	if userExist(userName) {
		response.IsValid = false
	} else {
		response.IsValid = true
	}

	json.NewEncoder(w).Encode(response)
}

func CreateResponse(message string, status int, isValid bool) Response {
	return Response{message, status, isValid}
}

func main() {
	cssHandle := http.FileServer(http.Dir("./Front/css/"))
	jsHandle := http.FileServer(http.Dir("./Front/js/"))
	mux := mux.NewRouter()
	mux.HandleFunc("/plaintext", hola).Methods("GET")
	mux.HandleFunc("/json", holaJson).Methods("GET")
	mux.HandleFunc("/", loadHtml).Methods("GET")
	mux.HandleFunc("/validate", validate).Methods("POST")
	http.Handle("/", mux)
	http.Handle("/css/", http.StripPrefix("/css/", cssHandle))
	http.Handle("/js/", http.StripPrefix("/js/", jsHandle))
	log.Println("El servidor se encuentra en en el puerto :4545")
	log.Fatal(http.ListenAndServe(":4545", nil))
}
