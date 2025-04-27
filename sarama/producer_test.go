package sarama

import (
	"github.com/IBM/sarama"
	"github.com/stretchr/testify/assert"
	"testing"
)

var addrs = []string{"localhost:9094"}

// 同步发送
func TestSyncProduce(t *testing.T) {
	cfg := sarama.NewConfig()
	cfg.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer(addrs, cfg)
	assert.NoError(t, err)
	_, _, err = producer.SendMessage(&sarama.ProducerMessage{
		Topic: "test_topic",
		//消息数据
		Value: sarama.StringEncoder("hello world"),
		//在生产者与消费者之间传递
		Headers: []sarama.RecordHeader{
			{
				Key:   []byte("trace_id"),
				Value: []byte("123456"),
			},
		},
		//之作用于发送过程，只与生产者有关
		Metadata: "metadata",
	})
	assert.NoError(t, err)
}

// 异步发送
func TestAsyncProduce(t *testing.T) {
	cfg := sarama.NewConfig()
	cfg.Producer.Return.Errors = true
	cfg.Producer.Return.Successes = true
	producer, err := sarama.NewAsyncProducer(addrs, cfg)
	assert.NoError(t, err)
	msgCh := producer.Input()
	go func() {
		for {
			msg := &sarama.ProducerMessage{
				Topic: "test_topic",
				//消息数据
				Key:   sarama.StringEncoder("123"),
				Value: sarama.StringEncoder("hello world"),
				//在生产者与消费者之间传递
				Headers: []sarama.RecordHeader{
					{
						Key:   []byte("trace_id"),
						Value: []byte("123456"),
					},
				},
				//之作用于发送过程，只与生产者有关
				Metadata: "metadata",
			}
			select {
			case msgCh <- msg:
				//default:
			}
		}
	}()
	errCh := producer.Errors()
	successCh := producer.Successes()

	for {
		select {
		case err := <-errCh:
			t.Log("发送出现错误", err.Err, err.Msg.Topic, err.Msg.Partition, err.Msg.Offset)
		case msg := <-successCh:
			t.Log("发送成功", msg.Topic, msg.Partition, msg.Offset)
		}
	}

}
