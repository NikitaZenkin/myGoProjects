package user

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"rclone/pkg/system"
	"time"
)

type TokenStruct struct {
	Token string `json:"token"`
}

type User struct {
	ID       string
	UserName string
	password string
}

func NewUser(userName string, password string) User {

	hashByte := md5.Sum([]byte(password))
	return User{
		ID:       system.NewId(),
		UserName: userName,
		password: fmt.Sprintf("%x", hashByte),
	}
}

func (user *User) CheckPassword(password string) bool {
	hashByte := md5.Sum([]byte(password))
	return user.password == fmt.Sprintf("%x", hashByte)
}

func (user User) JsonToken() []byte {
	userFields := jwt.MapClaims{
		"username": user.UserName,
		"id":       user.ID,
	}
	time := time.Now()
	endTime := time.AddDate(0, 1, 0)
	preToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp":  endTime.Unix(),
		"iat":  time.Unix(),
		"user": userFields,
	})
	tokenValue, _ := preToken.SignedString([]byte("key"))
	token := TokenStruct{Token: tokenValue}
	tokenByte, _ := json.Marshal(token)
	return tokenByte
}
