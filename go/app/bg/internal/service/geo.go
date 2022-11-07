package service

import (
	"bytes"
	"context"
	"encoding/binary"
	"net"
	"strconv"
	"strings"

	v1 "bingo/api/bg/v1"
	"bingo/app/bg/internal/conf"

	"github.com/go-redis/redis/v8"
	"google.golang.org/protobuf/types/known/emptypb"
)

type GeoService struct {
	v1.UnimplementedGeoServer
	redis *redis.Client
}

func NewGeoService(c *conf.Redis) (*GeoService, error) {
	var addr string = c.Addr
	if len(addr) == 0 {
		addr = "127.0.0.1:6379"
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	err := redisClient.Ping(context.Background()).Err()
	if err != nil {
		return nil, err
	}
	return &GeoService{redis: redisClient}, nil
}

func (s *GeoService) Location(
	ctx context.Context,
	in *v1.LocationRequest,
) (*v1.LocationReply, error) {

	net_ip := net.ParseIP(in.Ip)
	if net_ip == nil {
		return nil, v1.ErrorInvalidIpAddress("invalid ip address: %s", in.Ip)
	}
	var uint32_ip uint32
	binary.Read(bytes.NewBuffer(net_ip.To4()), binary.BigEndian, &uint32_ip)

	r := s.redis.ZRangeByScore(
		context.Background(),
		"ipv4",
		&redis.ZRangeBy{
			Min:    strconv.FormatUint(uint64(uint32_ip), 10),
			Max:    "+inf",
			Offset: 0,
			Count:  1,
		},
	).Val()

	if len(r) == 0 {
		return nil, v1.ErrorLocationNotFound("location not found: %s", in.Ip)
	}

	geo := strings.Split(r[0], "|")

	if len(geo) != 3 {
		return nil, v1.ErrorInternalServerError("internal server error: bad geo format")
	}

	if geo[1] == "" && geo[2] == "" {
		return nil, v1.ErrorLocationNotFound("location not found: %s", in.Ip)
	}
	return &v1.LocationReply{
		Ip:      in.Ip,
		Country: geo[1],
		City:    geo[2],
	}, nil
}

func (s *GeoService) Readiness(
	ctx context.Context,
	in *emptypb.Empty,
) (*v1.StatusReply, error) {
	return &v1.StatusReply{
		Status: "OK",
	}, nil
}

func (s *GeoService) Liveness(
	ctx context.Context,
	in *emptypb.Empty,
) (*v1.StatusReply, error) {
	return &v1.StatusReply{
		Status: "OK",
	}, nil
}
