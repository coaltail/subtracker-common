package testing

import (
	"context"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	tc "github.com/testcontainers/testcontainers-go"
	tcRabbit "github.com/testcontainers/testcontainers-go/modules/rabbitmq"
	"github.com/testcontainers/testcontainers-go/wait"
)

type cleanupFunction func()
func SetupRabbitMQTestContainer() (*amqp.Connection, cleanupFunction, error) {
	ctx := context.Background()

	container, err := tcRabbit.Run(
		ctx,
		"rabbitmq:3.12-alpine",
		tc.WithWaitStrategy(
			wait.ForLog("Server startup complete").
				WithOccurrence(1).
				WithStartupTimeout(30*time.Second),
		),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("start container: %w", err)
	}

	connStr, err := container.AmqpURL(ctx)
	if err != nil {
		_ = container.Terminate(ctx)
		return nil, nil, fmt.Errorf("get AMQP URL: %w", err)
	}

	conn, err := amqp.Dial(connStr)
	if err != nil {
		_ = container.Terminate(ctx)
		return nil, nil, fmt.Errorf("dial RabbitMQ: %w", err)
	}

	cleanup := func() {
		_ = conn.Close()
		_ = container.Terminate(ctx)
	}

	return conn, cleanup, nil
}
