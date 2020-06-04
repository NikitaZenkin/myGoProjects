package main

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"html/template"
	"math/rand"
	"net/http"
	"rclone/pkg/handlers"
	"rclone/pkg/middleware"
	"rclone/pkg/post"
	"rclone/pkg/session"
	"rclone/pkg/user"
	"time"
)

func main() {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	collection := client.Database("rclone").Collection("posts")

	dsn := "root:qwerty@tcp(localhost:3306)/rclone?"
	dsn += "&charset=utf8"
	dsn += "&interpolateParams=true"
	db, err := sql.Open("mysql", dsn)
	db.SetMaxOpenConns(10)
	err = db.Ping()
	if err != nil {
		fmt.Println(err)
		return
	}

	rand.Seed(time.Now().UnixNano())
	templates := template.Must(template.ParseFiles("./template/index.html"))
	staticHandler := http.StripPrefix(
		"/static/",
		http.FileServer(http.Dir("./template/static/")),
	)
	http.Handle("/static/", staticHandler)

	userRepoz := user.NewUserRepo(db)
	postRepoz := post.NewPostRepo(collection)
	sesManager := session.NewSesManager(db)
	userHandler := &handlers.UserHandler{
		Tmpl:     templates,
		UserRepo: &userRepoz,
		Sessions: &sesManager,
	}
	postHandler := &handlers.PostHandler{
		PostRepo: postRepoz,
	}

	r := mux.NewRouter()

	r.HandleFunc("/", userHandler.Template).Methods("GET")
	r.HandleFunc("/api/login", userHandler.Login).Methods("POST")
	r.HandleFunc("/api/register", userHandler.Register).Methods("POST")

	r.HandleFunc("/api/posts", postHandler.NewPost).Methods("POST")
	r.HandleFunc("/api/post/{id}", postHandler.ShowPost).Methods("GET")
	r.HandleFunc("/api/posts/", postHandler.ShowAllPosts).Methods("GET")
	r.HandleFunc("/api/posts/{category}", postHandler.ShowAllCategoryPosts).Methods("GET")
	r.HandleFunc("/api/user/{user_login}", postHandler.ShowAllUserPosts).Methods("GET")
	r.HandleFunc("/api/post/{id}", postHandler.DeletePost).Methods("DELETE")

	r.HandleFunc("/api/post/{post_id}", postHandler.NewComment).Methods("POST")
	r.HandleFunc("/api/post/{post_id}/{comment_id}", postHandler.DeleteComment).Methods("DELETE")

	r.HandleFunc("/api/post/{post_id}/upvote", postHandler.Upvote).Methods("GET")
	r.HandleFunc("/api/post/{post_id}/downvote", postHandler.Downvote).Methods("GET")
	r.HandleFunc("/api/post/{post_id}/unvote", postHandler.Unvote).Methods("GET")

	muxH := middleware.Auth(&sesManager, r)
	muxH = middleware.Panic(muxH)
	http.ListenAndServe(":8080", muxH)
}
