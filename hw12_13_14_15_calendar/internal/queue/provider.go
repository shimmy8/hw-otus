package internalqueue

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Provider struct {
	conn       *amqp.Connection
	queue      amqp.Queue
	channel    *amqp.Channel
	connString string
	queueName  string
}

type Message struct {
	EventID      string    `json:"event_id"`       //nolint:tagliatelle
	EventTitle   string    `json:"event_title"`    //nolint:tagliatelle
	EventStartDt time.Time `json:"event_start_dt"` //nolint:tagliatelle
	UserID       string    `json:"event_user_id"`  //nolint:tagliatelle
}

func New(username string, password string, host string, port int, queueName string) *Provider {
	return &Provider{
		connString: fmt.Sprintf(
			"amqp://%s:%s@%s:%d/",
			username,
			password,
			host,
			port,
		),
		queueName: queueName,
	}
}

func (p *Provider) Connect() (err error) {
	p.conn, err = amqp.Dial(p.connString)
	if err != nil {
		return err
	}

	p.channel, err = p.conn.Channel()
	if err != nil {
		return err
	}

	p.queue, err = p.channel.QueueDeclare(
		p.queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	return err
}

func (p *Provider) Close() {
	p.conn.Close()
}

func (p *Provider) Publish(ctx context.Context, msg *Message) error {
	msgJSON, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	return p.channel.PublishWithContext(
		ctx,
		"",
		p.queue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        msgJSON,
		})
}

type onMessageFunc = func(msg *Message)

func (p *Provider) Consume(ctx context.Context, onMessage onMessageFunc) error {
	msgChan, err := p.channel.ConsumeWithContext(
		ctx,
		p.queue.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		case rawMsg, ok := <-msgChan:
			if ok {
				msg := Message{}
				if err := json.Unmarshal(rawMsg.Body, &msg); err != nil {
					return err
				}
				onMessage(&msg)
			}
		}
	}
}
