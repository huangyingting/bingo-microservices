package data

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoBSStore struct {
	client   *mongo.Client
	database string
	ctx      context.Context
}

func NewMongoBSStore(
	dsn string,
	database string,
	username string,
	password string,
) (*MongoBSStore, error) {
	credential := options.Credential{
		Username: username,
		Password: password,
	}
	client, err := mongo.Connect(
		context.Background(),
		options.Client().ApplyURI(dsn).SetAuth(credential),
	)
	if err != nil {
		return nil, err
	}
	return &MongoBSStore{client: client, database: database, ctx: context.Background()}, nil
}

/*
func (q *MongoBSStore) Open() error {
	session, err := q.db.StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(q.ctx)

	err = mongo.WithSession(q.ctx, session, func(sessionContext mongo.SessionContext) error {
		if err = session.StartTransaction(); err != nil {
			return err
		}
		collection := q.db.Database("bingo").Collection("short_url")
		_, err := collection.Indexes().CreateOne(
			sessionContext,
			mongo.IndexModel{
				Keys: bson.D{
					primitive.E{Key: "alias", Value: 1},
					primitive.E{Key: "oid", Value: 1},
				},
				Options: options.Index().SetUnique(true),
			},
		)

		if err != nil {
			return err
		}

		collection = q.db.Database("bingo").Collection("visit")
		_, err = collection.Indexes().CreateOne(
			sessionContext,
			mongo.IndexModel{
				Keys: bson.D{
					primitive.E{Key: "alias", Value: 1},
				},
				Options: options.Index().SetUnique(true),
			},
		)

		if err != nil {
			return err
		}

		if err = session.CommitTransaction(sessionContext); err != nil {
			return err
		}
		return nil
	})
	return err
}
*/

func (q *MongoBSStore) Open() error {
	collection := q.client.Database(q.database).Collection("short_url")
	_, err := collection.Indexes().CreateOne(
		q.ctx,
		mongo.IndexModel{
			Keys: bson.D{
				primitive.E{Key: "alias", Value: 1},
				primitive.E{Key: "oid", Value: 1},
			},
			Options: options.Index().SetUnique(true),
		},
	)

	if err != nil {
		return err
	}

	return nil
}

func (q *MongoBSStore) Close() error {
	return q.client.Disconnect(q.ctx)
}

func (q *MongoBSStore) CreateShortUrl(alias string, customized bool, url string, oid string) error {
	flags := Bits(0).Set(FLAG_CUSTOMIZED, customized)
	collection := q.client.Database(q.database).Collection("short_url")
	_, err := collection.InsertOne(
		q.ctx,
		ShortUrl{Alias: alias, Url: url, Oid: oid, Flags: flags, CreatedAt: time.Now()},
	)
	return err
}

func (q *MongoBSStore) DeleteShortUrl(alias string, oid string) error {
	collection := q.client.Database(q.database).Collection("short_url")
	result, err := collection.DeleteOne(q.ctx, bson.M{"alias": alias, "oid": oid})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		err = ErrNoRowsDeleted
	}
	return err
}

func (q *MongoBSStore) GetShortUrl(alias string) (*ShortUrl, error) {
	shortUrl := ShortUrl{}
	collection := q.client.Database(q.database).Collection("short_url")
	err := collection.FindOne(q.ctx, bson.M{"alias": alias}).Decode(&shortUrl)
	return &shortUrl, err
}

func (q *MongoBSStore) GetShortUrlByOid(alias string, oid string) (*ShortUrl, error) {
	shortUrl := ShortUrl{}
	collection := q.client.Database(q.database).Collection("short_url")
	err := collection.FindOne(q.ctx, bson.M{"alias": alias, "oid": oid}).Decode(&shortUrl)
	return &shortUrl, err
}

func (q *MongoBSStore) ListShortUrl(oid string, limit int64, offset int64) ([]*ShortUrl, error) {
	collection := q.client.Database(q.database).Collection("short_url")
	findOpt := options.FindOptions{Limit: &limit, Skip: &offset}
	cur, err := collection.Find(q.ctx, bson.M{"oid": oid}, &findOpt)
	if err != nil {
		return nil, err
	}
	var shortUrls []*ShortUrl
	for cur.Next(q.ctx) {
		var shortUrl ShortUrl
		if err := cur.Decode(&shortUrl); err != nil {
			return nil, err
		}
		shortUrls = append(shortUrls, &shortUrl)
	}
	return shortUrls, nil
}

func (q *MongoBSStore) UpdateShortUrl(
	alias string,
	oid string,
	updateShortUrl UpdateShortUrl,
) error {
	collection := q.client.Database(q.database).Collection("short_url")
	shortUrl := ShortUrl{}
	err := collection.FindOne(q.ctx, bson.M{"alias": alias}).Decode(&shortUrl)
	if err != nil {
		return err
	}

	flags := (shortUrl.Flags & 1) + updateShortUrl.Flags
	result, err := collection.UpdateOne(
		q.ctx,
		bson.M{"alias": alias, "oid": oid},
		bson.M{
			"$set": bson.M{
				"url":          updateShortUrl.Url,
				"title":        updateShortUrl.Title,
				"tags":         updateShortUrl.Tags,
				"flags":        flags,
				"utm_source":   updateShortUrl.UtmSource,
				"utm_medium":   updateShortUrl.UtmMedium,
				"utm_campaign": updateShortUrl.UtmCampaign,
				"utm_term":     updateShortUrl.UtmTerm,
				"utm_content":  updateShortUrl.UtmContent,
			},
		},
	)
	if result.ModifiedCount == 0 {
		return ErrNoRowsUpdated
	}
	return err
}
