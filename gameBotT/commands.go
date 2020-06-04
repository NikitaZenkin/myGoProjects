package main

//здесь описываются все команды доступные игроку

//команда идти
func goTo(commands []string) string {
	if len(commands) < 2 {
		return "нехватает аргументов команды"
	}
	if player.place.outdoor && commands[1] == "домой" {
		commands[1] = player.place.areaToGo[0].name
	}
	for _, newPlace := range player.place.areaToGo {
		if commands[1] == newPlace.name {
			if !newPlace.open {
				return "дверь закрыта"
			} else {
				player.place = newPlace
				return player.place.intro()
			}
		}
	}
	return "нет пути в " + commands[1]
}

//команда осмотреться
func lookAround(commands []string) string {
	return player.place.info()
}

//команда надеть
func wear(commands []string) string {
	if len(commands) < 2 {
		return "нехватает аргументов команды"
	}
	return player.wear(commands[1])
}

//команда взять
func take(commands []string) string {
	if len(commands) < 2 {
		return "нехватает аргументов команды"
	}
	return player.take(commands[1])
}

//команда применить
func use(commands []string) string {
	if len(commands) < 3 {
		return "нехватает аргументов команды"
	}
	exist, currentItem := player.hasItems(commands[1])
	if !exist {
		return "нет предмета в инвентаре - " + commands[1]
	}
	if currentItem.usable {
		switch currentItem.name {
		case "ключи":
			{
				switch commands[2] {
				case "дверь":
					player.place.openDoors()
					return "дверь открыта"
				default:
					return "не к чему применить"
				}
			}
		}
	}
	return "этот предмет нельзя использовать"
}
