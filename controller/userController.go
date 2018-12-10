package controller

import (
	"github.com/gorilla/mux"
	"github.com/perseph0ne/DSS-UI-WEB/model"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func GetUsers(w http.ResponseWriter, r *http.Request,userLogged model.User, users []model.User){
	tmpl,_ := template.ParseFiles("./view/user.html")
	appResult :=model.AppResultUser{UserLogged:userLogged,Admin:userLogged.Admin,MsgResult:"",Users:users}
	_ = tmpl.Execute(w, appResult)
}

func CreateUser(w http.ResponseWriter, r *http.Request,userLogged *model.User, users *[]model.User){

	userName := r.FormValue("username")
	password := r.FormValue("password")
	email := r.FormValue("email")
	listUser:= *users
	id:= listUser[len(listUser)-1].ID +1
	userNew := model.User{ID:id,UserName:userName,Password:password,Email:email,CreatedAT:time.Now()}
	*users= append(*users,userNew)
	tmpl,_ := template.ParseFiles("./view/user.html")
	appResult :=model.AppResultUser{UserLogged:*userLogged,Admin:userLogged.Admin,MsgResult:"",Users:*users}
	_ = tmpl.Execute(w, appResult)

}
func DeleteUserById(w http.ResponseWriter, r *http.Request,userLogged *model.User, users *[]model.User){

	params := mux.Vars(r)
	id,_:= strconv.Atoi(params["id"])
	var usersRemove [] model.User
	for _, user := range *users {
		if user.ID != int64(id) {

			usersRemove= append(usersRemove,user)
		}
	}
	*users = usersRemove
	tmpl,_ := template.ParseFiles("./view/user.html")
	appResult :=model.AppResultUser{UserLogged:*userLogged,Admin:userLogged.Admin,MsgResult:"",Users:usersRemove}
	_ = tmpl.Execute(w, appResult)
}
func LoginUser(w http.ResponseWriter, r *http.Request,userLogged *model.User, users []model.User){
	userName := r.FormValue("username")
	password := r.FormValue("password")
	link :="./view/login.html"
	for _, user := range users {
		if strings.ToLower( user.UserName) == strings.ToLower(userName) {
			err:=model.ValidatePassword(password,user.PasswordEncryp)
			if VerifyError(w,err){
				return
			}
			*userLogged=user
			link="./view/index.html"
			break
		}
	}

	appResult :=model.AppResult{UserLogged:*userLogged,Admin:userLogged.Admin}
	tmpl,_ := template.ParseFiles(link)
	_ = tmpl.Execute(w, appResult)
}
func LogoutUser(w http.ResponseWriter, r *http.Request,userLogged model.User){
	tmpl,_ := template.ParseFiles("./view/login.html")
	appResult :=model.AppResult{UserLogged:userLogged,Admin:userLogged.Admin}
	_ = tmpl.Execute(w, appResult)
}

