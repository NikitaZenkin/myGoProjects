package handlers

import (
	"bytes"
	"fmt"
	"github.com/golang/mock/gomock"
	"html/template"
	"net/http/httptest"
	"rclone/pkg/user"
	"testing"
)

func TestUserHandlerRegister(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	fakeRepa := NewMockUserRepositoryInterface(ctrl)
	fakeManager := NewMockSessionsManagerInterface(ctrl)
	service := &UserHandler{
		UserRepo: fakeRepa,
		Sessions: fakeManager,
		Tmpl:     template.Must(template.ParseFiles("./template_for_test/index.html")),
	}
	rqBody := bytes.Buffer{}

	rqBody.Write([]byte(`bad json`))
	req := httptest.NewRequest("POST", "/api/register", &rqBody)
	w := httptest.NewRecorder()
	service.Register(w, req)
	resp := w.Result()
	if resp.StatusCode != 500 {
		t.Errorf("expected resp status 500, got %d", resp.StatusCode)
		return
	}

	rqBody.Write([]byte(`{"username": "Nikita", "password": "nope"}`))
	userc := user.NewUser("Nikita", "happy")
	fakeRepa.EXPECT().UserExist("Nikita").Return(&userc, nil)
	req = httptest.NewRequest("POST", "/api/register", &rqBody)
	w = httptest.NewRecorder()
	service.Register(w, req)
	resp = w.Result()
	if resp.StatusCode != 422 {
		t.Errorf("expected resp status 422, got %d", resp.StatusCode)
		return
	}

	rqBody.Write([]byte(`{"username": "Nikita", "password": "nope"}`))
	userc = user.NewUser("Nikita", "nope")
	fakeRepa.EXPECT().UserExist("Nikita").Return(&userc, fmt.Errorf("noUser"))
	fakeRepa.EXPECT().AddUser(gomock.Any()).Return(fmt.Errorf("can't Add"))
	req = httptest.NewRequest("POST", "/api/register", &rqBody)
	w = httptest.NewRecorder()
	service.Register(w, req)
	resp = w.Result()
	if resp.StatusCode != 500 {
		t.Errorf("expected resp status 500, got %d", resp.StatusCode)
		return
	}

	rqBody.Write([]byte(`{"username": "Nikita", "password": "nope"}`))
	userc = user.NewUser("Nikita", "nope")
	req = httptest.NewRequest("POST", "/api/register", &rqBody)
	w = httptest.NewRecorder()
	fakeRepa.EXPECT().UserExist("Nikita").Return(&userc, fmt.Errorf("noUser"))
	fakeRepa.EXPECT().AddUser(gomock.Any()).Return(nil)
	fakeManager.EXPECT().CreateSession(w, gomock.Any(), userc.UserName).Return(fmt.Errorf("can't Create"))
	service.Register(w, req)
	resp = w.Result()
	if resp.StatusCode != 500 {
		t.Errorf("expected resp status 500, got %d", resp.StatusCode)
		return
	}

	rqBody.Write([]byte(`{"username": "Nikita", "password": "nope"}`))
	userc = user.NewUser("Nikita", "nope")
	req = httptest.NewRequest("POST", "/api/register", &rqBody)
	w = httptest.NewRecorder()
	fakeRepa.EXPECT().UserExist("Nikita").Return(&userc, fmt.Errorf("noUser"))
	fakeRepa.EXPECT().AddUser(gomock.Any()).Return(nil)
	fakeManager.EXPECT().CreateSession(w, gomock.Any(), userc.UserName).Return(nil)
	service.Register(w, req)
	resp = w.Result()
	if resp.StatusCode != 200 {
		t.Errorf("expected resp status 200, got %d", resp.StatusCode)
		return
	}

}

func TestUserHandlerLogin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	fakeRepa := NewMockUserRepositoryInterface(ctrl)
	fakeManager := NewMockSessionsManagerInterface(ctrl)
	service := &UserHandler{
		UserRepo: fakeRepa,
		Sessions: fakeManager,
		Tmpl:     template.Must(template.ParseFiles("./template_for_test/index.html")),
	}
	rqBody := bytes.Buffer{}

	rqBody.Write([]byte(`bad json`))
	req := httptest.NewRequest("POST", "/api/loginn", &rqBody)
	w := httptest.NewRecorder()
	service.Login(w, req)
	resp := w.Result()
	if resp.StatusCode != 500 {
		t.Errorf("expected resp status 500, got %d", resp.StatusCode)
		return
	}

	rqBody.Write([]byte(`{"username": "Nikita", "password": "nope"}`))
	fakeRepa.EXPECT().UserExist("Nikita").Return(nil, fmt.Errorf("noUser"))
	req = httptest.NewRequest("POST", "/api/login", &rqBody)
	w = httptest.NewRecorder()
	service.Login(w, req)
	resp = w.Result()
	if resp.StatusCode != 401 {
		t.Errorf("expected resp status 401, got %d", resp.StatusCode)
		return
	}

	rqBody.Write([]byte(`{"username": "Nikita", "password": "nope"}`))
	req = httptest.NewRequest("POST", "/api/login", &rqBody)
	w = httptest.NewRecorder()
	fakeRepa.EXPECT().UserExist("Nikita").Return(&user.User{}, nil)
	service.Login(w, req)
	resp = w.Result()
	if resp.StatusCode != 401 {
		t.Errorf("expected resp status 401, got %d", resp.StatusCode)
		return
	}

	rqBody.Write([]byte(`{"username": "Nikita", "password": "happy"}`))
	userc := user.NewUser("Nikita", "happy")
	req = httptest.NewRequest("POST", "/api/login", &rqBody)
	w = httptest.NewRecorder()
	fakeRepa.EXPECT().UserExist("Nikita").Return(&userc, nil)
	fakeManager.EXPECT().CreateSession(w, userc.ID, userc.UserName).Return(fmt.Errorf("notCreatedr"))
	service.Login(w, req)
	resp = w.Result()
	if resp.StatusCode != 500 {
		t.Errorf("expected resp status 401, got %d", resp.StatusCode)
		return
	}

	rqBody.Write([]byte(`{"username": "Nikita", "password": "happy"}`))
	userc = user.NewUser("Nikita", "happy")
	req = httptest.NewRequest("POST", "/api/login", &rqBody)
	w = httptest.NewRecorder()
	fakeRepa.EXPECT().UserExist("Nikita").Return(&userc, nil)
	fakeManager.EXPECT().CreateSession(w, userc.ID, userc.UserName).Return(nil)
	service.Login(w, req)
	resp = w.Result()
	if resp.StatusCode != 200 {
		t.Errorf("expected resp status 200, got %d", resp.StatusCode)
		return
	}

}

func TestUserHandlerTemplate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	fakeRepa := NewMockUserRepositoryInterface(ctrl)
	fakeManager := NewMockSessionsManagerInterface(ctrl)
	service := &UserHandler{
		UserRepo: fakeRepa,
		Sessions: fakeManager,
		Tmpl:     template.Must(template.ParseFiles("./template_for_test/index.html")),
	}

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	service.Template(w, req)
}
