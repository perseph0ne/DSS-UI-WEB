package controller

import (
	"html/template"
	"math/rand"
	"net/http"
)


func ShowError(w http.ResponseWriter, errorMessage error)  {
	tmpl,_ := template.ParseFiles("./view/error.html")
	type MsgError struct {
		MessageError string
	}
	var MessageError = MsgError{errorMessage.Error()}
	_ = tmpl.Execute(w, MessageError)

}
func VerifyError(w http.ResponseWriter,  errorMessage error) bool{
	if errorMessage!=nil{
		ShowError(w,errorMessage)
		return true
	}
	return false
}
func randomString(l int) string {
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		bytes[i] = byte(randInt(65, 90))
	}
	return string(bytes)
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}