package main

// сюда писать код

import (
	"context"
	"fmt"
	"gopkg.in/telegram-bot-api.v4"
	"net/http"
	"os"
)

var (
	// @BotFather в телеграме даст вам это
	BotToken = "1046936324:AAENfOSfpAE0k3fvwdUmDBzcMX-uZAsYfTc"

	// урл выдаст вам игрок или хероку
	WebhookURL = "https://zenkin-sphere-bot.herokuapp.com/"
)

func startTaskBot(ctx context.Context) error {

	bot, err := tgbotapi.NewBotAPI(BotToken)
	if err != nil {
		panic(err)
	}
	bot.Debug = true
	fmt.Printf(WebhookURL)
	fmt.Printf("Authorized on account %s\n", bot.Self.UserName)
	_, err = bot.SetWebhook(tgbotapi.NewWebhook(WebhookURL))
	if err != nil {
		panic(err)
	}
	updates := bot.ListenForWebhook("/")
	port := os.Getenv("PORT")
	go http.ListenAndServe(":"+port, nil)

	var mapStringToSend map[int64]string
	for update := range updates {
		checkUser(update)
		mapStringToSend, err = handleCommand(update)
		if err == nil {
			for num, nextUserStr := range mapStringToSend {
				bot.Send(tgbotapi.NewMessage(num, nextUserStr))
			}
		} else {
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Нет такой команды"))
		}
	}
	return nil
}

func main() {
	err := startTaskBot(context.Background())
	if err != nil {
		panic(err)
	}
}
