package data

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoBIStore struct {
	db  *mongo.Client
	ctx context.Context
}

func NewMongoBIStore(dsn string, username string, password string) (*MongoBIStore, error) {
	credential := options.Credential{
		Username: username,
		Password: password,
	}
	db, err := mongo.Connect(
		context.Background(),
		options.Client().ApplyURI(dsn).SetAuth(credential),
	)
	if err != nil {
		return nil, err
	}
	return &MongoBIStore{db: db, ctx: context.Background()}, nil
}

func (s *MongoBIStore) Open() error {
	collection := s.db.Database("bingo").Collection("clicks")
	_, err := collection.Indexes().CreateOne(
		s.ctx,
		mongo.IndexModel{
			Keys: bson.D{
				primitive.E{Key: "alias", Value: 1},
			},
			Options: options.Index().SetUnique(false),
		},
	)

	if err != nil {
		return err
	}

	return nil
}

func (s *MongoBIStore) Close() error {
	return s.db.Disconnect(s.ctx)
}

func (s *MongoBIStore) Create(click *Click) error {
	collection := s.db.Database("bingo").Collection("clicks")
	_, err := collection.InsertOne(s.ctx, click)
	return err
}

func (s *MongoBIStore) Clicks(alias string) (uint64, error) {
	collection := s.db.Database("bingo").Collection("clicks")
	clicks, err := collection.CountDocuments(s.ctx, bson.M{"alias": alias})
	return uint64(clicks), err
}
