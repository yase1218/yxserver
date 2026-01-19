package dao

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type MongoReader struct {
	client *mongo.Client
}

func NewMongoReader(client *mongo.Client) *MongoReader {
	return &MongoReader{client: client}
}

func (m *MongoReader) FindOne(ctx context.Context, database, collection string, filter interface{}) *mongo.SingleResult {
	defer func() {
		ctx.Done()
	}()
	return m.client.Database(database).Collection(collection).FindOne(ctx, filter)
}

func (m *MongoReader) Find(ctx context.Context, database, collection string, filter interface{}) (*mongo.Cursor, error) {
	defer func() {
		ctx.Done()
	}()
	return m.client.Database(database).Collection(collection).Find(ctx, filter)
}

func (m *MongoReader) Aggregate(ctx context.Context, database, collection string, pipeline interface{}) (*mongo.Cursor, error) {
	defer func() {
		ctx.Done()
	}()
	return m.client.Database(database).Collection(collection).Aggregate(ctx, pipeline)
}
