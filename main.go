package main

import (
	"github.com/gorilla/mux"
	"github.com/perseph0ne/DSS-UI-WEB/controller"
	"github.com/perseph0ne/DSS-UI-WEB/model"
	"html/template"
	"log"
	"net/http"
)

var users = model.GenerateUser()
var userLogged = model.User{ID: 0, UserName: "", Admin: false}

func main() {
	router := mux.NewRouter()
	//load files--------------------
	fsAsset := http.FileServer(http.Dir("./assets/"))
	fsView := http.FileServer(http.Dir("./view/"))
	router.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", fsAsset))
	router.PathPrefix("/view/").Handler(http.StripPrefix("/view/", fsView))
	//-------------//------------------

	router.HandleFunc("/", index).Methods(http.MethodGet)
	router.HandleFunc("/", index).Methods(http.MethodPost)
	router.HandleFunc("/users", getListUsers).Methods(http.MethodGet)
	router.HandleFunc("/user/{id}", deleteUser).Methods(http.MethodGet)
	router.HandleFunc("/user", registerUser).Methods(http.MethodPost)
	router.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
		appResult := model.AppResult{UserLogged: userLogged, Admin: userLogged.Admin}
		tmpl, _ := template.ParseFiles("./view/usercreate.html")
		_ = tmpl.Execute(w, appResult)
	}).Methods(http.MethodGet)
	router.HandleFunc("/documents", getListDocuments).Methods(http.MethodGet)
	router.HandleFunc("/document", createDocument).Methods(http.MethodPost)
	router.HandleFunc("/document/{id}", deleteDocument).Methods(http.MethodGet)
	router.HandleFunc("/login", login).Methods(http.MethodPost)
	router.HandleFunc("/logout", logout).Methods(http.MethodGet)
	log.Fatal(http.ListenAndServe(":9001", router))
}

func index(w http.ResponseWriter, r *http.Request) {
	if userLogged.ID == 0 {
		appResult := model.AppResult{UserLogged: userLogged, Admin: false}
		tmpl, _ := template.ParseFiles("./view/login.html")
		_ = tmpl.Execute(w, appResult)
	} else {
		appResult := model.AppResult{UserLogged: userLogged, Admin: userLogged.Admin}
		tmpl, _ := template.ParseFiles("./view/index.html")
		_ = tmpl.Execute(w, appResult)
	}
}
func login(w http.ResponseWriter, r *http.Request) {
	controller.LoginUser(w, r, &userLogged, users)
}
func logout(w http.ResponseWriter, r *http.Request) {
	userLogged = model.User{ID: 0, UserName: "", Admin: false}
	controller.LogoutUser(w, r, userLogged)
}
func getListUsers(w http.ResponseWriter, r *http.Request) {
	if userLogged.ID > 0 {
		controller.GetUsers(w, r, userLogged, users)
	} else {
		appResult := model.AppResult{UserLogged: userLogged, Admin: userLogged.Admin}
		tmpl, _ := template.ParseFiles("./view/index.html")
		_ = tmpl.Execute(w, appResult)
	}

}
func registerUser(w http.ResponseWriter, r *http.Request) {
	controller.CreateUser(w, r, &userLogged, &users)
}
func deleteUser(w http.ResponseWriter, r *http.Request) {
	controller.DeleteUserById(w, r, &userLogged, &users)
}
func getListDocuments(w http.ResponseWriter, r *http.Request) {
	controller.GetDocumentsServer(w, r, userLogged)
}
func createDocument(w http.ResponseWriter, r *http.Request) {
	controller.UploadDocumentServer(w, r, userLogged, users)
}
func deleteDocument(w http.ResponseWriter, r *http.Request) {
	controller.DeleteDocumentServer(w, r, userLogged, users)
}
