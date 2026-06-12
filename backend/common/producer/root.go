package producer

import (
	"backend/common/tlsconfig"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/IBM/sarama"
	"github.com/google/uuid"
)

type Producer struct {
	syncProducer sarama.SyncProducer
}

func NewProducer(clientIdPrefix string) *Producer {
	syncProducer, err := createProducer(clientIdPrefix)
	if err != nil {
		slog.Error("fail to create producer",
			"err", err,
		)
		panic(err)
	}
	log.Print("success to create kafka producer")
	kp := Producer{syncProducer}
	return &kp
}

func createProducer(clientIdPrefix string) (sarama.SyncProducer, error) {
	cfg := sarama.NewConfig()
	id, err := uuid.NewV7()
	if err != nil {
		slog.Error("fail to create uuid for kafka producer client id")
		return nil, err
	}

	if os.Getenv("PROFILE") == "production" {
		tlsConfig, err1 := tlsconfig.Create("kafka/user/user.crt", "kafka/user/user.key", "kafka/cluster/ca.crt")
		if err1 != nil {
			return nil, err1
		}
		cfg.Net.TLS.Enable = true
		cfg.Net.TLS.Config = tlsConfig
	}

	cfg.ClientID = clientIdPrefix + id.String()
	//cfg.Net.SASL.Enable = true
	//cfg.Net.SASL.Version = 1
	//cfg.Net.SASL.Mechanism = sarama.SASLTypePlaintext
	//cfg.Net.SASL.User = <api-key>
	//cfg.Net.SASL.Password = <secret>

	cfg.Producer.Return.Successes = true
	cfg.Producer.Return.Errors = true
	cfg.Producer.Compression = sarama.CompressionSnappy
	cfg.Producer.RequiredAcks = sarama.WaitForLocal
	cfg.Producer.Idempotent = false
	cfg.Producer.Flush.Messages = 100
	cfg.Producer.Flush.Frequency = time.Millisecond * 5
	cfg.Producer.Retry.Max = 3
	cfg.Producer.Retry.Backoff = time.Millisecond * 250
	cfg.Net.MaxOpenRequests = 5

	return sarama.NewSyncProducer([]string{os.Getenv("KAFKA_ADDRESS")}, cfg)
}

func (p *Producer) PushMessage(topic string, key, value []byte) error {
	msg := sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(value),
	}
	if key != nil {
		msg.Key = sarama.ByteEncoder(key)
	}

	partition, _, err := p.syncProducer.SendMessage(&msg)
	if err != nil {
		log.Print("Failed to produce payload", "err", err)
		return err
	}
	log.Print("Success to produce payload", "partition", partition)
	return nil
}

func (p *Producer) Close() error {
	return p.syncProducer.Close()
}

//func (p *Producer) PushDeadLetter(reason error, originalTopic string, value []byte) error {
//	_ = []sarama.RecordHeader{
//		{Key: []byte("x-error-message"), Value: []byte(reason.Error())},
//		{Key: []byte("x-original-topic"), Value: []byte(originalTopic)},
//		{Key: []byte("x-failed-at"), Value: []byte(time.Now().Format(time.RFC3339))},
//	}
//	err := p.PushMessage("dlq", nil, value)
//	if err != nil {
//		slog.Error("fail to publish dead letter", "err", err)
//		return err
//	}
//	return nil
//}
