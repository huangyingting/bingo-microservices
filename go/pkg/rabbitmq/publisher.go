package rabbitmq

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io/ioutil"
	"net"
	"os"
	"sync"
	"syscall"
	"time"

	"bingo/pkg/retry"

	"github.com/go-kratos/kratos/v2/log"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Publisher struct {
	c       *PublisherConfig
	h       *log.Helper
	mutex   sync.Mutex
	conn    *amqp.Connection
	channel *amqp.Channel
	ctx     context.Context
	cancel  context.CancelFunc
}

func NewPublisher(
	c *PublisherConfig,
	h *log.Helper,
) (*Publisher, error) {

	_, err := amqp.ParseURI(c.AmqpUri)
	if err != nil {
		return nil, err
	}

	if c.RetryInterval.AsDuration() == 0 {
		return nil, errors.New("retry interval is zero")
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &Publisher{
		c:      c,
		h:      h,
		conn:   nil,
		ctx:    ctx,
		cancel: cancel,
	}, nil
}

func (p *Publisher) Publish(message []byte) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	var i uint = 0
	return retry.Do(
		func() error {
			var err error
			i++
			// Reuse exiting channel
			if p.channel != nil {
				return p.publish(message)
			}
			if err = p.connect(); err != nil {
				p.h.Error(err)
				return err
			}
			if err = p.setup(); err != nil {
				p.h.Error(err)
				return err
			}
			return p.publish(message)
		},
		retry.Attempts(p.c.RetryAttempt),
		retry.Delay(p.c.RetryInterval.AsDuration()),
		retry.OnRetry(func(n uint32, err error) {
			p.h.Infof("publish attempt: %d/%d",
				i,
				p.c.RetryAttempt)
		}),
		retry.RetryIf(p.shouldRetry),
	)
}

func (p *Publisher) dail(network, addr string) (net.Conn, error) {
	conn, err := net.DialTimeout(network, addr, p.c.ConnectTimeout.AsDuration())
	if err != nil {
		return nil, err
	}
	if err := conn.SetDeadline(time.Now().Add(p.c.ConnectTimeout.AsDuration())); err != nil {
		return nil, err
	}
	return conn, nil
}

func (p *Publisher) publish(message []byte) error {
	args := amqp.Table{"x-queue-mode": p.c.QueueMode}
	if err := p.channel.Publish(
		p.c.ExchangeName,
		p.c.RoutingKey,
		false,
		false,
		amqp.Publishing{
			Headers:         args,
			ContentType:     "text/plain",
			ContentEncoding: "",
			Body:            message,
			DeliveryMode:    amqp.Transient, // 1=non-persistent, 2=persistent
			Priority:        0,              // 0-9)
		}); err != nil {
		// Close channel and connection
		p.channel.Close()
		p.channel = nil
		p.conn.Close()
		p.conn = nil
		return err
	}
	return nil
}

func (p *Publisher) shouldRetry(err error) bool {
	// Retry always on ampq error
	_, ok := err.(*amqp.Error)
	if ok {
		return true
	}

	netErr, ok := err.(net.Error)
	if !ok {
		return false
	}
	if netErr.Timeout() {
		return true
	}
	opErr, ok := netErr.(*net.OpError)
	if !ok {
		return false
	}
	switch t := opErr.Err.(type) {
	case *net.DNSError:
		return true
	case *os.SyscallError:
		if errno, ok := t.Err.(syscall.Errno); ok {
			switch errno {
			case syscall.ECONNREFUSED:
				return true
			case syscall.ETIMEDOUT:
				return true
			}
		}
	}
	return false
}

func (p *Publisher) setup() error {
	chn, err := p.conn.Channel()
	if err != nil {
		p.h.Errorf("setup - open channel: %v", err)
		return err
	}

	p.h.Debug("setup - exchange declare")

	if err := chn.ExchangeDeclare(
		p.c.ExchangeName,
		p.c.ExchangeType,
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		p.h.Errorf("setup - exchange declare: %v", err)
		return err
	}

	p.channel = chn
	return nil
}

func (p *Publisher) connect() error {
	var err error
	uri, err := amqp.ParseURI(p.c.AmqpUri)
	if err != nil {
		return err
	}

	var tlsCnf *tls.Config = nil
	if uri.Scheme == "amqps" {
		tlsCnf := &tls.Config{}
		if p.c.CaCert != "" {
			tlsCnf.RootCAs = x509.NewCertPool()
			if ca, err := ioutil.ReadFile(p.c.CaCert); err == nil {
				tlsCnf.RootCAs.AppendCertsFromPEM(ca)
			}
		}
		if cert, err := tls.LoadX509KeyPair(p.c.ClientCert, p.c.ClientKey); err == nil {
			tlsCnf.Certificates = append(tlsCnf.Certificates, cert)
		}
	}

	p.conn, err = amqp.DialConfig(p.c.AmqpUri,
		amqp.Config{
			Properties:      amqp.Table{"connection_name": p.c.Name},
			TLSClientConfig: tlsCnf,
			Dial:            p.dail,
		},
	)
	if err != nil {
		return err
	}

	return nil
}
