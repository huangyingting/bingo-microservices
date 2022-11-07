package service

import (
	"context"
	"fmt"

	v1 "bingo/api/bi/v1"
	"bingo/app/bi/internal/conf"
	"bingo/app/bi/internal/data"

	"google.golang.org/protobuf/types/known/emptypb"
)

// BIService is a BI service.
type BIService struct {
	v1.UnimplementedBIServer
	store data.IBIStore
}

// NewBIService new a BI service.
func NewBIService(c *conf.Store) (*BIService, error) {
	var store data.IBIStore = nil
	var err error = nil
	switch c.Driver {
	case "sqlite":
		dsn := "file::memory:?cache=shared"
		store, err = data.NewSqliteBIStore(dsn)
	case "mysql":
		dsn := fmt.Sprintf(
			"%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
			c.Username,
			c.Password,
			c.Host,
			c.Port,
			c.Database)
		store, err = data.NewMysqlBIStore(dsn)
	case "postgres":
		dsn := fmt.Sprintf(
			"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			c.Host,
			c.Port,
			c.Username,
			c.Password,
			c.Database)
		store, err = data.NewPostgresBIStore(dsn)
	case "sqlserver":
		dsn := fmt.Sprintf(
			"sqlserver://%s:%s@%s:%d?data=%s",
			c.Username,
			c.Password,
			c.Host,
			c.Port,
			c.Database)
		store, err = data.NewSQLServerBIStore(dsn)
	case "mongo":
		dsn := fmt.Sprintf(
			"mongodb://%s:%d/%s",
			c.Host,
			c.Port,
			c.Database,
		)
		store, err = data.NewMongoBIStore(dsn, c.Username, c.Password)
	default:
		return nil, fmt.Errorf("unsupported store driver: %s", c.Driver)
	}

	if err != nil {
		return nil, err
	}

	if err = store.Open(); err != nil {
		return nil, err
	}

	return &BIService{store: store}, nil
}

// Clicks implements Clicks.Clicks.
func (s *BIService) Clicks(
	ctx context.Context,
	in *v1.ClicksRequest,
) (*v1.ClicksReply, error) {

	clicks, err := s.store.Clicks(in.Alias)
	if err != nil {
		return nil, v1.ErrorDbError(err.Error())
	}
	return &v1.ClicksReply{Clicks: clicks}, nil
}

func (s *BIService) Create(ctx context.Context, click *data.Click) error {
	return s.store.Create(click)
}

func (s *BIService) Readiness(
	ctx context.Context,
	in *emptypb.Empty,
) (*v1.StatusReply, error) {
	return &v1.StatusReply{
		Status: "OK",
	}, nil
}

func (s *BIService) Liveness(
	ctx context.Context,
	in *emptypb.Empty,
) (*v1.StatusReply, error) {
	return &v1.StatusReply{
		Status: "OK",
	}, nil
}
