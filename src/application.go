package src

import (
	"context"
	"log"

	"main/src/ampq/consumers"
	"main/src/config"
	"os"
	"os/signal"
	"sync"
	"syscall"

	go_amqp_lib "github.com/lan143/go-amqp-lib"
	go_healthcheck_lib "github.com/lan143/go-healthcheck-lib"
)

type Application struct {
	amqpClient  *go_amqp_lib.Client
	consuming   *go_amqp_lib.Consuming
	config      *config.Config
	healthCheck *go_healthcheck_lib.HealthCheck

	getUserConsumer *consumers.GetUserConsumer

	wg                    sync.WaitGroup
	consumingShutdownChan chan<- interface{}
	sigs                  chan os.Signal
}

func (a *Application) Init(ctx context.Context) error {
	log.Printf("application: init")

	a.sigs = make(chan os.Signal, 1)
	a.consumingShutdownChan = make(chan interface{}, 1)

	err := a.config.Init()
	if err != nil {
		return err
	}

	signal.Notify(a.sigs, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	a.amqpClient.Init(a.config.GetAmqpConfig(), 10, &a.wg)

	a.consuming.Init(a.config.GetConsumingMaxJobsCount(), &a.wg, a.consumingShutdownChan)
	a.consuming.AddConsumer(a.getUserConsumer)

	a.healthCheck.Init(a.config.GetHealthCheckAddress(), &a.wg)
	a.healthCheck.AddReadinessProbe(a.amqpClient)

	return nil
}

func (a *Application) Run(ctx context.Context) error {
	log.Printf("application.run: start")

	cancelCtx, cancelFunc := context.WithCancel(ctx)
	go a.processSignals(cancelFunc)

	a.healthCheck.Run(cancelCtx)

	err := a.amqpClient.Run(cancelCtx)
	if err != nil {
		return err
	}

	a.consuming.Run(cancelCtx)

	log.Println("application.run: running")

	a.wg.Wait()

	log.Println("application: graceful shutdown.")

	return nil
}

func (a *Application) processSignals(cancelFun context.CancelFunc) {
	select {
	case <-a.sigs:
		log.Println("application: received shutdown signal from OS")
		cancelFun()
		break
	}
}

func NewApplication(amqpClient *go_amqp_lib.Client,
	consuming *go_amqp_lib.Consuming,
	config *config.Config,
	healthCheck *go_healthcheck_lib.HealthCheck,
	// Consumers
	getUserConsumer *consumers.GetUserConsumer,

) *Application {
	return &Application{
		amqpClient:  amqpClient,
		consuming:   consuming,
		config:      config,
		healthCheck: healthCheck,
		// Consumers
		getUserConsumer: getUserConsumer,
	}
}
