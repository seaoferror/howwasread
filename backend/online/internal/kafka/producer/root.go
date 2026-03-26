package producer

import (
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/IBM/sarama"
	"github.com/google/uuid"
)

type KafkaProducer struct {
	producer sarama.AsyncProducer
}

func NewKafkaProducer() *KafkaProducer {
	producer, err := createProducer()
	if err != nil {
		slog.Error("fail to create producer",
			"err", err,
		)
		panic(err)
	}
	log.Print("success to create kafka producer")
	kp := KafkaProducer{producer}

	return &kp
}

func createProducer() (sarama.AsyncProducer, error) {
	cfg := sarama.NewConfig()
	id, err := uuid.NewV7()
	if err != nil {
		slog.Error("fail to create uuid for kafka producer client id")
		return nil, err
	}

	cfg.ClientID = "auth.producer." + id.String()
	//cfg.Net.SASL.Enable = true
	//cfg.Net.SASL.Version = 1
	//cfg.Net.SASL.Mechanism = sarama.SASLTypePlaintext
	//cfg.Net.SASL.User = <api-key>
	//cfg.Net.SASL.Password = <secret>
	//cfg.Net.TLS.Enable = true
	//cfg.Net.SASL.Handshake = true

	cfg.Producer.Return.Successes = true
	cfg.Producer.Return.Errors = true
	cfg.Producer.Compression = sarama.CompressionZSTD
	cfg.Producer.RequiredAcks = sarama.WaitForLocal
	cfg.Producer.Retry.Max = 30
	cfg.Producer.Retry.Backoff = time.Millisecond * 10
	//cfg.Producer.Idempotent = true
	//cfg.Producer.RequiredAcks = sarama.WaitForAll
	//cfg.Net.MaxOpenRequests = 1

	return sarama.NewAsyncProducer([]string{os.Getenv("KAFKA_URL")}, cfg)
}

func (kp *KafkaProducer) PushMessage(topic string, value []byte) error {

	msg := sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(value),
	}

	kp.producer.Input() <- &msg

	select {
	case succeedMsg := <-kp.producer.Successes():
		log.Print("Success to produce message",
			"partition", succeedMsg.Partition)
		return nil
	case err := <-kp.producer.Errors():
		log.Print("Failed to produce message",
			"err", err)
		return err
	}
}

func (kp *KafkaProducer) Close() error {
	return kp.producer.Close()
}
