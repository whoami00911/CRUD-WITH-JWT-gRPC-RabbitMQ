package rabbitmq

import (
	"log"
	"webPractice1/pkg/logger"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/spf13/viper"
	"github.com/whoami00911/gRPC-server/pkg/grpcPb"
	"go.mongodb.org/mongo-driver/bson"
)

type rabbitmqServer struct {
	conn   *amqp.Connection
	ch     *amqp.Channel
	queue  amqp.Queue
	logger *logger.Logger
}

func NewRabbitMQProducer(logger *logger.Logger) *rabbitmqServer {
	conn, err := amqp.Dial(viper.GetString("rabbitmq.url"))
	if err != nil {
		logger.Errorf("amqp Dial error: %s", err)
		log.Fatalf("amqp Dial error: %s", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		logger.Errorf("RabbitMq channel error: %s", err)
		log.Fatalf("RabbitMq channel error: %s", err)
	}

	queue, err := ch.QueueDeclare(
		viper.GetString("rabbitmq.queueName"), //имя очереди
		false,                                 // не durable
		false,                                 // удалять очередь, если не используется
		false,                                 // не эксклюзивная
		false,                                 // не ждать подтверждения
		nil,                                   // дополнительные аргументы
	)
	if err != nil {
		logger.Errorf("QueueDeclare error: %s", err)
		log.Fatalf("QueueDeclare error: %s", err)
	}
	return &rabbitmqServer{
		conn:   conn,
		ch:     ch,
		queue:  queue,
		logger: logger,
	}
}

func (r *rabbitmqServer) Produce(log grpcPb.LogItem) error {
	marshalLog, err := bson.Marshal(log)
	if err != nil {
		return err
	}
	if err := r.ch.Publish(
		"",                                    // используем дефолтный обменник
		viper.GetString("rabbitmq.queueName"), // ключ маршрутизации – имя очереди
		false,                                 // не принудительно (mandatory)
		false,                                 // не немедленно (immediate)
		amqp.Publishing{
			ContentType: "application/bson",
			Body:        marshalLog,
		}); err != nil {
		return err
	}
	return nil
}
