package main

import (
	"reflect"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pchchv/golog"
)

const startmessage = "Hi, I am a bot for creating one-time secret notes, to create a note send me the command /send. To retrieve a note, send /get."

func tgbot() {
	bot, err := tgbotapi.NewBotAPI(getEnvValue("TG_BOT_TOKEN"))
	if err != nil {
		golog.Panic(err.Error())
	}

	golog.Info("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}
		// Check that a text message was received from the user
		if reflect.TypeOf(update.Message.Text).Kind() == reflect.String && update.Message.Text != "" {
			switch update.Message.Text {
			case "/start":
				//Отправлем сообщение
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, startmessage)
				bot.Send(msg)
			}
		}
	}
}
