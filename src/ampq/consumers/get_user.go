package consumers

import (
	"context"
	"encoding/json"
	"log"
	"main/src/ampq/requests"
	"main/src/ampq/responses"
	"main/src/services"
	"time"

	go_amqp_lib "github.com/lan143/go-amqp-lib"
	"github.com/streadway/amqp"
)

type GetUserConsumer struct {
	service *services.UserService
}

func (consumer *GetUserConsumer) GetQueueName() string {
	return "user.get"
}

func (consumer *GetUserConsumer) IsQuorum() bool {
	return false
}

func (consumer *GetUserConsumer) Handle(
	ctx context.Context,
	channel *amqp.Channel,
	delivery *amqp.Delivery,
	client *go_amqp_lib.Client,
) {
	log.Printf("%s: %s", consumer.GetQueueName(), delivery.Body)
	defer func(delivery *amqp.Delivery, multiple bool) {
		err := delivery.Ack(multiple)
		if err != nil {
			log.Printf("%s: %s", consumer.GetQueueName(), err)
		}
	}(delivery, false)

	request := go_amqp_lib.Request[requests.UserGetByIdRequest]{}
	err := json.Unmarshal(delivery.Body, &request)
	if err != nil {
		consumer.sendResponse(
			client,
			channel,
			delivery.Headers["reply-to"].(string),
			go_amqp_lib.Response[responses.UserResponse]{
				Success: false,
				Message: err.Error(),
			},
		)
		return
	}

	// message is outdated
	if time.Now().After(request.ExpiresAt) {
		log.Printf("%s.recv: message is outdated", consumer.GetQueueName())
		return
	}

	response, err := consumer.service.GetUserById(ctx, &request.Payload)
	if err != nil {
		consumer.sendResponse(
			client,
			channel,
			delivery.Headers["reply-to"].(string),
			go_amqp_lib.Response[responses.UserResponse]{
				Success: false,
				Message: err.Error(),
			},
		)
	} else {
		consumer.sendResponse(
			client,
			channel,
			delivery.Headers["reply-to"].(string),
			go_amqp_lib.Response[responses.UserResponse]{
				Success: true,
				Payload: *response,
			},
		)
	}
}

func (consumer *GetUserConsumer) sendResponse(
	client *go_amqp_lib.Client,
	channel *amqp.Channel,
	route string,
	response go_amqp_lib.Response[responses.UserResponse],
) {
	responseBytes, err := json.Marshal(response)
	if err != nil {
		log.Println(err.Error())
		return
	}

	log.Printf("%s.send: %s", consumer.GetQueueName(), responseBytes)

	err = client.Publish(channel, route, responseBytes, nil)
	if err != nil {
		log.Println(err.Error())
	}
}

func NewGetUserConsumer(service *services.UserService) *GetUserConsumer {
	return &GetUserConsumer{service: service}
}
