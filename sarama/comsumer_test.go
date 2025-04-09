package sarama

import (
	"context"
	"encoding/json"
	"github.com/IBM/sarama"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestConsumer(t *testing.T) {
	cfg := sarama.NewConfig()
	consumer, err := sarama.NewConsumerGroup(addrs,
		"test_group", cfg)
	require.NoError(t, err)
	err = consumer.Consume(context.Background(),
		[]string{"test_topic"}, testConsumerGroupHandler{})
	t.Log(err)
}

type testConsumerGroupHandler struct {
	// 这个是 sarama 的 ConsumerGroupHandler
}

func (t testConsumerGroupHandler) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (t testConsumerGroupHandler) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

type MyBizMsg struct {
	Name string
}

func (t testConsumerGroupHandler) ConsumeClaim(
	session sarama.ConsumerGroupSession,
	claim sarama.ConsumerGroupClaim) error {
	msg := claim.Messages()
	for msg := range msg {
		var bizMsg MyBizMsg
		err := json.Unmarshal(msg.Value, &bizMsg)
		if err != nil {
			continue
		}
		// 处理消息
		session.MarkMessage(msg, "")
	}
	return nil
}
