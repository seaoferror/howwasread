package common

import (
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/IBM/sarama"
	"github.com/google/uuid"
)

type Producer interface {
	PushMessage(topic string, key, value []byte, headers []sarama.RecordHeader)
	Close() error
}

type producer struct {
	asyncProducer sarama.AsyncProducer
}

func NewProducer(clientIdPrefix string) Producer {
	asyncProducer, err := createProducer(clientIdPrefix)
	if err != nil {
		slog.Error("fail to create producer",
			"err", err,
		)
		panic(err)
	}
	log.Print("success to create kafka producer")
	kp := producer{asyncProducer}

	go kp.drainErrorChannel()

	return &kp
}

func createProducer(clientIdPrefix string) (sarama.AsyncProducer, error) {
	cfg := sarama.NewConfig()
	id, err := uuid.NewV7()
	if err != nil {
		slog.Error("fail to create uuid for kafka producer client id")
		return nil, err
	}

	tlsConfig, err1 := CreateTlSConfig(os.Getenv("KAFKA_USER_CERT_PATH"), os.Getenv("KAFKA_USER_KEY_PATH"), os.Getenv("KAFKA_CA_CERT_PATH"))
	if err1 != nil {
		return nil, err1
	}
	cfg.Net.TLS.Config = tlsConfig
	cfg.Net.TLS.Enable = true

	cfg.ClientID = clientIdPrefix + id.String()
	if os.Getenv("KAFKA_API_KEY") != "" {
		cfg.Net.SASL.Enable = true
		cfg.Net.SASL.Version = 1
		cfg.Net.SASL.Mechanism = sarama.SASLTypePlaintext
		cfg.Net.SASL.User = os.Getenv("KAFKA_API_KEY")
		cfg.Net.SASL.Password = os.Getenv("KAFKA_API_SECRET")
		cfg.Net.SASL.Handshake = true
	}

	cfg.Producer.Return.Successes = false
	cfg.Producer.Return.Errors = true
	cfg.Producer.Compression = sarama.CompressionSnappy
	cfg.Producer.RequiredAcks = sarama.WaitForAll
	cfg.Producer.Idempotent = false
	cfg.Producer.Flush.Messages = 100
	cfg.Producer.Flush.Frequency = time.Millisecond * 5
	cfg.Producer.Retry.Max = 5
	cfg.Producer.Retry.Backoff = time.Millisecond * 300
	cfg.Net.MaxOpenRequests = 5

	return sarama.NewAsyncProducer([]string{os.Getenv("KAFKA_ADDRESS")}, cfg)
}

func (p *producer) PushMessage(topic string, key, value []byte, headers []sarama.RecordHeader) {
	msg := sarama.ProducerMessage{
		Topic:   topic,
		Headers: headers,
		Value:   sarama.ByteEncoder(value),
	}
	if key != nil {
		msg.Key = sarama.ByteEncoder(key)
	}

	p.asyncProducer.Input() <- &msg
}

func (p *producer) Close() error {
	return p.asyncProducer.Close()
}

func (p *producer) drainErrorChannel() {
	for err := range p.asyncProducer.Errors() {
		log.Print(
			"Failed to produce payload",
			"err", err.Err,
			"topic", err.Msg.Topic,
		)
		//TODO: dlq? or send notification to developer directly
	}
}
