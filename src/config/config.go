package config

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	go_amqp_lib "github.com/lan143/go-amqp-lib"
)

type Config struct {
	amqpConfig            go_amqp_lib.AmqpConfig
	healthCheckAddress    string
	isDebug               bool
	consumingMaxJobsCount uint64
}

func (c *Config) Init() error {
	log.Println("config: init")

	_ = godotenv.Load()

	err := c.loadAmqpConfig()
	if err != nil {
		return err
	}

	err = c.loadHealthCheckAddress()
	if err != nil {
		return err
	}

	err = c.loadIsDebug()
	if err != nil {
		return err
	}

	err = c.loadConsumingMaxJobsCount()
	if err != nil {
		return err
	}

	log.Println("config.init: successful")

	return nil
}

func (c *Config) GetAmqpConfig() go_amqp_lib.AmqpConfig {
	return c.amqpConfig
}

func (c *Config) GetHealthCheckAddress() string {
	return c.healthCheckAddress
}

func (c *Config) IsDebug() bool {
	return c.isDebug
}

func (c *Config) GetConsumingMaxJobsCount() uint64 {
	return c.consumingMaxJobsCount
}

func (c *Config) loadAmqpConfig() error {
	c.amqpConfig = go_amqp_lib.AmqpConfig{
		Host:     os.Getenv("AMQP_HOST"),
		Port:     os.Getenv("AMQP_PORT"),
		Username: os.Getenv("AMQP_USER"),
		Password: os.Getenv("AMQP_PASSWORD"),
		VHost:    os.Getenv("AMQP_VHOST"),
	}

	return nil
}

func (c *Config) loadHealthCheckAddress() error {
	c.healthCheckAddress = os.Getenv("HEALTH_CHECK_ADDR")
	if len(c.healthCheckAddress) == 0 {
		c.healthCheckAddress = "0.0.0.0:8080"
	}

	return nil
}

func (c *Config) loadIsDebug() error {
	isDebugVal := os.Getenv("DEBUG")

	if strings.Compare(isDebugVal, "true") == 0 || strings.Compare(isDebugVal, "1") == 0 {
		c.isDebug = true
	} else {
		c.isDebug = false
	}

	return nil
}

func (c *Config) loadConsumingMaxJobsCount() error {
	var err error
	c.consumingMaxJobsCount, err = strconv.ParseUint(os.Getenv("CONSUMING_MAX_JOBS_COUNT"), 10, 64)
	if err != nil {
		c.consumingMaxJobsCount = 10
	}

	return nil
}

func NewConfig() *Config {
	return &Config{}
}
