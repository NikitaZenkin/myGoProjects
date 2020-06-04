package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"rclone/pkg/session"
)

func (handler *PostHandler) Upvote(w http.ResponseWriter, r *http.Request) {
	session, exist := session.SessionFromContext(r.Context())
	if !exist {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	vars := mux.Vars(r)
	postId := vars["post_id"]
	curPost, err := handler.PostRepo.UpvotePost(postId, session.UserId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	postByte, _ := json.Marshal(curPost)
	w.Write(postByte)
}

func (handler *PostHandler) Downvote(w http.ResponseWriter, r *http.Request) {
	session, exist := session.SessionFromContext(r.Context())
	if !exist {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	vars := mux.Vars(r)
	postId := vars["post_id"]
	curPost, err := handler.PostRepo.DownvotePost(postId, session.UserId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	postByte, _ := json.Marshal(curPost)
	w.Write(postByte)
}

func (handler *PostHandler) Unvote(w http.ResponseWriter, r *http.Request) {
	session, exist := session.SessionFromContext(r.Context())
	if !exist {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	vars := mux.Vars(r)
	postId := vars["post_id"]
	curPost, err := handler.PostRepo.UnvotePost(postId, session.UserId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	postByte, _ := json.Marshal(curPost)
	w.Write(postByte)
}
