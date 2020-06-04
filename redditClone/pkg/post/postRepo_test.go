package post

import (
	"context"
	"fmt"
	"github.com/golang/mock/gomock"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"testing"
	"time"
)

func TestNewPostRepo(t *testing.T) {
	NewPostRepo(&mongo.Collection{})
}

func TestPostRepoDeletePost(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	context, _ := context.WithTimeout(context.Background(), 10*time.Second)
	fakeCollection := NewMockCollectionHelper(ctrl)
	servise := &PostRepo{Mongo: fakeCollection, ctx: context}

	fakeCollection.EXPECT().DeleteOne(context, bson.M{"_id": "1"}).Return(int64(0), nil)
	err := servise.DeletePost("1")
	if err != nil {
		t.Errorf("unexpected error")
		return
	}
}

func TestPostRepoFindPost(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	context, _ := context.WithTimeout(context.Background(), 10*time.Second)
	fakeCollection := NewMockCollectionHelper(ctrl)
	fakeSingle := NewMockSingleResultHelper(ctrl)
	servise := &PostRepo{Mongo: fakeCollection, ctx: context}

	fakeCollection.EXPECT().FindOne(context, bson.M{"_id": "2"}).Return(fakeSingle)
	fakeSingle.EXPECT().Decode(&Post{}).Return(fmt.Errorf("didn't decode"))
	_, err := servise.FindPost("2")
	if err == nil {
		t.Errorf("expected error didn't decode")
		return
	}

	fakeCollection.EXPECT().FindOne(context, bson.M{"_id": "2"}).Return(fakeSingle)
	fakeSingle.EXPECT().Decode(&Post{}).Return(nil)
	fakeCollection.EXPECT().FindOneAndReplace(context, bson.M{"_id": "2"}, gomock.Any())
	_, err = servise.FindPost("2")
	if err != nil {
		t.Errorf("unexpected error")
		return
	}
}

func TestPostRepoGetAllPosts(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	context, _ := context.WithTimeout(context.Background(), 10*time.Second)
	fakeCollection := NewMockCollectionHelper(ctrl)
	fakeCursor := NewMockCursorHelper(ctrl)
	servise := &PostRepo{Mongo: fakeCollection, ctx: context}

	fakeCollection.EXPECT().Find(context, bson.M{}).Return(fakeCursor, fmt.Errorf("didn't find"))
	_, err := servise.GetAllPosts()
	if err == nil {
		t.Errorf("expected error didn't find")
		return
	}

	fakeCollection.EXPECT().Find(context, bson.M{}).Return(fakeCursor, nil)
	fakeCursor.EXPECT().Next(context).Return(true)
	fakeCursor.EXPECT().Decode(&Post{}).Return(fmt.Errorf("didn't decode"))
	_, err = servise.GetAllPosts()
	if err == nil {
		t.Errorf("expected error didn't decode")
		return
	}

	fakeCollection.EXPECT().Find(context, bson.M{}).Return(fakeCursor, nil)
	fakeCursor.EXPECT().Next(context).Return(true)
	fakeCursor.EXPECT().Decode(&Post{}).Return(nil)
	fakeCursor.EXPECT().Next(context).Return(false)
	_, err = servise.GetAllPosts()
	if err != nil {
		t.Errorf("unexpexted error")
		return
	}
}

func TestPostRepoGetAllPostsGyLogin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	context, _ := context.WithTimeout(context.Background(), 10*time.Second)
	fakeCollection := NewMockCollectionHelper(ctrl)
	fakeCursor := NewMockCursorHelper(ctrl)
	servise := &PostRepo{Mongo: fakeCollection, ctx: context}

	fakeCollection.EXPECT().Find(context, bson.M{"author.username": "Nikita"}).Return(fakeCursor, fmt.Errorf("didn't find"))
	_, err := servise.GetAllPostsByLogin("Nikita")
	if err == nil {
		t.Errorf("expected error didn't find")
		return
	}

	fakeCollection.EXPECT().Find(context, bson.M{"author.username": "Nikita"}).Return(fakeCursor, nil)
	fakeCursor.EXPECT().Next(context).Return(true)
	fakeCursor.EXPECT().Decode(&Post{}).Return(fmt.Errorf("didn't decode"))
	_, err = servise.GetAllPostsByLogin("Nikita")
	if err == nil {
		t.Errorf("expected error didn't decode")
		return
	}

	fakeCollection.EXPECT().Find(context, bson.M{"author.username": "Nikita"}).Return(fakeCursor, nil)
	fakeCursor.EXPECT().Next(context).Return(true)
	fakeCursor.EXPECT().Decode(&Post{}).Return(nil)
	fakeCursor.EXPECT().Next(context).Return(false)
	_, err = servise.GetAllPostsByLogin("Nikita")
	if err != nil {
		t.Errorf("unexpexted error")
		return
	}
}

func TestPostRepoGetAllPostsGyCategory(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	context, _ := context.WithTimeout(context.Background(), 10*time.Second)
	fakeCollection := NewMockCollectionHelper(ctrl)
	fakeCursor := NewMockCursorHelper(ctrl)
	servise := &PostRepo{Mongo: fakeCollection, ctx: context}

	fakeCollection.EXPECT().Find(context, bson.M{"category": "music"}).Return(fakeCursor, fmt.Errorf("didn't find"))
	_, err := servise.GetAllPostsByCategory("music")
	if err == nil {
		t.Errorf("expected error didn't find")
		return
	}

	fakeCollection.EXPECT().Find(context, bson.M{"category": "music"}).Return(fakeCursor, nil)
	fakeCursor.EXPECT().Next(context).Return(true)
	fakeCursor.EXPECT().Decode(&Post{}).Return(fmt.Errorf("didn't decode"))
	_, err = servise.GetAllPostsByCategory("music")
	if err == nil {
		t.Errorf("expected error didn't decode")
		return
	}

	fakeCollection.EXPECT().Find(context, bson.M{"category": "music"}).Return(fakeCursor, nil)
	fakeCursor.EXPECT().Next(context).Return(true)
	fakeCursor.EXPECT().Decode(&Post{}).Return(nil)
	fakeCursor.EXPECT().Next(context).Return(false)
	_, err = servise.GetAllPostsByCategory("music")
	if err != nil {
		t.Errorf("unexpexted error")
		return
	}
}

func TestPostRepoAddPost(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	context, _ := context.WithTimeout(context.Background(), 10*time.Second)
	fakeCollection := NewMockCollectionHelper(ctrl)
	servise := &PostRepo{Mongo: fakeCollection, ctx: context}

	err := servise.AddPost(nil)
	if err == nil {
		t.Errorf("expected error bad post")
		return
	}

	fakeCollection.EXPECT().InsertOne(context, gomock.Any()).Return(gomock.Nil(), nil)
	err = servise.AddPost(&Post{})
	if err != nil {
		t.Errorf("unexpected error")
		return
	}
}

func TestPostRepoAddCommentToPost(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	context, _ := context.WithTimeout(context.Background(), 10*time.Second)
	fakeCollection := NewMockCollectionHelper(ctrl)
	fakeSingle := NewMockSingleResultHelper(ctrl)
	servise := &PostRepo{Mongo: fakeCollection, ctx: context}

	fakeCollection.EXPECT().FindOne(context, bson.M{"_id": "2"}).Return(fakeSingle)
	fakeSingle.EXPECT().Decode(&Post{}).Return(fmt.Errorf("didn't decode"))
	_, err := servise.AddCommentToPost("2", Comment{})
	if err == nil {
		t.Errorf("expected error didn't decode")
		return
	}

	fakeCollection.EXPECT().FindOne(context, bson.M{"_id": "2"}).Return(fakeSingle)
	fakeSingle.EXPECT().Decode(&Post{}).Return(nil)
	fakeCollection.EXPECT().FindOneAndReplace(context, bson.M{"_id": "2"}, gomock.Any())
	_, err = servise.AddCommentToPost("2", Comment{})
	if err != nil {
		t.Errorf("unexpected error")
		return
	}
}

func TestPostRepoDeleteCommentFromPost(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	context, _ := context.WithTimeout(context.Background(), 10*time.Second)
	fakeCollection := NewMockCollectionHelper(ctrl)
	fakeSingle := NewMockSingleResultHelper(ctrl)
	servise := &PostRepo{Mongo: fakeCollection, ctx: context}

	fakeCollection.EXPECT().FindOne(context, bson.M{"_id": "2"}).Return(fakeSingle)
	fakeSingle.EXPECT().Decode(&Post{}).Return(fmt.Errorf("didn't decode"))
	_, err := servise.DeleteCommentFromPost("2", "3")
	if err == nil {
		t.Errorf("expected error didn't decode")
		return
	}

	fakeCollection.EXPECT().FindOne(context, bson.M{"_id": "2"}).Return(fakeSingle)
	fakeSingle.EXPECT().Decode(&Post{}).Return(nil)
	fakeCollection.EXPECT().FindOneAndReplace(context, bson.M{"_id": "2"}, gomock.Any())
	_, err = servise.DeleteCommentFromPost("2", "3")
	if err != nil {
		t.Errorf("unexpected error")
		return
	}
}

func TestPostRepoUpvotePost(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	context, _ := context.WithTimeout(context.Background(), 10*time.Second)
	fakeCollection := NewMockCollectionHelper(ctrl)
	fakeSingle := NewMockSingleResultHelper(ctrl)
	servise := &PostRepo{Mongo: fakeCollection, ctx: context}

	fakeCollection.EXPECT().FindOne(context, bson.M{"_id": "2"}).Return(fakeSingle)
	fakeSingle.EXPECT().Decode(&Post{}).Return(fmt.Errorf("didn't decode"))
	_, err := servise.UpvotePost("2", "3")
	if err == nil {
		t.Errorf("expected error didn't decode")
		return
	}

	fakeCollection.EXPECT().FindOne(context, bson.M{"_id": "2"}).Return(fakeSingle)
	fakeSingle.EXPECT().Decode(&Post{}).Return(nil)
	fakeCollection.EXPECT().FindOneAndReplace(context, bson.M{"_id": "2"}, gomock.Any())
	_, err = servise.UpvotePost("2", "3")
	if err != nil {
		t.Errorf("unexpected error")
		return
	}
}

func TestPostRepoUnvotePost(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	context, _ := context.WithTimeout(context.Background(), 10*time.Second)
	fakeCollection := NewMockCollectionHelper(ctrl)
	fakeSingle := NewMockSingleResultHelper(ctrl)
	servise := &PostRepo{Mongo: fakeCollection, ctx: context}

	fakeCollection.EXPECT().FindOne(context, bson.M{"_id": "2"}).Return(fakeSingle)
	fakeSingle.EXPECT().Decode(&Post{}).Return(fmt.Errorf("didn't decode"))
	_, err := servise.UnvotePost("2", "3")
	if err == nil {
		t.Errorf("expected error didn't decode")
		return
	}

	fakeCollection.EXPECT().FindOne(context, bson.M{"_id": "2"}).Return(fakeSingle)
	fakeSingle.EXPECT().Decode(&Post{}).Return(nil)
	fakeCollection.EXPECT().FindOneAndReplace(context, bson.M{"_id": "2"}, gomock.Any())
	_, err = servise.UnvotePost("2", "3")
	if err != nil {
		t.Errorf("unexpected error")
		return
	}
}

func TestPostRepoDownvotePost(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	context, _ := context.WithTimeout(context.Background(), 10*time.Second)
	fakeCollection := NewMockCollectionHelper(ctrl)
	fakeSingle := NewMockSingleResultHelper(ctrl)
	servise := &PostRepo{Mongo: fakeCollection, ctx: context}

	fakeCollection.EXPECT().FindOne(context, bson.M{"_id": "2"}).Return(fakeSingle)
	fakeSingle.EXPECT().Decode(&Post{}).Return(fmt.Errorf("didn't decode"))
	_, err := servise.DownvotePost("2", "3")
	if err == nil {
		t.Errorf("expected error didn't decode")
		return
	}

	fakeCollection.EXPECT().FindOne(context, bson.M{"_id": "2"}).Return(fakeSingle)
	fakeSingle.EXPECT().Decode(&Post{}).Return(nil)
	fakeCollection.EXPECT().FindOneAndReplace(context, bson.M{"_id": "2"}, gomock.Any())
	_, err = servise.DownvotePost("2", "3")
	if err != nil {
		t.Errorf("unexpected error")
		return
	}
}
