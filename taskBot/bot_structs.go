package main

import "strconv"

var userMap = map[int64]User{}
var taskMap = map[int64]Task{}
var taskNum int64

type User struct {
	name string
	id   int64
}

type Task struct {
	number    int64
	name      string
	creatorId int64
	assignId  int64
}

func (task Task) toString() string {
	return strconv.Itoa(int(task.number)) + ". " + task.name + " by " + userMap[task.creatorId].name
}

func (task Task) toShortString() string {
	return "Задача \"" + task.name + "\" "
}

func (task Task) userOptionsToString(userId int64) string {
	if task.assignId == 0 {
		return "\n/assign_" + strconv.Itoa(int(task.number))
	}
	if task.assignId == userId {
		return "\n/unassign_" + strconv.Itoa(int(task.number)) + " /resolve_" + strconv.Itoa(int(task.number))
	} else {
		return "\nassignee: " + userMap[task.assignId].name
	}
}

func (task Task) infoForUser(userId int64, withAssigneeOnMe bool) string {
	if withAssigneeOnMe && task.assignId == userId {
		return task.toString() + "\nassignee: я" + task.userOptionsToString(userId)
	} else {
		return task.toString() + task.userOptionsToString(userId)
	}
}
