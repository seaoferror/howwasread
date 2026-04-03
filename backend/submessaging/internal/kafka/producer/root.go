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

	cfg.ClientID = "subsignal.producer." + id.String()
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
	cfg.Producer.Retry.Max = 1_000_000_000
	cfg.Producer.Retry.Backoff = time.Millisecond * 10
	//cfg.Producer.Idempotent = true
	//cfg.Producer.RequiredAcks = sarama.WaitForAll
	//cfg.Net.MaxOpenRequests = 1

	return sarama.NewAsyncProducer([]string{os.Getenv("KAFKA_URL")}, cfg)
}

func (kp *KafkaProducer) PushMessage(topic string, header []sarama.RecordHeader, value []byte) error {

	msg := sarama.ProducerMessage{
		Topic:   topic,
		Headers: header,
		Value:   sarama.ByteEncoder(value),
	}

	kp.producer.Input() <- &msg

	select {
	case succeedMsg := <-kp.producer.Successes():
		log.Print("Success to produce payload",
			"partition", succeedMsg.Partition)
		return nil
	case err := <-kp.producer.Errors():
		log.Print("Failed to produce payload",
			"err", err)
		return err
	}
}

func (kp *KafkaProducer) Close() error {
	return kp.producer.Close()
}

func (kp *KafkaProducer) PushDeadLetter(reason error, originalTopic string, value []byte) error {
	headers := []sarama.RecordHeader{
		{Key: []byte("x-error-message"), Value: []byte(reason.Error())},
		{Key: []byte("x-original-topic"), Value: []byte(originalTopic)},
		{Key: []byte("x-failed-at"), Value: []byte(time.Now().Format(time.RFC3339))},
	}
	err := kp.PushMessage("dlq", headers, value)
	if err != nil {
		slog.Error("fail to publish dead letter", "err", err)
		return err
	}
	return nil
}
