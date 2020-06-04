package handlers

import (
	"context"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"net/http/httptest"
	"rclone/pkg/session"
	"testing"
)

func TestPostHandlerDownvote(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	fakeRepa := NewMockPostRepositoryInterface(ctrl)
	service := &PostHandler{PostRepo: fakeRepa}

	req := httptest.NewRequest("GET", "/api/post/{post_id}/downvote", nil)
	w := httptest.NewRecorder()
	service.Downvote(w, req)
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
	req = httptest.NewRequest("GET", "/api/post/{post_id}/downvote", nil)
	req = mux.SetURLVars(req, map[string]string{"post_id": "345"})
	w = httptest.NewRecorder()
	ctx := context.WithValue(req.Context(), "sessionKey", &sess)
	fakeRepa.EXPECT().DownvotePost("345", sess.UserId).Return(nil, fmt.Errorf("post not exist"))
	service.Downvote(w, req.WithContext(ctx))
	resp = w.Result()
	if resp.StatusCode != 500 {
		t.Errorf("expected resp status 500, got %d", resp.StatusCode)
		return
	}

	req = httptest.NewRequest("GET", "/api/post/{post_id}/downvote", nil)
	req = mux.SetURLVars(req, map[string]string{"post_id": "345"})
	w = httptest.NewRecorder()
	ctx = context.WithValue(req.Context(), "sessionKey", &sess)
	fakeRepa.EXPECT().DownvotePost("345", sess.UserId).Return(nil, nil)
	service.Downvote(w, req.WithContext(ctx))
	resp = w.Result()
	if resp.StatusCode != 200 {
		t.Errorf("expected resp status 200, got %d", resp.StatusCode)
		return
	}
}

func TestPostHandlerUpvote(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	fakeRepa := NewMockPostRepositoryInterface(ctrl)
	service := &PostHandler{PostRepo: fakeRepa}

	req := httptest.NewRequest("GET", "/api/post/{post_id}/upvote", nil)
	w := httptest.NewRecorder()
	service.Upvote(w, req)
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
	req = httptest.NewRequest("GET", "/api/post/{post_id}/upvote", nil)
	req = mux.SetURLVars(req, map[string]string{"post_id": "346"})
	w = httptest.NewRecorder()
	ctx := context.WithValue(req.Context(), "sessionKey", &sess)
	fakeRepa.EXPECT().UpvotePost("346", sess.UserId).Return(nil, fmt.Errorf("post not exist"))
	service.Upvote(w, req.WithContext(ctx))
	resp = w.Result()
	if resp.StatusCode != 500 {
		t.Errorf("expected resp status 500, got %d", resp.StatusCode)
		return
	}

	req = httptest.NewRequest("GET", "/api/post/{post_id}/upvote", nil)
	req = mux.SetURLVars(req, map[string]string{"post_id": "346"})
	w = httptest.NewRecorder()
	ctx = context.WithValue(req.Context(), "sessionKey", &sess)
	fakeRepa.EXPECT().UpvotePost("346", sess.UserId).Return(nil, nil)
	service.Upvote(w, req.WithContext(ctx))
	resp = w.Result()
	if resp.StatusCode != 200 {
		t.Errorf("expected resp status 200, got %d", resp.StatusCode)
		return
	}
}

func TestPostHandlerUnvote(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	fakeRepa := NewMockPostRepositoryInterface(ctrl)
	service := &PostHandler{PostRepo: fakeRepa}

	req := httptest.NewRequest("GET", "/api/post/{post_id}/unvote", nil)
	w := httptest.NewRecorder()
	service.Unvote(w, req)
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
	req = httptest.NewRequest("GET", "/api/post/{post_id}/unvote", nil)
	req = mux.SetURLVars(req, map[string]string{"post_id": "346"})
	w = httptest.NewRecorder()
	ctx := context.WithValue(req.Context(), "sessionKey", &sess)
	fakeRepa.EXPECT().UnvotePost("346", sess.UserId).Return(nil, fmt.Errorf("post not exist"))
	service.Unvote(w, req.WithContext(ctx))
	resp = w.Result()
	if resp.StatusCode != 500 {
		t.Errorf("expected resp status 500, got %d", resp.StatusCode)
		return
	}

	req = httptest.NewRequest("GET", "/api/post/{post_id}/unvote", nil)
	req = mux.SetURLVars(req, map[string]string{"post_id": "346"})
	w = httptest.NewRecorder()
	ctx = context.WithValue(req.Context(), "sessionKey", &sess)
	fakeRepa.EXPECT().UnvotePost("346", sess.UserId).Return(nil, nil)
	service.Unvote(w, req.WithContext(ctx))
	resp = w.Result()
	if resp.StatusCode != 200 {
		t.Errorf("expected resp status 200, got %d", resp.StatusCode)
		return
	}
}
