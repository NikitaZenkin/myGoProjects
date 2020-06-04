package system

import (
	"encoding/json"
	"math/rand"
	"strconv"
)

type Message struct {
	Message string `json:"message"`
}

func (msg Message) ToJson() []byte {
	jsonByte, _ := json.Marshal(msg)
	return jsonByte
}

type Error struct {
	Location string `json:"location"`
	Param    string `json:"param"`
	Value    string `json:"value"`
	Message  string `json:"msg"`
}

type Errors struct {
	Errors []Error `json:"errors"`
}

func (errors Errors) ToJson() []byte {
	jsonByte, _ := json.Marshal(errors)
	return jsonByte
}

func NewId() string {
	return strconv.Itoa(rand.Intn(MaxId))
}
