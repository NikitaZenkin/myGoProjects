package main

import "strings"

//В этом файле описываются струтура и методы персонажа

//структура песонажа
type person struct {
	hasBag bool
	place  *area
	items  []*item
}

//функция проверки наличия предмета у персонажа
func (per person) hasItems(itemName string) (bool, *item) {
	for _, nextItem := range per.items {
		if itemName == nextItem.name {
			return true, nextItem
		}
	}
	return false, nil
}

//функция возвращает список мест куда сейчас можно пройти
func (per person) placesToGo() string {
	var stringToReturn []string
	stringToReturn = append(stringToReturn, "можно пройти - ")
	if per.place.outdoor {
		stringToReturn = append(stringToReturn, "домой")
		stringToReturn = append(stringToReturn, ", ")
	} else {
		for _, nextPlace := range per.place.areaToGo {
			stringToReturn = append(stringToReturn, nextPlace.name)
			stringToReturn = append(stringToReturn, ", ")
		}
	}
	return strings.Join(stringToReturn[:len(stringToReturn)-1], "")
}

//функция для взятия предметов
func (per *person) take(itemName string) string {
	if !per.hasBag {
		return "некуда класть"
	}
	for _, nextFurniture := range per.place.furnitureHere {
		for i, nextItem := range nextFurniture.itemsHere {
			if nextItem.name == itemName {
				if !nextItem.takeable {
					return "нельзя взять: " + itemName
				} else {
					per.items = append(per.items, nextItem)
					nextFurniture.itemsHere = append(nextFurniture.itemsHere[:i], nextFurniture.itemsHere[i+1:]...)
					return "предмет добавлен в инвентарь: " + itemName
				}
			}
		}
	}
	return "нет такого"
}

//функция для надевания предметов
func (per *person) wear(itemName string) string {
	for _, nextFurniture := range per.place.furnitureHere {
		for i, nextItem := range nextFurniture.itemsHere {
			if nextItem.name == itemName {
				if !nextItem.wearable {
					return "нельзя надеть: " + itemName
				} else {
					per.items = append(per.items, nextItem)
					nextFurniture.itemsHere = append(nextFurniture.itemsHere[:i], nextFurniture.itemsHere[i+1:]...)
					if itemName == "рюкзак" {
						per.hasBag = true
					}
					return "вы надели: " + itemName
				}
			}
		}
	}
	return "нет такого"
}
