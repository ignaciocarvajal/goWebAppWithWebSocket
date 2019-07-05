package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
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
	userName  string `json: "name"`
	Websocket *websocket.Conn
}

func createUser(userName string, ws *websocket.Conn) User {
	return User{userName, ws}
}

func AddUser(user User) {
	Users.Lock()
	defer Users.Unlock()
	Users.m[user.userName] = user
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

func removeUser(userName string) {
	Users.Lock()
	defer Users.Unlock()
	delete(Users.m, userName)
}

func sendMessage(typeMessage int, message []byte) {
	Users.RLock()
	defer Users.RUnlock()

	for _, user := range Users.m {
		if err := user.Websocket.WriteMessage(typeMessage, message); err != nil {
			return
		}
	}
}

func toArrayByte(value string) []byte {
	return []byte(value)
}

func concatMessage(userName string, array []byte) string {
	return userName + " : " + string(array[:])
}

func webSocket(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userName := vars["user_name"]
	ws, err := websocket.Upgrade(w, r, nil, 1024, 1024)
	if err != nil {
		log.Println(err)
		return
	}
	currentUser := createUser(userName, ws)
	AddUser(currentUser)
	log.Println("Nuevo Ususario agregado")

	for {
		typeMessage, message, err := ws.ReadMessage()
		if err != nil {
			removeUser(userName)
			return
		}
		finalMessage := concatMessage(userName, message)
		sendMessage(typeMessage, toArrayByte(finalMessage))
	}
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
	mux.HandleFunc("/chat/{user_name}", webSocket).Methods("GET")
	http.Handle("/", mux)
	http.Handle("/css/", http.StripPrefix("/css/", cssHandle))
	http.Handle("/js/", http.StripPrefix("/js/", jsHandle))
	log.Println("El servidor se encuentra en en el puerto :4545")
	log.Fatal(http.ListenAndServe(":4545", nil))
}
