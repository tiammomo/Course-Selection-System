package mq

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"course_select/internal/config"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
)

// Client RocketMQ 客户端
type Client struct {
	producer   rocketmq.Producer
	consumer   rocketmq.PushConsumer
	cfg        *config.RocketMQConfig
	isProducer bool
	isConsumer bool
}

// BookingMessage 选课消息
type BookingMessage struct {
	StudentID string    `json:"student_id"`
	CourseID  string    `json:"course_id"`
	Timestamp time.Time `json:"timestamp"`
}

// New 创建 RocketMQ 客户端
func New(cfg *config.RocketMQConfig) (*Client, error) {
	client := &Client{
		cfg: cfg,
	}

	// 创建生产者
	p, err := rocketmq.NewProducer(
		producer.WithNameServer([]string{cfg.NameServer}),
		producer.WithGroupName(cfg.GroupID),
		producer.WithInstanceName(cfg.InstanceName),
		producer.WithRetry(2),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create producer: %w", err)
	}
	client.producer = p
	client.isProducer = true

	// 启动生产者
	if err := client.producer.Start(); err != nil {
		return nil, fmt.Errorf("failed to start producer: %w", err)
	}

	return client, nil
}

// NewConsumer 创建 RocketMQ 消费者
func NewConsumer(cfg *config.RocketMQConfig, handler func(*BookingMessage) error) (*Client, error) {
	client := &Client{
		cfg: cfg,
	}

	// 创建消费者
	c, err := rocketmq.NewPushConsumer(
		consumer.WithNameServer([]string{cfg.NameServer}),
		consumer.WithGroupName(cfg.GroupID),
		consumer.WithConsumerModel(consumer.Clustering),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer: %w", err)
	}
	client.consumer = c
	client.isConsumer = true

	// 注册消息处理函数
	err = c.Subscribe(cfg.Topic, consumer.MessageSelector{}, func(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
		for _, msg := range msgs {
			var bookingMsg BookingMessage
			if err := json.Unmarshal(msg.Body, &bookingMsg); err != nil {
				fmt.Printf("failed to unmarshal message: %v\n", err)
				return consumer.ConsumeRetryLater, err
			}

			if err := handler(&bookingMsg); err != nil {
				fmt.Printf("failed to handle message: %v\n", err)
				return consumer.ConsumeRetryLater, err
			}
		}
		return consumer.ConsumeSuccess, nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to subscribe: %w", err)
	}

	// 启动消费者
	if err := c.Start(); err != nil {
		return nil, fmt.Errorf("failed to start consumer: %w", err)
	}

	return client, nil
}

// SendBookingMessage 发送选课消息
func (c *Client) SendBookingMessage(ctx context.Context, msg *BookingMessage) error {
	if !c.isProducer {
		return fmt.Errorf("producer not initialized")
	}

	body, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	mqMsg := &primitive.Message{
		Topic: c.cfg.Topic,
		Body:  body,
	}

	_, err = c.producer.SendSync(ctx, mqMsg)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}

// Close 关闭连接
func (c *Client) Close() error {
	if c.isProducer && c.producer != nil {
		c.producer.Shutdown()
	}
	if c.isConsumer && c.consumer != nil {
		c.consumer.Shutdown()
	}
	return nil
}
