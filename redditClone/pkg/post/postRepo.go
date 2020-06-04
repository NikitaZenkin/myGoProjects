package post

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type PostRepo struct {
	Mongo CollectionHelper
	ctx   context.Context
}

func NewPostRepo(mongo *mongo.Collection) PostRepo {
	context, _ := context.WithTimeout(context.Background(), 1000*time.Second)
	col := mongoCollection{sr: mongo}
	return PostRepo{
		Mongo: &col,
		ctx:   context,
	}
}

func (repa PostRepo) GetPostFromDB(id string) (*Post, error) {
	var post Post
	err := repa.Mongo.FindOne(repa.ctx, bson.M{"_id": id}).Decode(&post)
	return &post, err
}

func (repa PostRepo) DeletePost(id string) error {
	_, err := repa.Mongo.DeleteOne(repa.ctx, bson.M{"_id": id})
	return err
}

func (repa PostRepo) FindPost(id string) (*Post, error) {
	post, err := repa.GetPostFromDB(id)
	if err != nil {
		return nil, err
	}
	post.Views++
	bytePost, _ := bson.Marshal(*post)
	repa.Mongo.FindOneAndReplace(repa.ctx, bson.M{"_id": id}, bytePost)
	return post, err
}

func (repa PostRepo) GetAllPosts() (*[]Post, error) {
	var posts []Post
	var post Post
	cursor, err := repa.Mongo.Find(repa.ctx, bson.M{})
	if err != nil {
		return &posts, err
	}
	for cursor.Next(repa.ctx) {
		post = Post{}
		err = cursor.Decode(&post)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return &posts, err
}

func (repa PostRepo) GetAllPostsByCategory(category string) (*[]Post, error) {
	var posts []Post
	var post Post
	cursor, err := repa.Mongo.Find(repa.ctx, bson.M{"category": category})
	if err != nil {
		return nil, err
	}
	for cursor.Next(repa.ctx) {
		post = Post{}
		err = cursor.Decode(&post)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return &posts, nil
}

func (repa PostRepo) GetAllPostsByLogin(login string) (*[]Post, error) {
	var posts []Post
	var post Post
	cursor, err := repa.Mongo.Find(repa.ctx, bson.M{"author.username": login})
	if err != nil {
		return nil, err
	}
	for cursor.Next(repa.ctx) {
		post = Post{}
		err = cursor.Decode(&post)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return &posts, err
}

func (repa PostRepo) AddPost(newPost *Post) error {
	bytePost, err := bson.Marshal(newPost)
	if err != nil {
		return err
	}
	_, err = repa.Mongo.InsertOne(repa.ctx, bytePost)
	return err
}

func (repa PostRepo) AddCommentToPost(posiId string, newComment Comment) (*Post, error) {
	post, err := repa.GetPostFromDB(posiId)
	if err != nil {
		return nil, err
	}
	post.AddComment(newComment)
	bytePost, _ := bson.Marshal(post)
	repa.Mongo.FindOneAndReplace(repa.ctx, bson.M{"_id": posiId}, bytePost)
	return post, err
}

func (repa PostRepo) DeleteCommentFromPost(posiId string, comentId string) (*Post, error) {
	post, err := repa.GetPostFromDB(posiId)
	if err != nil {
		return nil, err
	}
	post.DeleteComment(comentId)
	bytePost, _ := bson.Marshal(post)
	repa.Mongo.FindOneAndReplace(repa.ctx, bson.M{"_id": posiId}, bytePost)
	return post, err
}

func (repa PostRepo) UpvotePost(posiId string, userId string) (*Post, error) {
	post, err := repa.GetPostFromDB(posiId)
	if err != nil {
		return nil, err
	}
	post.Upvote(userId)
	bytePost, _ := bson.Marshal(post)
	repa.Mongo.FindOneAndReplace(repa.ctx, bson.M{"_id": posiId}, bytePost)
	return post, err
}

func (repa PostRepo) UnvotePost(posiId string, userId string) (*Post, error) {
	post, err := repa.GetPostFromDB(posiId)
	if err != nil {
		return nil, err
	}
	post.Unvote(userId)
	bytePost, _ := bson.Marshal(post)
	repa.Mongo.FindOneAndReplace(repa.ctx, bson.M{"_id": posiId}, bytePost)
	return post, err
}

func (repa PostRepo) DownvotePost(posiId string, userId string) (*Post, error) {
	post, err := repa.GetPostFromDB(posiId)
	if err != nil {
		return nil, err
	}
	post.Downvote(userId)
	bytePost, _ := bson.Marshal(post)
	repa.Mongo.FindOneAndReplace(repa.ctx, bson.M{"_id": posiId}, bytePost)
	return post, err
}
