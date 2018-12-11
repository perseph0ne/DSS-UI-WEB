package controller

import (
	"encoding/json"
	"errors"
	"github.com/perseph0ne/DSS-UI-WEB/model"
	"github.com/streadway/amqp"
)

func generateBodySendEmailRPC(email model.Email, act model.RequestBase) (bodySend []byte) {

	if act.Action == "send" {
		body := model.RequestSendMail{Base: act, From: email.From, To: email.To, Message: email.Message}
		bodySend, _ = json.Marshal(body)
	}
	return
}
func emailRPC(bodySend []byte) (err error) {
	conn, err := amqp.Dial("amqp://test:Password123@68.183.130.209:5672/")
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

	_, err = ch.QueueDeclare(
		"yoyo_request_mail", // name
		true,                // durable
		false,               // delete when usused
		false,               // exclusive
		false,               // noWait
		nil,                 // arguments
	)
	if err != nil {
		err = errors.New("Failed to declare a queue " + err.Error())
		return
	}
	//failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		"yoyo_response_mail", // queue
		"",                   // consumer
		true,                 // auto-ack
		false,                // exclusive
		false,                // no-local
		false,                // no-wait
		nil,                  // args
	)
	if err != nil {
		err = errors.New("Failed to register a consumer " + err.Error())
		return
	}
	//failOnError(err, "Failed to register a consumer")

	corrId := randomString(32)

	err = ch.Publish(
		"",                  // exchange
		"yoyo_request_mail", // routing key
		false,               // mandatory
		false,               // immediate
		amqp.Publishing{
			ContentType:   "application/json",
			CorrelationId: corrId,
			ReplyTo:       "yoyo_response_mail",
			Body:          bodySend,
		})
	if err != nil {
		err = errors.New("Failed to publish a message " + err.Error())
		return
	}
	//failOnError(err, "Failed to publish a message")

	for d := range msgs {
		if corrId == d.CorrelationId {
			//bodyJson, _ := json.Marshal(d.Body)
			var respRPC model.Result
			json.Unmarshal(d.Body, &respRPC)
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
func sendEmail(email model.Email) (err error) {
	bodySend := generateBodySendEmailRPC(email, model.RequestBase{Action: "send"})
	err = emailRPC(bodySend)
	return
}
