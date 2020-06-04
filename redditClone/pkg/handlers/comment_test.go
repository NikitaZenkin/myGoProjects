package handlers

import (
	"bytes"
	"context"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"net/http/httptest"
	"rclone/pkg/session"
	"testing"
)

func TestPostHandlerNewComment(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	fakeRepa := NewMockPostRepositoryInterface(ctrl)
	service := &PostHandler{PostRepo: fakeRepa}

	req := httptest.NewRequest("POST", "/api/post/{post_id}", nil)
	w := httptest.NewRecorder()
	service.NewComment(w, req)
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
	req = httptest.NewRequest("POST", "/api/post/{post_id}", nil)
	w = httptest.NewRecorder()
	req = mux.SetURLVars(req, map[string]string{"post_id": "346"})
	ctx := context.WithValue(req.Context(), "sessionKey", &sess)
	service.NewComment(w, req.WithContext(ctx))
	resp = w.Result()
	if resp.StatusCode != 500 {
		t.Errorf("expected resp status 500, got %d", resp.StatusCode)
		return
	}

	rqBody := bytes.Buffer{}
	rqBody.Write([]byte(`{"comment": "help"}`))
	req = httptest.NewRequest("POST", "/api/post/{post_id}", &rqBody)
	w = httptest.NewRecorder()
	req = mux.SetURLVars(req, map[string]string{"post_id": "346"})
	ctx = context.WithValue(req.Context(), "sessionKey", &sess)
	fakeRepa.EXPECT().AddCommentToPost("346", gomock.Any()).Return(nil, fmt.Errorf("didn't creae comment"))
	service.NewComment(w, req.WithContext(ctx))
	resp = w.Result()
	if resp.StatusCode != 500 {
		t.Errorf("expected resp status 500, got %d", resp.StatusCode)
		return
	}

	rqBody.Write([]byte(`{"comment": "help"}`))
	req = httptest.NewRequest("POST", "/api/post/{post_id}", &rqBody)
	w = httptest.NewRecorder()
	req = mux.SetURLVars(req, map[string]string{"post_id": "346"})
	ctx = context.WithValue(req.Context(), "sessionKey", &sess)
	fakeRepa.EXPECT().AddCommentToPost("346", gomock.Any()).Return(nil, nil)
	service.NewComment(w, req.WithContext(ctx))
	resp = w.Result()
	if resp.StatusCode != 201 {
		t.Errorf("expected resp status 201, got %d", resp.StatusCode)
		return
	}
}

func TestPostHandlerDeleteComment(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	fakeRepa := NewMockPostRepositoryInterface(ctrl)
	service := &PostHandler{PostRepo: fakeRepa}

	req := httptest.NewRequest("DELETE", "/api/post/{post_id}/{comment_id}", nil)
	req = mux.SetURLVars(req, map[string]string{"post_id": "346", "comment_id": "725"})
	w := httptest.NewRecorder()
	fakeRepa.EXPECT().DeleteCommentFromPost("346", "725").Return(nil, fmt.Errorf("didn't delete"))
	service.DeleteComment(w, req)
	resp := w.Result()
	if resp.StatusCode != 500 {
		t.Errorf("expected resp status 500, got %d", resp.StatusCode)
		return
	}

	req = httptest.NewRequest("DELETE", "/api/post/{post_id}/{comment_id}", nil)
	req = mux.SetURLVars(req, map[string]string{"post_id": "346", "comment_id": "725"})
	w = httptest.NewRecorder()
	fakeRepa.EXPECT().DeleteCommentFromPost("346", "725").Return(nil, nil)
	service.DeleteComment(w, req)
	resp = w.Result()
	if resp.StatusCode != 200 {
		t.Errorf("expected resp status 200, got %d", resp.StatusCode)
		return
	}

}
