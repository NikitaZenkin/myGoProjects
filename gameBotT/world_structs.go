package main

import "strings"

//В этом файле описываются струтуры и методы игровых элементов (помещения, мебель, предметы)

//структура помещения
type area struct {
	name          string
	intro         func() string
	info          func() string
	areaToGo      []*area
	furnitureHere []*furniture
	open          bool
	outdoor       bool
}

//струтура мебели
type furniture struct {
	name      string
	itemsHere []*item
}

//струтура предмета
type item struct {
	name     string
	takeable bool
	wearable bool
	usable   bool
}

//вывод списка всех предсетов в данном помещении
func (place area) allItems() string {
	var stringToReturn []string
	for _, nextFurniture := range place.furnitureHere {
		if len(nextFurniture.itemsHere) != 0 {
			stringToReturn = append(stringToReturn, nextFurniture.name+": ")
			for _, nextItem := range nextFurniture.itemsHere {
				stringToReturn = append(stringToReturn, nextItem.name)
				stringToReturn = append(stringToReturn, ", ")
			}
		}
	}
	if stringToReturn == nil {
		return "пустая " + place.name
	}
	return strings.Join(stringToReturn[:len(stringToReturn)-1], "")
}

//вывод списка всех предсетов в данном помещении, которые можно взять или надеть
func (place area) allInterestingItems() string {
	var stringToReturn []string
	for _, nextFurniture := range place.furnitureHere {
		if len(nextFurniture.itemsHere) != 0 && nextFurniture.hasInterestingItems() {
			stringToReturn = append(stringToReturn, nextFurniture.name+": ")
			for _, item := range nextFurniture.itemsHere {
				if item.wearable || item.takeable {
					stringToReturn = append(stringToReturn, item.name)
					stringToReturn = append(stringToReturn, ", ")
				}
			}
		}
	}
	if stringToReturn == nil {
		return "ничего интересного"
	}
	return strings.Join(stringToReturn[:len(stringToReturn)-1], "")
}

//функция открытя дверей в помещении
func (place *area) openDoors() {
	for _, nextPlace := range place.areaToGo {
		nextPlace.open = true
	}
}

//функция проверяющая, что у данной мебели есть предметы которые можно взять или надеть
func (furn furniture) hasInterestingItems() bool {
	flag := false
	for _, nextItem := range furn.itemsHere {
		flag = flag || nextItem.takeable
		flag = flag || nextItem.wearable
	}
	return flag
}
