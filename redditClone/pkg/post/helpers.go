package post

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
)

type CollectionHelper interface {
	FindOne(ctx context.Context, filter interface{}) SingleResultHelper
	InsertOne(ctx context.Context, filter interface{}) (interface{}, error)
	DeleteOne(ctx context.Context, filter interface{}) (int64, error)
	Find(ctx context.Context, filter interface{}) (CursorHelper, error)
	FindOneAndReplace(ctx context.Context, filter interface{}, v interface{})
}

type SingleResultHelper interface {
	Decode(v interface{}) error
}

type CursorHelper interface {
	Next(ctx context.Context) bool
	Decode(v interface{}) error
}

type mongoCollection struct {
	sr *mongo.Collection
}

type mongoSingleResult struct {
	sr *mongo.SingleResult
}

type mongoCursor struct {
	sr *mongo.Cursor
}

func (mc *mongoCollection) FindOne(ctx context.Context, filter interface{}) SingleResultHelper {
	singleResult := mc.sr.FindOne(ctx, filter)
	return &mongoSingleResult{sr: singleResult}
}

func (mc *mongoCollection) InsertOne(ctx context.Context, filter interface{}) (interface{}, error) {
	id, err := mc.sr.InsertOne(ctx, filter)
	return id.InsertedID, err
}

func (mc *mongoCollection) DeleteOne(ctx context.Context, filter interface{}) (int64, error) {
	count, err := mc.sr.DeleteOne(ctx, filter)
	return count.DeletedCount, err
}

func (mc *mongoCollection) Find(ctx context.Context, filter interface{}) (CursorHelper, error) {
	cur, err := mc.sr.Find(ctx, filter)
	return &mongoCursor{sr: cur}, err
}

func (mc *mongoCollection) FindOneAndReplace(ctx context.Context, filter interface{}, v interface{}) {
	mc.sr.FindOneAndReplace(ctx, filter, v)
}

func (sr *mongoSingleResult) Decode(v interface{}) error {
	return sr.sr.Decode(v)
}

func (sr *mongoCursor) Next(ctx context.Context) bool {
	return sr.sr.Next(ctx)
}

func (sr *mongoCursor) Decode(v interface{}) error {
	return sr.sr.Decode(v)
}
