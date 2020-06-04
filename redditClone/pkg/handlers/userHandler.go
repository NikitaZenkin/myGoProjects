package handlers

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"net/http"
	"rclone/pkg/system"
	"rclone/pkg/user"
)

type UserRepositoryInterface interface {
	UserExist(userName string) (*user.User, error)
	AddUser(newUser *user.User) error
}

type SessionsManagerInterface interface {
	CreateSession(w http.ResponseWriter, userID string, userName string) error
}

type UserHandler struct {
	Tmpl     *template.Template
	UserRepo UserRepositoryInterface
	Sessions SessionsManagerInterface
}

type PreUser struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

func (handler *UserHandler) Template(w http.ResponseWriter, r *http.Request) {
	handler.Tmpl.ExecuteTemplate(w, "index.html", struct{}{})
}

func (handler *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	preUser := PreUser{}
	err := json.Unmarshal(body, &preUser)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	curUser, err := handler.UserRepo.UserExist(preUser.UserName)
	if err != nil {
		mes := system.Message{Message: "user not found"}
		jsonMes := mes.ToJson()
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(jsonMes)
		return
	}
	if !curUser.CheckPassword(preUser.Password) {
		mes := system.Message{Message: "invalid password"}
		jsonMes := mes.ToJson()
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(jsonMes)
		return
	}
	jsonToken := curUser.JsonToken()
	err = handler.Sessions.CreateSession(w, curUser.ID, curUser.UserName)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(jsonToken)
}

func (handler *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	preUser := PreUser{}
	err := json.Unmarshal(body, &preUser)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	curUser, err := handler.UserRepo.UserExist(preUser.UserName)
	if err == nil {
		newError := system.Error{
			Location: "body",
			Param:    "username",
			Value:    curUser.UserName,
			Message:  "already exists",
		}
		errors := system.Errors{Errors: []system.Error{newError}}
		jsonMes := errors.ToJson()
		w.WriteHeader(422)
		w.Write(jsonMes)
		return
	}
	newUser := user.NewUser(preUser.UserName, preUser.Password)
	jsonToken := newUser.JsonToken()
	err = handler.UserRepo.AddUser(&newUser)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = handler.Sessions.CreateSession(w, newUser.ID, newUser.UserName)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(jsonToken)
}
