package handlers

import (
	"bytes"
	"context"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"net/http/httptest"
	"rclone/pkg/post"
	"rclone/pkg/session"
	"testing"
)

func TestPostHandlerNewPost(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	fakeRepa := NewMockPostRepositoryInterface(ctrl)
	service := &PostHandler{PostRepo: fakeRepa}

	req := httptest.NewRequest("POST", "/api/posts", nil)
	w := httptest.NewRecorder()
	service.NewPost(w, req)
	resp := w.Result()
	if resp.StatusCode != 401 {
		t.Errorf("expected resp status 401, got %d", resp.StatusCode)
		return
	}

	sess := session.Session{
		ID:       "5",
		UserId:   "35",
		UserName: "Nikita",
	}
	rqBody := bytes.Buffer{}
	ctx := context.WithValue(req.Context(), "sessionKey", &sess)
	rqBody.Write([]byte(`bad json`))
	req = httptest.NewRequest("POST", "/api/posts", &rqBody)
	w = httptest.NewRecorder()
	service.NewPost(w, req.WithContext(ctx))
	resp = w.Result()
	if resp.StatusCode != 500 {
		t.Errorf("expected resp status 500, got %d", resp.StatusCode)
		return
	}

	rqBody.Write([]byte(`{"category":"music","type":"text","title":"elctro","text":"beat"}`))
	req = httptest.NewRequest("POST", "/api/posts", &rqBody)
	w = httptest.NewRecorder()
	fakeRepa.EXPECT().AddPost(gomock.Any()).Return(fmt.Errorf("didn't add"))
	service.NewPost(w, req.WithContext(ctx))
	resp = w.Result()
	if resp.StatusCode != 500 {
		t.Errorf("expected resp status 500, got %d", resp.StatusCode)
		return
	}

	rqBody.Write([]byte(`{"category":"music","type":"text","title":"elctro","text":"beat"}`))
	req = httptest.NewRequest("POST", "/api/posts", &rqBody)
	w = httptest.NewRecorder()
	fakeRepa.EXPECT().AddPost(gomock.Any()).Return(nil)
	service.NewPost(w, req.WithContext(ctx))
	resp = w.Result()
	if resp.StatusCode != 201 {
		t.Errorf("expected resp status 201, got %d", resp.StatusCode)
		return
	}
}

func TestPostHandlerShowPost(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	fakeRepa := NewMockPostRepositoryInterface(ctrl)
	service := &PostHandler{PostRepo: fakeRepa}

	req := httptest.NewRequest("GET", "/api/post/{id}", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "346"})
	w := httptest.NewRecorder()
	fakeRepa.EXPECT().FindPost("346").Return(nil, fmt.Errorf("didn't find"))
	service.ShowPost(w, req)
	resp := w.Result()
	if resp.StatusCode != 404 {
		t.Errorf("expected resp status 404, got %d", resp.StatusCode)
		return
	}

	req = httptest.NewRequest("GET", "/api/post/{id}", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "346"})
	w = httptest.NewRecorder()
	fakeRepa.EXPECT().FindPost("346").Return(nil, nil)
	service.ShowPost(w, req)
	resp = w.Result()
	if resp.StatusCode != 200 {
		t.Errorf("expected resp status 200, got %d", resp.StatusCode)
		return
	}

}

func TestPostHandlerDeletePost(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	fakeRepa := NewMockPostRepositoryInterface(ctrl)
	service := &PostHandler{PostRepo: fakeRepa}

	req := httptest.NewRequest("DELETE", "/api/post/{id}", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "346"})
	w := httptest.NewRecorder()
	fakeRepa.EXPECT().DeletePost("346").Return(fmt.Errorf("didn't dlete"))
	service.DeletePost(w, req)
	resp := w.Result()
	if resp.StatusCode != 404 {
		t.Errorf("expected resp status 404, got %d", resp.StatusCode)
		return
	}

	req = httptest.NewRequest("GET", "/api/post/{id}", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "346"})
	w = httptest.NewRecorder()
	fakeRepa.EXPECT().DeletePost("346").Return(nil)
	service.DeletePost(w, req)
	resp = w.Result()
	if resp.StatusCode != 200 {
		t.Errorf("expected resp status 200, got %d", resp.StatusCode)
		return
	}

}

func TestPostHandlerShowAllPosts(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	fakeRepa := NewMockPostRepositoryInterface(ctrl)
	service := &PostHandler{PostRepo: fakeRepa}

	req := httptest.NewRequest("GET", "/api/posts/", nil)
	w := httptest.NewRecorder()
	fakeRepa.EXPECT().GetAllPosts().Return(&[]post.Post{}, fmt.Errorf("didn't get"))
	service.ShowAllPosts(w, req)
	resp := w.Result()
	if resp.StatusCode != 500 {
		t.Errorf("expected resp status 500, got %d", resp.StatusCode)
		return
	}

	req = httptest.NewRequest("GET", "/api/posts/", nil)
	w = httptest.NewRecorder()
	fakeRepa.EXPECT().GetAllPosts().Return(&[]post.Post{}, nil)
	service.ShowAllPosts(w, req)
	resp = w.Result()
	if resp.StatusCode != 200 {
		t.Errorf("expected resp status 200, got %d", resp.StatusCode)
		return
	}
}

func TestPostHandlerShowAllCategoryPosts(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	fakeRepa := NewMockPostRepositoryInterface(ctrl)
	service := &PostHandler{PostRepo: fakeRepa}

	req := httptest.NewRequest("GET", "/api/posts/{category}", nil)
	req = mux.SetURLVars(req, map[string]string{"category": "music"})
	w := httptest.NewRecorder()
	fakeRepa.EXPECT().GetAllPostsGyCategory("music").Return(&[]post.Post{}, fmt.Errorf("didn't get"))
	service.ShowAllCategoryPosts(w, req)
	resp := w.Result()
	if resp.StatusCode != 500 {
		t.Errorf("expected resp status 500, got %d", resp.StatusCode)
		return
	}

	req = httptest.NewRequest("GET", "/api/posts/{category}", nil)
	req = mux.SetURLVars(req, map[string]string{"category": "music"})
	w = httptest.NewRecorder()
	fakeRepa.EXPECT().GetAllPostsGyCategory("music").Return(&[]post.Post{}, nil)
	service.ShowAllCategoryPosts(w, req)
	resp = w.Result()
	if resp.StatusCode != 200 {
		t.Errorf("expected resp status 200, got %d", resp.StatusCode)
		return
	}
}

func TestPostHandlerShowAllUserPosts(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	fakeRepa := NewMockPostRepositoryInterface(ctrl)
	service := &PostHandler{PostRepo: fakeRepa}

	req := httptest.NewRequest("GET", "/api/user/{user_login}", nil)
	req = mux.SetURLVars(req, map[string]string{"user_login": "Nikita"})
	w := httptest.NewRecorder()
	fakeRepa.EXPECT().GetAllPostsGyLogin("Nikita").Return(&[]post.Post{}, fmt.Errorf("didn't get"))
	service.ShowAllUserPosts(w, req)
	resp := w.Result()
	if resp.StatusCode != 500 {
		t.Errorf("expected resp status 500, got %d", resp.StatusCode)
		return
	}

	req = httptest.NewRequest("GET", "/api/user/{user_login}", nil)
	req = mux.SetURLVars(req, map[string]string{"user_login": "Nikita"})
	w = httptest.NewRecorder()
	fakeRepa.EXPECT().GetAllPostsGyLogin("Nikita").Return(&[]post.Post{}, nil)
	service.ShowAllUserPosts(w, req)
	resp = w.Result()
	if resp.StatusCode != 200 {
		t.Errorf("expected resp status 200, got %d", resp.StatusCode)
		return
	}
}
