package main

import (
	"os"
	"time"

	v1 "bingo/api/bi/v1"
	"bingo/pkg/rabbitmq"

	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	Name string = "bi-c"
	// Version is the version of the compiled software.
	Version string
	// flagconf is the config flag.
	flagconf string

	id, _ = os.Hostname()
)

func main() {
	logger := log.With(log.NewStdLogger(os.Stdout),
		"ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
		"service.id", id,
		"service.name", Name,
		"service.version", Version,
	)

	h := log.NewHelper(log.NewFilter(logger, log.FilterLevel(log.LevelDebug)))

	srv, err := rabbitmq.NewPublisher(&rabbitmq.PublisherConfig{
		AmqpUri:        "amqp://guest:guest@localhost:5672/",
		Name:           "",
		CaCert:         "",
		ClientCert:     "",
		ClientKey:      "",
		ExchangeName:   "bingo",
		ExchangeType:   "direct",
		QueueMode:      "default",
		RoutingKey:     "clicks",
		RetryAttempt:   3,
		RetryInterval:  durationpb.New(5 * time.Second),
		ConnectTimeout: durationpb.New(5 * time.Second),
	}, h)

	if err != nil {
		h.Errorf("new publisher error: %v", err)
		return
	}

	var i int = 0
	clickEvent := v1.ClickEvent{
		Alias:     "microsoft",
		Ip:        "167.220.255.52",
		Ua:        "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.4951.41 Safari/537.36",
		Referer:   "https://www.microsoft.com",
		CreatedAt: timestamppb.New(time.Now()),
	}

	for {
		i++
		bytes, err := proto.Marshal(&clickEvent)
		if err == nil {
			srv.Publish(bytes)
		}
		time.Sleep(1 * time.Second)
	}

}
