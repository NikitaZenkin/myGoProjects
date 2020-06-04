package main

import (
	"fmt"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
	"strconv"
	"strings"
)

func checkUser(update tgbotapi.Update) {
	if _, exist := userMap[update.Message.Chat.ID]; !exist {
		newUser := User{
			name: "@" + update.Message.Chat.UserName,
			id:   update.Message.Chat.ID,
		}
		userMap[update.Message.Chat.ID] = newUser
	}
}

func handleCommand(update tgbotapi.Update) (map[int64]string, error) {
	listOfCommands := strings.Split(update.Message.Text, " ")
	curId := update.Message.Chat.ID
	extraList := strings.Split(listOfCommands[0], "_")
	switch extraList[0] {
	case "/tasks":
		return showTasks(curId), nil
	case "/new":
		return newTask(listOfCommands, curId), nil
	case "/assign":
		return assignTasks(listOfCommands, curId), nil
	case "/unassign":
		return unassignTasks(listOfCommands, curId), nil
	case "/resolve":
		return resolveTasks(listOfCommands, curId), nil
	case "/my":
		return showMyTasks(curId), nil
	case "/owner":
		return showTasksMadeByMe(curId), nil
	default:
		return map[int64]string{}, fmt.Errorf("Недопустимая команда")
	}
}

func showTasks(currentUserId int64) map[int64]string {
	mapOfResponsesToUsers := make(map[int64]string)
	if len(taskMap) == 0 {
		mapOfResponsesToUsers[currentUserId] = "Нет задач"
		return mapOfResponsesToUsers
	}
	responseString := make([]string, len(taskMap))
	var i int
	for _, task := range taskMap {
		responseString[i] = task.infoForUser(currentUserId, true)
		i++
	}
	mapOfResponsesToUsers[currentUserId] = strings.Join(responseString, "\n\n")
	return mapOfResponsesToUsers
}

func newTask(commands []string, currentUserId int64) map[int64]string {
	mapOfResponsesToUsers := make(map[int64]string)
	if len(commands) == 1 {
		mapOfResponsesToUsers[currentUserId] = "Введите название таски"
		return mapOfResponsesToUsers
	}
	taskNum++
	task := Task{
		number:    taskNum,
		name:      strings.Join(commands[1:], " "),
		creatorId: currentUserId,
	}
	taskMap[taskNum] = task
	mapOfResponsesToUsers[currentUserId] = task.toShortString() + "создана, id=" + strconv.Itoa(int(task.number))
	return mapOfResponsesToUsers
}

func assignTasks(commands []string, currentUserId int64) map[int64]string {
	mapOfResponsesToUsers := make(map[int64]string)
	list := strings.Split(commands[0], "_")
	numTask, err := strconv.Atoi(list[1])
	if err != nil {
		mapOfResponsesToUsers[currentUserId] = "Нужно указать номер таски"
		return mapOfResponsesToUsers
	}
	task, ex := taskMap[int64(numTask)]
	if !ex {
		mapOfResponsesToUsers[currentUserId] = "Такой таски не существет"
		return mapOfResponsesToUsers
	}
	if task.assignId == currentUserId {
		mapOfResponsesToUsers[currentUserId] = "Таска уже на вас"
		return mapOfResponsesToUsers
	}
	mapOfResponsesToUsers[currentUserId] = task.toShortString() + "назначена на вас"
	if task.assignId == 0 && task.creatorId != currentUserId {
		mapOfResponsesToUsers[task.creatorId] = task.toShortString() + "назначена на " + userMap[currentUserId].name
	} else if task.assignId != 0 {
		mapOfResponsesToUsers[task.assignId] = task.toShortString() + "назначена на " + userMap[currentUserId].name
	}
	taskMap[int64(numTask)] = Task{
		number:    task.number,
		name:      task.name,
		creatorId: task.creatorId,
		assignId:  currentUserId,
	}
	return mapOfResponsesToUsers
}

func unassignTasks(commands []string, currentUserId int64) map[int64]string {
	mapOfResponsesToUsers := make(map[int64]string)
	list := strings.Split(commands[0], "_")
	numTask, err := strconv.Atoi(list[1])
	if err != nil {
		mapOfResponsesToUsers[currentUserId] = "Нужно указать номер таски"
		return mapOfResponsesToUsers
	}
	task, ex := taskMap[int64(numTask)]
	if !ex {
		mapOfResponsesToUsers[currentUserId] = "Такой таски не существет"
		return mapOfResponsesToUsers
	}
	if task.assignId != currentUserId {
		mapOfResponsesToUsers[currentUserId] = "Задача не на вас"
		return mapOfResponsesToUsers
	}
	mapOfResponsesToUsers[currentUserId] = "Принято"
	mapOfResponsesToUsers[task.creatorId] = task.toShortString() + "осталась без исполнителя"
	taskMap[int64(numTask)] = Task{
		number:    task.number,
		name:      task.name,
		creatorId: task.creatorId,
		assignId:  0,
	}
	return mapOfResponsesToUsers
}

func resolveTasks(commands []string, currentUserId int64) map[int64]string {
	mapOfResponsesToUsers := make(map[int64]string)
	list := strings.Split(commands[0], "_")
	numTask, err := strconv.Atoi(list[1])
	if err != nil {
		mapOfResponsesToUsers[currentUserId] = "Нужно указать номер таски"
		return mapOfResponsesToUsers
	}
	task, ex := taskMap[int64(numTask)]
	if !ex {
		mapOfResponsesToUsers[currentUserId] = "Такой таски не существет"
		return mapOfResponsesToUsers
	}
	if task.assignId != currentUserId {
		mapOfResponsesToUsers[currentUserId] = "Задача не на вас"
		return mapOfResponsesToUsers
	}
	mapOfResponsesToUsers[currentUserId] = task.toShortString() + "выполнена"
	if currentUserId != task.creatorId {
		mapOfResponsesToUsers[task.creatorId] = task.toShortString() + "выполнена " + userMap[currentUserId].name
	}
	delete(taskMap, task.number)
	return mapOfResponsesToUsers
}

func showTasksMadeByMe(currentUserId int64) map[int64]string {
	mapOfResponsesToUsers := make(map[int64]string)
	var responseString []string
	for _, task := range taskMap {
		if task.creatorId == currentUserId {
			responseString = append(responseString, task.infoForUser(currentUserId, true))
		}
	}
	if len(responseString) == 0 {
		mapOfResponsesToUsers[currentUserId] = "Нет созданных вами задач"
		return mapOfResponsesToUsers
	}
	mapOfResponsesToUsers[currentUserId] = strings.Join(responseString, "\n")
	return mapOfResponsesToUsers
}

func showMyTasks(currentUserId int64) map[int64]string {
	mapOfResponsesToUsers := make(map[int64]string)
	var responseString []string
	for _, task := range taskMap {
		if task.assignId == currentUserId {
			responseString = append(responseString, task.infoForUser(currentUserId, false))
		}
	}
	if len(responseString) == 0 {
		mapOfResponsesToUsers[currentUserId] = "Нет задач назанченных на вас"
		return mapOfResponsesToUsers
	}
	mapOfResponsesToUsers[currentUserId] = strings.Join(responseString, "\n")
	return mapOfResponsesToUsers
}
