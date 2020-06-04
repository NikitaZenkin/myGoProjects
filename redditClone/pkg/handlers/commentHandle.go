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

type CreateCooment struct {
	Comment string
}

func (handler *PostHandler) NewComment(w http.ResponseWriter, r *http.Request) {
	session, exist := session.SessionFromContext(r.Context())
	if !exist {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	preComment := CreateCooment{}
	err := json.Unmarshal(body, &preComment)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	newComment := post.Comment{
		Crated: time.Now(),
		Author: post.Author{
			UserName: session.UserName,
			ID:       session.UserId,
		},
		Body: preComment.Comment,
		ID:   system.NewId(),
	}
	vars := mux.Vars(r)
	postId := vars["post_id"]
	curPost, err := handler.PostRepo.AddCommentToPost(postId, newComment)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	postByte, _ := json.Marshal(curPost)
	w.WriteHeader(http.StatusCreated)
	w.Write(postByte)
}

func (handler *PostHandler) DeleteComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postId := vars["post_id"]
	commentId := vars["comment_id"]
	curPost, err := handler.PostRepo.DeleteCommentFromPost(postId, commentId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	postByte, _ := json.Marshal(curPost)
	w.Write(postByte)
}
