package server

import (
	bgv1 "bingo/api/bg/v1"
	v1 "bingo/api/bi/v1"
	"bingo/app/bi/internal/conf"
	"bingo/app/bi/internal/data"
	"bingo/app/bi/internal/service"
	"bingo/pkg/rabbitmq"
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	transgrpc "github.com/go-kratos/kratos/v2/transport/grpc"
	transhttp "github.com/go-kratos/kratos/v2/transport/http"
	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/protobuf/proto"
)

func location(g *conf.Geo, clickEvent *v1.ClickEvent) (rsp *bgv1.LocationReply, err error) {
	if g.Protocol == conf.Geo_HTTP {
		conn, err := transhttp.NewClient(
			context.Background(),
			transhttp.WithMiddleware(
				recovery.Recovery(),
			),
			transhttp.WithEndpoint(g.HttpAddr),
		)

		if err != nil {
			return nil, err
		}

		defer conn.Close()
		client := bgv1.NewGeoHTTPClient(conn)
		return client.Location(context.Background(), &bgv1.LocationRequest{Ip: clickEvent.Ip})
	}

	conn, err := transgrpc.DialInsecure(
		context.Background(),
		transgrpc.WithEndpoint(g.HttpAddr),
		transgrpc.WithMiddleware(
			recovery.Recovery(),
		),
	)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := bgv1.NewGeoClient(conn)
	return client.Location(context.Background(), &bgv1.LocationRequest{Ip: clickEvent.Ip})
}

func NewRabbitmqServer(
	c *rabbitmq.SubscriberConfig,
	g *conf.Geo,
	bi *service.BIService,
	logger log.Logger,
) *rabbitmq.Subscriber {

	h := log.NewHelper(log.NewFilter(logger, log.FilterLevel(log.LevelDebug)))
	consume := func(event amqp.Delivery) error {
		h.Debugf(
			"consume - [%s][%v][messge id: %s]",
			event.ConsumerTag,
			event.DeliveryTag,
			event.MessageId,
		)
		var clickEvent = v1.ClickEvent{}
		if err := proto.Unmarshal(event.Body, &clickEvent); err != nil {
			h.Errorf("consume - protobuf unmarshal failed: %v", err)
			return err
		}

		r, err := location(g, &clickEvent)
		if err != nil {
			h.Warnf("consume - get location failed: %v", err)
			return bi.Create(context.TODO(), &data.Click{
				Alias:     clickEvent.Alias,
				IP:        clickEvent.Ip,
				UA:        clickEvent.Ua,
				Referer:   clickEvent.Referer,
				CreatedAt: clickEvent.CreatedAt.AsTime(),
			})
		}

		h.Debugf("consume - get location: %v", r)

		err = bi.Create(context.TODO(), &data.Click{
			Alias:     clickEvent.Alias,
			IP:        clickEvent.Ip,
			UA:        clickEvent.Ua,
			Referer:   clickEvent.Referer,
			Country:   r.Country,
			City:      r.City,
			CreatedAt: clickEvent.CreatedAt.AsTime(),
		})

		if err != nil {
			h.Warnf("consume - store to database failed: %v", err)
		}

		return err
	}

	subscriber, err := rabbitmq.NewSubscriber(c, h, consume)
	if err != nil {
		return nil
	}
	return subscriber
}
