package consumer

import (
	"backend/messagepreprocess/internal/service"
	"backend/payload"
	"context"
	"encoding/json"
	"errors"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/IBM/sarama"
	"github.com/google/uuid"

	_ "github.com/joho/godotenv/autoload"
)

type Consumer struct {
	consumerGroup sarama.ConsumerGroup
	service       *service.Service
}

func NewConsumer(s *service.Service) *Consumer {
	consumerGroup, err := connectConsumer("preprocess_message")
	if err != nil {
		log.Panicf("fail to create consumer group client: %v", err)
	}
	return &Consumer{
		consumerGroup: consumerGroup,
		service:       s,
	}
}

func connectConsumer(groupID string) (sarama.ConsumerGroup, error) {
	cfg := sarama.NewConfig()
	id, err := uuid.NewV7()
	if err != nil {
		slog.Error("fail to create uuid for kafka client uuid")
		return nil, err
	}
	cfg.ClientID = "preprocess_message." + id.String()
	//cfg.Net.SASL.Enable = true
	//cfg.Net.SASL.Version = 1
	//cfg.Net.SASL.Mechanism = sarama.SASLTypePlaintext
	//cfg.Net.SASL.User = <api-key>
	//cfg.Net.SASL.Password = <secret>
	//cfg.Net.TLS.Enable = true
	//cfg.Net.SASL.Handshake = true

	cfg.Consumer.Return.Errors = true
	cfg.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategySticky()}
	//if balance strategy need to be change flexible, use switch-case with config di
	cfg.Consumer.Offsets.Initial = sarama.OffsetOldest
	//this setting make possible to consume payload which is stored but not consumed for certain reason like worker internal down

	return sarama.NewConsumerGroup([]string{os.Getenv("KAFKA_URL")}, groupID, cfg)
}

func (c *Consumer) Setup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (c *Consumer) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (c *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case msg := <-claim.Messages():
			err := c.distinguishMessage(session.Context(), msg)
			if err != nil {
				log.Printf("Fail to manage message: %v", err)
				return err
			}
			session.MarkMessage(msg, "")
			continue
		case <-session.Context().Done():
			return nil
		}
	}
}

func (c *Consumer) GetMessage(topics []string) error {
	ctx, cancel := context.WithCancel(context.Background())

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			if err := c.consumerGroup.Consume(ctx, topics, c); err != nil {
				return
			}
			if ctx.Err() != nil {
				return
			}
		}
	}()

	log.Println("Sarama consumer up and running")

	sigusr1 := make(chan os.Signal, 1)
	signal.Notify(sigusr1, syscall.SIGUSR1)

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

	keepRunning := true
	consumptionIsPaused := false
	for keepRunning {
		select {
		case <-ctx.Done():
			log.Println("terminating: context cancelled")
			keepRunning = false
		case <-sigterm:
			log.Println("terminating: via signal")
			keepRunning = false
		case <-sigusr1:
			toggleConsumptionFlow(c.consumerGroup, &consumptionIsPaused)
		}
	}
	cancel()
	wg.Wait()

	if err := c.consumerGroup.Close(); err != nil {
		log.Printf("Error closing client: %v", err)
		return err
	}
	return nil
}

func toggleConsumptionFlow(client sarama.ConsumerGroup, isPaused *bool) {
	if *isPaused {
		client.ResumeAll()
		log.Println("Resuming consumption")
	} else {
		client.PauseAll()
		log.Println("Pausing consumption")
	}

	*isPaused = !*isPaused
}

func (c *Consumer) distinguishMessage(ctx context.Context, message *sarama.ConsumerMessage) error {
	if message.Topic == "chat.message" {
		var p payload.ChatMessage
		err := json.Unmarshal(message.Value, &p)
		if err != nil {
			slog.Error("fail to unmarshal payload value",
				"err", err,
				"payload.Value", message.Value)
			return err
		}
		err = c.service.ManageMessage(ctx, uuid.UUID(p.Id), uuid.UUID(p.FromId), p.ToIdType, uuid.UUID(p.ToId), p.ContentType, p.Contents)
		if err != nil {
			return err
		}
		return nil
	}
	return errors.New("this topic does not exist")
}
