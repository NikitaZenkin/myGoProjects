package main

import (
	tgbotapi "gopkg.in/telegram-bot-api.v4"
	"net/http"
	"os"
	"strings"
)

//в этом файле происходит инциализация игрового мира и запускается игра

//глобальная переменная - игрок
var player person

//функция инциализации игрового мира
func initGame() {
	//создаем кухню
	kitchen := area{name: "кухня", open: true}
	func() {
		table := furniture{name: "на столе"}
		func() {
			tea := item{name: "чай"}
			table.itemsHere = append(table.itemsHere, &tea)
		}()
		kitchen.furnitureHere = append(kitchen.furnitureHere, &table)
		kitchen.intro = func() string {
			return "кухня, " + kitchen.allInterestingItems() + ". " + player.placesToGo()
		}
		kitchen.info = func() string {
			var stringToReturn []string
			stringToReturn = append(stringToReturn, "ты находишься на кухне, ")
			stringToReturn = append(stringToReturn, kitchen.allItems())
			if exist, _ := player.hasItems("конспекты"); exist {
				stringToReturn = append(stringToReturn, ", надо идти в универ. ")
			} else {
				stringToReturn = append(stringToReturn, ", надо собрать рюкзак и идти в универ. ")
			}
			stringToReturn = append(stringToReturn, player.placesToGo())
			return strings.Join(stringToReturn, "")
		}
	}()

	//создаем коридор
	hole := area{name: "коридор", open: true}
	func() {
		hole.intro = func() string {
			return hole.allInterestingItems() + ". " + player.placesToGo()
		}
		hole.info = func() string {
			return hole.allInterestingItems() + ". " + player.placesToGo()
		}
	}()

	//создаем комнату
	bedRoom := area{name: "комната", open: true}
	func() {
		tableInRoom := furniture{name: "на столе"}
		func() {
			key := item{name: "ключи", takeable: true, usable: true}
			tableInRoom.itemsHere = append(tableInRoom.itemsHere, &key)
			conspects := item{name: "конспекты", takeable: true}
			tableInRoom.itemsHere = append(tableInRoom.itemsHere, &conspects)
		}()
		bedRoom.furnitureHere = append(bedRoom.furnitureHere, &tableInRoom)

		chair := furniture{name: "на стуле"}
		func() {
			bag := item{name: "рюкзак", wearable: true}
			chair.itemsHere = append(chair.itemsHere, &bag)
		}()
		bedRoom.furnitureHere = append(bedRoom.furnitureHere, &chair)
		bedRoom.intro = func() string {
			return "ты в своей комнате. " + player.placesToGo()
		}
		bedRoom.info = func() string {
			return bedRoom.allItems() + ". " + player.placesToGo()
		}
	}()

	//создаем улицу
	street := area{name: "улица", open: false, outdoor: true}
	func() {
		street.intro = func() string {
			return "на улице весна. " + player.placesToGo()
		}
		street.info = func() string {
			return street.allInterestingItems() + ". " + player.placesToGo()
		}
	}()

	//создаем граф игрового мира
	func() {
		kitchen.areaToGo = append(kitchen.areaToGo, &hole)
		hole.areaToGo = append(hole.areaToGo, &kitchen)
		hole.areaToGo = append(hole.areaToGo, &bedRoom)
		hole.areaToGo = append(hole.areaToGo, &street)
		bedRoom.areaToGo = append(bedRoom.areaToGo, &hole)
		street.areaToGo = append(street.areaToGo, &hole)
	}()

	//начальное сотсояние игрока
	player.place = &kitchen
	player.hasBag = false
	player.items = nil
}

//функция принимающая очередную команду
func handleCommand(command string) string {
	listOfCommands := strings.Split(command, " ")
	switch listOfCommands[0] {
	case "идти":
		return goTo(listOfCommands)
	case "осмотреться":
		return lookAround(listOfCommands)
	case "надеть":
		return wear(listOfCommands)
	case "взять":
		return take(listOfCommands)
	case "применить":
		return use(listOfCommands)
	default:
		return "неизвестная команда"
	}
}

var (
	BotToken   = "1029459525:AAHrOvge3ELvKfloM4XUOnrRrZZADJUCtlY"
	WebhookURL = "https://zenkin-game-bot.herokuapp.com/"
)

func main() {
	bot, err := tgbotapi.NewBotAPI(BotToken)
	if err != nil {
		panic(err)
	}
	bot.Debug = true
	_, err = bot.SetWebhook(tgbotapi.NewWebhook(WebhookURL))
	if err != nil {
		panic(err)
	}
	updates := bot.ListenForWebhook("/")
	port := os.Getenv("PORT")
	go http.ListenAndServe(":"+port, nil)
	initGame()
	for update := range updates {
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, handleCommand(update.Message.Text)))
	}
}
