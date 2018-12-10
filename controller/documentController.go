package controller

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"github.com/perseph0ne/DSS-UI-WEB/model"
	"github.com/streadway/amqp"
	"html/template"
	"io/ioutil"
	"net/http"
	"strings"
)

func generateBodySendDocumentRPC(doc model.Document, act model.RequestBase) (bodySend []byte) {

	if act.Action == "create" {
		body := model.RequestCreateDocument{Base: act, Name: doc.Name, Content: doc.File}
		bodySend, _ = json.Marshal(body)
		return
	} else if act.Action == "remove" {
		body := model.RequestDeleteDocument{Base: act, ID: doc.ID}
		bodySend, _ = json.Marshal(body)
		return
	} else if act.Action == "get" {
		body := model.RequestGetDocument{Base: act, ID: doc.ID}
		bodySend, _ = json.Marshal(body)
		return
	} else if act.Action == "list" {
		body := model.RequestListDocument{Base: act}
		bodySend, _ = json.Marshal(body)
		return
	}
	return
}
func documentRPC(bodySend []byte) (docs []model.Document, err error) {
	conn, err := amqp.Dial("amqp://test:Password123@68.183.130.209:15672/")
	if err != nil {
		err = errors.New("Failed to connect to RabbitMQ " + err.Error())
		return
	}
	//failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		err = errors.New("Failed to open a channel " + err.Error())
		return
	}
	//failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"",    // name
		true,  // durable
		false, // delete when usused
		false, // exclusive
		false, // noWait
		nil,   // arguments
	)
	if err != nil {
		err = errors.New("Failed to declare a queue " + err.Error())
		return
	}
	//failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		err = errors.New("Failed to register a consumer " + err.Error())
		return
	}
	//failOnError(err, "Failed to register a consumer")

	corrId := randomString(32)

	err = ch.Publish(
		"",              // exchange
		"yoyo_response", // routing key
		false,           // mandatory
		false,           // immediate
		amqp.Publishing{
			DeliveryMode:  amqp.Persistent,
			ContentType:   "application/json",
			CorrelationId: corrId,
			ReplyTo:       q.Name,
			Body:          bodySend,
		})
	if err != nil {
		err = errors.New("Failed to publish a message " + err.Error())
		return
	}
	//failOnError(err, "Failed to publish a message")

	for d := range msgs {
		if corrId == d.CorrelationId {
			bodyJson, _ := json.Marshal(d.Body)
			var respRPC model.Result
			json.Unmarshal([]byte(bodyJson), &respRPC)
			json.Unmarshal([]byte(respRPC.Json), &docs)
			if respRPC.Code == 0 {
				err = nil
			} else {
				err = errors.New(respRPC.Message)
			}

			break
		}
	}

	return
}
func UploadDocumentServer(w http.ResponseWriter, r *http.Request, userLogged model.User, users []model.User) {
	r.ParseMultipartForm(32 << 20)
	file, handler, err := r.FormFile("file")
	if VerifyError(w, err) == false {
		defer file.Close()
		data, err := ioutil.ReadAll(file)
		if VerifyError(w, err) {
			return
		} else {
			var doc = model.Document{Name: handler.Filename, File: data}
			var act = model.RequestBase{Action: "create"}
			bodySend := generateBodySendDocumentRPC(doc, act)
			_, err := documentRPC(bodySend)
			if VerifyError(w, err) {
				return
			} else {
				message := "The user " + userLogged.UserName + " has uploaded the document " + doc.Name
				email := model.Email{From: userLogged.Email, To: userLogged.GetEmailsToSend(users), Message: message}
				email.ToStr = strings.Join(email.To, ",")
				err := sendEmail(email)
				if VerifyError(w, err) {
					return
				} else {
					var msgResult = "The following email has been sent correctly to all recipients."
					tmpl, _ := template.ParseFiles("./view/email.html")
					_ = tmpl.Execute(w, model.AppResultEmail{UserLogged: userLogged, Admin: false, MsgResult: msgResult, Email: email})
				}
			}
		}
	} else {
		defer file.Close()
		return
	}

}
func DeleteDocumentServer(w http.ResponseWriter, r *http.Request, userLogged model.User, users []model.User) {
	params := mux.Vars(r)
	id := string(params["id"])
	var doc = model.Document{ID: id}
	var act = model.RequestBase{Action: "remove"}
	_, err := documentRPC(generateBodySendDocumentRPC(doc, act))
	if VerifyError(w, err) {
		return
	} else {
		message := "The user " + userLogged.UserName + " has deleted the document with ID" + id
		email := model.Email{From: userLogged.Email, To: userLogged.GetEmailsToSend(users), Message: message}
		email.ToStr = strings.Join(email.To, ",")
		err := sendEmail(email)
		if VerifyError(w, err) {
			return
		} else {
			var msgResult = "The following email has been sent correctly to all recipients."
			tmpl, _ := template.ParseFiles("./view/email.html")
			_ = tmpl.Execute(w, model.AppResultEmail{UserLogged: userLogged, Admin: false, MsgResult: msgResult, Email: email})
		}
	}

}
func GetDocumentsServer(w http.ResponseWriter, r *http.Request, userLogged model.User) {

	var doc = model.Document{}
	var act = model.RequestBase{Action: "list"}
	bodySend := generateBodySendDocumentRPC(doc, act)
	docs, err := documentRPC(bodySend)
	if VerifyError(w, err) {
		return
	} else {
		tmpl, _ := template.ParseFiles("./view/document.html")
		appResult := model.AppResultDocument{UserLogged: userLogged, Admin: userLogged.Admin, MsgResult: "", Docs: docs}
		_ = tmpl.Execute(w, appResult)
	}

}
