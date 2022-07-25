package src

import (
	"main/src/ampq/consumers"
	"main/src/config"
	"main/src/repositories"
	"main/src/services"

	go_amqp_lib "github.com/lan143/go-amqp-lib"
	go_healthcheck_lib "github.com/lan143/go-healthcheck-lib"
	"go.uber.org/dig"
)

func BuildContainer() *dig.Container {
	container := dig.New()

	processError(container.Provide(NewApplication))
	processError(container.Provide(config.NewConfig))

	// 3rd party
	processError(container.Provide(go_amqp_lib.NewClient))
	processError(container.Provide(go_amqp_lib.NewConsuming))
	processError(container.Provide(go_healthcheck_lib.NewHealthCheck))

	// Consumers
	processError(container.Provide(consumers.NewGetUserConsumer))

	// Repositories
	processError(container.Provide(repositories.NewUsersRepository))

	// Services
	processError(container.Provide(services.NewUserService))

	return container
}

func processError(err error) {
	if err != nil {
		panic(err)
	}
}
