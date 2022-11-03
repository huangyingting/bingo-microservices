package rabbitmq

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"sync"
	"time"

	"github.com/go-kratos/kratos/v2/log"

	"bingo/pkg/retry"

	amqp "github.com/rabbitmq/amqp091-go"
)

type ConsumeFunc func(message amqp.Delivery) error

type Subscriber struct {
	c       *SubscriberConfig
	wg      sync.WaitGroup
	h       *log.Helper
	conn    *amqp.Connection
	consume ConsumeFunc
	ctx     context.Context
	cancel  context.CancelFunc
}

func NewSubscriber(
	c *SubscriberConfig,
	h *log.Helper,
	consume ConsumeFunc,
) (*Subscriber, error) {

	_, err := amqp.ParseURI(c.AmqpUri)
	if err != nil {
		return nil, err
	}

	if c.RetryInterval.AsDuration() == 0 {
		return nil, errors.New("retry interval is zero")
	}

	if c.WorkerCount == 0 {
		return nil, errors.New("worker count is zero")
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &Subscriber{
		c:       c,
		h:       h,
		wg:      sync.WaitGroup{},
		conn:    nil,
		ctx:     ctx,
		cancel:  cancel,
		consume: consume,
	}, nil
}

func (s *Subscriber) Start(ctx context.Context) error {
	if err := s.start(false); err != nil {
		s.stop()
		return err
	}
	return nil
}

func (s *Subscriber) Stop(ctx context.Context) error {
	return s.stop()
}

func (s *Subscriber) start(reconnect bool) error {
	attempts := s.c.ConnectAttempt
	if reconnect {
		attempts = s.c.ReconnectAttempt
	}
	delayType := retry.FixedDelay
	if reconnect {
		delayType = retry.BackOffDelay
	}

	op := "connect"
	if reconnect {
		op = "reconnect"
	}

	var i uint = 0
	err := retry.Do(
		func() error {
			i++
			if err := s.connect(); err != nil {
				return err
			}
			return s.setup()
		},
		retry.Context(s.ctx),
		retry.Attempts(attempts),
		retry.DelayType(delayType),
		retry.Delay(s.c.RetryInterval.AsDuration()),
		retry.OnRetry(func(n uint32, err error) {
			s.h.Infof("%s attempt: %d/%d",
				op,
				i,
				attempts)
		}),
	)

	if err != nil {
		return err
	}

	return s.runWorkers()
}

func (s *Subscriber) stop() error {
	s.cancel()
	s.wg.Wait()
	if s.conn != nil {
		if err := s.conn.Close(); err != nil {
			s.h.Errorf("stop: %v", err)
			return err
		}
	}
	return nil
}

func (s *Subscriber) setup() error {

	s.h.Debug("setup - open channel")
	chn, err := s.conn.Channel()
	if err != nil {
		s.h.Errorf("setup - open channel: %v", err)
		return err
	}
	defer chn.Close()

	s.h.Debug("setup - exchange declare")
	if err := chn.ExchangeDeclare(
		s.c.ExchangeName,
		s.c.ExchangeType,
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		s.h.Errorf("setup - exchange declare: %v", err)
		return err
	}

	args := amqp.Table{"x-queue-mode": s.c.QueueMode}
	if s.c.DlExchangeName != "" {
		args["x-dead-letter-exchange"] = s.c.DlExchangeName
	}

	s.h.Debug("setup - queue declare")
	if _, err := chn.QueueDeclare(
		s.c.QueueName,
		true,
		false,
		false,
		false,
		args,
	); err != nil {
		s.h.Errorf("setup - queue declare: %v", err)
		return err
	}

	s.h.Debugf("setup - queue bind")
	if err := chn.QueueBind(
		s.c.QueueName,
		s.c.RoutingKey,
		s.c.ExchangeName,
		false,
		nil,
	); err != nil {
		s.h.Errorf("setup - queue bind: %v", err)
		return err
	}

	return nil
}

func (s *Subscriber) runWorkers() error {
	var startChs []chan struct{}
	var i uint32 = 0
	for ; i < s.c.WorkerCount; i++ {
		startCh := make(chan struct{})
		if err := s.runWorker(i+1, startCh); err != nil {
			return err
		} else {
			startChs = append(startChs, startCh)
		}
	}
	for _, ch := range startChs {
		ch <- struct{}{}
	}
	return nil
}

func (s *Subscriber) runWorker(
	workerId uint32, startCh chan struct{},
) error {

	chn, err := s.conn.Channel()
	if err != nil {
		s.h.Errorf("runWorker - open channel: %v", err)
		return err
	}

	if err := chn.Qos(int(s.c.PrefetchCount), 0, false); err != nil {
		s.h.Errorf(
			"runWorker - qos channel: %v", err)
		return err
	}

	messages, err := chn.Consume(
		s.c.QueueName,
		fmt.Sprintf(
			"%s(%d/%d)",
			s.c.Name,
			workerId,
			s.c.WorkerCount,
		),
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		s.h.Errorf(
			"runWorker - consume channel: %s(%d/%d): %v",
			s.c.Name,
			workerId,
			s.c.WorkerCount,
			err)
		return err
	}

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()

		select {
		case <-startCh:
		case <-s.ctx.Done():
			return
		}

		s.onChannelClose(chn, workerId, startCh)
		s.h.Debugf("runWorker - started: %s(%d/%d)",
			s.c.Name,
			workerId,
			s.c.WorkerCount)
		for {
			select {
			case message, ok := <-messages:
				if !ok {
					defer s.h.Debugf("runWorker - stopped: %s(%d/%d)",
						s.c.Name,
						workerId,
						s.c.WorkerCount)
					return
				}
				// Requeue if error reported at first time, otherwise send to dead letter queue
				if err := s.consume(message); err != nil {
					if message.Redelivered {
						s.h.Debugf("runWorker - consume: reject redelivered")
						message.Reject(false)
					} else {
						s.h.Debugf("runWorker - consume: reject")
						message.Reject(true)
					}
				} else {
					if err := message.Ack(false); err != nil {
						s.h.Errorf("runWorker - ack error: %v", err)
					}
				}
			case <-s.ctx.Done():
				s.h.Debugf("runWorker - canceled")
				return
			}
		}
	}()
	return nil
}

func (s *Subscriber) dail(network, addr string) (net.Conn, error) {
	conn, err := net.DialTimeout(network, addr, s.c.ConnectTimeout.AsDuration())
	if err != nil {
		return nil, err
	}
	if err := conn.SetDeadline(time.Now().Add(s.c.ConnectTimeout.AsDuration())); err != nil {
		return nil, err
	}
	return conn, nil
}

func (s *Subscriber) connect() error {
	var err error
	uri, err := amqp.ParseURI(s.c.AmqpUri)
	if err != nil {
		return err
	}

	var tlsCnf *tls.Config = nil

	if uri.Scheme == "amqps" {
		tlsCnf := &tls.Config{}
		if s.c.CaCert != "" {
			tlsCnf.RootCAs = x509.NewCertPool()
			if ca, err := ioutil.ReadFile(s.c.CaCert); err == nil {
				tlsCnf.RootCAs.AppendCertsFromPEM(ca)
			}
		}
		if cert, err := tls.LoadX509KeyPair(s.c.ClientCert, s.c.ClientKey); err == nil {
			tlsCnf.Certificates = append(tlsCnf.Certificates, cert)
		}
	}

	s.conn, err = amqp.DialConfig(s.c.AmqpUri,
		amqp.Config{
			Properties:      amqp.Table{"connection_name": s.c.Name},
			TLSClientConfig: tlsCnf,
			Dial:            s.dail,
		},
	)
	if err != nil {
		return err
	}
	s.onClose()
	return nil
}

func (s *Subscriber) onClose() {
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		select {
		case connClose := <-s.conn.NotifyClose(make(chan *amqp.Error)):
			if connClose != nil {
				s.h.Debugf("onClose - closed: %v", connClose)
				err := s.start(true)
				if err == nil {
					s.h.Debug("onClose - reconnected")
				} else {
					s.h.Error("onClose - reconnect failed")
				}
			} else {
				s.h.Debug("onClose - closed explicitly")
			}
		case <-s.ctx.Done():
			s.h.Debug("onClose - context canceled")
		}
	}()
}

func (s *Subscriber) onChannelClose(chn *amqp.Channel, workerId uint32, startCh chan struct{}) {
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		select {
		case amqpError := <-chn.NotifyClose(make(chan *amqp.Error)):
			s.h.Debugf("onChannelClose - closed: %v", amqpError)
			if amqpError != nil && amqpError.Code == amqp.ConnectionForced {
				s.h.Debug("onChannelClose - won't resume worker")
				return
			}
			s.h.Debug("onChannelClose - resume worker")
			if err := s.runWorker(workerId, startCh); err != nil {
				s.h.Errorf("onChannelClose - resume worker: %v", err)
			} else {
				startCh <- struct{}{}
			}
		case <-s.ctx.Done():
			s.h.Debug("onChannelClose - context canceled")
		}
	}()
}
