package kafka

import (
	"context"
	"encoding/json"
	"log"

	"digit-oss/noc-services/internal/config"
	"digit-oss/noc-services/internal/domain"
	"digit-oss/noc-services/internal/service/notification"

	"github.com/Shopify/sarama"
)

// NOCConsumer translates NOCConsumer.java.
// Listens on the persister topics (save, update, update-workflow) and
// triggers the notification service on each message.
type NOCConsumer struct {
	Cfg                 *config.Config
	NotificationService *notification.NOCNotificationService
}

// Topics returns the list of Kafka topics to subscribe to.
// Matches: @KafkaListener(topics = { "${persister.save.noc.topic}",
//   "${persister.update.noc.topic}", "${persister.update.noc.workflow.topic}" })
func (c *NOCConsumer) Topics() []string {
	return []string{
		c.Cfg.SaveTopic,           // save-noc-application
		c.Cfg.UpdateTopic,         // update-noc-application
		c.Cfg.UpdateWorkflowTopic, // update-noc-workflow
	}
}

// Start launches a sarama ConsumerGroup that listens on the configured topics.
func (c *NOCConsumer) Start(ctx context.Context, brokers []string, groupID string) error {
	cfg := sarama.NewConfig()
	cfg.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	cfg.Consumer.Offsets.Initial = sarama.OffsetOldest
	cfg.Version = sarama.V2_0_0_0

	group, err := sarama.NewConsumerGroup(brokers, groupID, cfg)
	if err != nil {
		return err
	}

	handler := &consumerHandler{consumer: c}

	go func() {
		for {
			if err := group.Consume(ctx, c.Topics(), handler); err != nil {
				log.Printf("[NOCConsumer] error: %v", err)
			}
			if ctx.Err() != nil {
				return
			}
		}
	}()

	log.Printf("[NOCConsumer] listening on topics=%v group=%s", c.Topics(), groupID)
	return nil
}

// ──────────────────────────────────────────────────────────────────────────────
// sarama ConsumerGroupHandler implementation
// ──────────────────────────────────────────────────────────────────────────────

type consumerHandler struct {
	consumer *NOCConsumer
}

func (h *consumerHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (h *consumerHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (h *consumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		log.Printf("[NOCConsumer] received topic=%s partition=%d offset=%d",
			msg.Topic, msg.Partition, msg.Offset)

		var nocReq domain.NocRequest
		if err := json.Unmarshal(msg.Value, &nocReq); err != nil {
			log.Printf("[NOCConsumer] failed to deserialize message on topic %s: %v", msg.Topic, err)
			session.MarkMessage(msg, "")
			continue
		}

		if nocReq.Noc != nil {
			log.Printf("[NOCConsumer] processing applicationNo=%s", nocReq.Noc.ApplicationNo)
		}

		// Trigger notification service (same as Java: notificationService.process(nocRequest))
		h.consumer.NotificationService.Process(&nocReq)

		session.MarkMessage(msg, "")
	}
	return nil
}
