package rabbit

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dmitriibb/go-common/logging"
	"github.com/dmitriibb/go-common/utils"
	amqp "github.com/rabbitmq/amqp091-go"
	"time"
)

var logger = logging.NewLogger("CommonRabbit")
var uri = GetUri()

func failOnError(err error, msg string) {
	if err != nil {
		logger.Error("%s: %s", msg, err)
	}
}

func TestConnection() {
	conn, err := amqp.Dial(uri)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()
	logger.Debug("test connection to rabbit")
	ch, err := conn.Channel()
	failOnError(err, "failed to get channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body := "Hello World2!"
	err = ch.PublishWithContext(ctx,
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	failOnError(err, "Failed to publish a message")
	logger.Info(" [x] Sent %v\n", body)

}

func SendToQueue(config RabbitQueueConfig, message interface{}) {
	conn, err := amqp.Dial(uri)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "failed to get channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		config.Name,       // name
		config.Durable,    // durable
		config.AutoDelete, // delete when unused
		false,             // exclusive
		false,             // no-wait
		nil,               // arguments
	)
	failOnError(err, "Failed to declare a queue")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	deliveryMode := amqp.Transient
	if config.Persistent {
		deliveryMode = amqp.Persistent
	}

	body, err := json.Marshal(message)
	failOnError(err, fmt.Sprintf("can't marshal '%s' to json", message))

	err = ch.PublishWithContext(ctx,
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			DeliveryMode: deliveryMode,
			ContentType:  "text/json",
			Body:         body,
		})
	failOnError(err, "Failed to publish a message")
	logger.Info("sent: %s", message)
}

func SubscribeToQueue(ctx context.Context, config RabbitQueueConfig, subscription func(delivery amqp.Delivery)) {
	conn, err := amqp.Dial(uri)
	utils.PanicOnError(err, "Failed to connect to RabbitMQ")

	ch, err := conn.Channel()
	utils.PanicOnError(err, "Failed to open a channel")

	q, err := ch.QueueDeclare(
		config.Name,       // name
		config.Durable,    // durable
		config.AutoDelete, // delete when unused
		false,             // exclusive
		false,             // no-wait
		nil,               // arguments
	)
	utils.PanicOnError(err, "Failed to declare a queue")

	messages, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	utils.PanicOnError(err, "Failed to register a consumer")

	go func() {
		<-ctx.Done()
		logger.Debug("stop listening to '%v' rabbit queue", config.Name)
		conn.Close()
		ch.Close()
	}()

	go func() {
		for d := range messages {
			subscription(d)
		}
	}()
}
