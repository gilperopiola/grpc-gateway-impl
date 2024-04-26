package mongo

import (
	"context"

	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoDB struct {
	DB   *mongo.Client
	Name string

	core.MongoDatabaseAPI
}

func NewMongoDB(dbName string, mongoClient *mongo.Client) *mongoDB {
	return &mongoDB{
		DB:   mongoClient,
		Name: dbName,
	}
}

func (m *mongoDB) Close(ctx context.Context) {
	m.DB.Disconnect(ctx)
}

func (m *mongoDB) InsertOne(ctx context.Context, colName string, document any) (*mongo.InsertOneResult, error) {
	collection := m.DB.Database(m.Name).Collection(colName)
	result, err := collection.InsertOne(ctx, document)
	return result, err
}

func (m *mongoDB) Find(ctx context.Context, colName string, filter any, limit, offset int) (*mongo.Cursor, error) {
	opts := options.Find().SetLimit(int64(limit)).SetSkip(int64(offset))
	collection := m.DB.Database(m.Name).Collection(colName)
	cursor, err := collection.Find(ctx, filter, opts)
	return cursor, err
}

func (m *mongoDB) FindOne(ctx context.Context, colName string, filter any) *mongo.SingleResult {
	collection := m.DB.Database(m.Name).Collection(colName)
	result := collection.FindOne(ctx, filter)
	return result
}

func (m *mongoDB) DeleteOne(ctx context.Context, colName string, filter any) (*mongo.DeleteResult, error) {
	collection := m.DB.Database(m.Name).Collection(colName)
	result, err := collection.DeleteOne(ctx, filter)
	return result, err
}

func (m *mongoDB) Count(ctx context.Context, colName string, filter any) (int64, error) {
	collection := m.DB.Database(m.Name).Collection(colName)
	count, err := collection.CountDocuments(ctx, filter)
	return count, err
}
