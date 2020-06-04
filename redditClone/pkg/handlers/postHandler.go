package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"rclone/pkg/post"
	"rclone/pkg/session"
	"rclone/pkg/system"
	"time"
)

type PostRepositoryInterface interface {
	DeletePost(id string) error
	FindPost(id string) (*post.Post, error)
	GetAllPosts() (*[]post.Post, error)
	GetAllPostsByCategory(category string) (*[]post.Post, error)
	GetAllPostsByLogin(login string) (*[]post.Post, error)
	AddPost(newPost *post.Post) error
	AddCommentToPost(posiId string, newComment post.Comment) (*post.Post, error)
	DeleteCommentFromPost(posiId string, comentId string) (*post.Post, error)
	UpvotePost(posiId string, userId string) (*post.Post, error)
	UnvotePost(posiId string, userId string) (*post.Post, error)
	DownvotePost(posiId string, userId string) (*post.Post, error)
}

type PostHandler struct {
	PostRepo PostRepositoryInterface
}

type CreatePost struct {
	Category string `json:"category"`
	Type     string `json:"type"`
	Title    string `json:"title"`
	Text     string `json:"text"`
	URL      string `json:"url"`
}

func (handler *PostHandler) NewPost(w http.ResponseWriter, r *http.Request) {
	session, exist := session.SessionFromContext(r.Context())
	if !exist {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	prePost := CreatePost{}
	err := json.Unmarshal(body, &prePost)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	newPost := post.Post{
		Score: 1,
		Views: 0,
		Type:  prePost.Type,
		Title: prePost.Title,
		Author: post.Author{
			UserName: session.UserName,
			ID:       session.UserId,
		},
		Category: prePost.Category,
		Text:     prePost.Text,
		Url:      prePost.URL,
		Votes: []post.Vote{post.Vote{
			User: session.UserId,
			Vote: 1,
		}},
		Comments:         []post.Comment{},
		Crated:           time.Now(),
		UpvotePercentage: 0,
		ID:               system.NewId(),
	}
	postByte, _ := json.Marshal(newPost)
	newPost.CalculateUpvotePercentage()
	err = handler.PostRepo.AddPost(&newPost)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(postByte)
}

func (handler *PostHandler) ShowPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	strinId := vars["id"]
	curPost, err := handler.PostRepo.FindPost(strinId)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	postByte, _ := json.Marshal(curPost)
	w.Write(postByte)
}

func (handler *PostHandler) DeletePost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	strinId := vars["id"]
	err := handler.PostRepo.DeletePost(strinId)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	mes := system.Message{Message: "success"}
	jsonMes := mes.ToJson()
	w.Write(jsonMes)
	return
}

func (handler *PostHandler) ShowAllPosts(w http.ResponseWriter, r *http.Request) {
	posts, err := handler.PostRepo.GetAllPosts()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	postsByte, _ := json.Marshal(*posts)
	w.Write(postsByte)
}

func (handler *PostHandler) ShowAllUserPosts(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	login := vars["user_login"]
	posts, err := handler.PostRepo.GetAllPostsByLogin(login)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	postsByte, _ := json.Marshal(*posts)
	w.Write(postsByte)
}

func (handler *PostHandler) ShowAllCategoryPosts(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	category := vars["category"]
	posts, err := handler.PostRepo.GetAllPostsByCategory(category)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	postsByte, _ := json.Marshal(*posts)
	w.Write(postsByte)
}
