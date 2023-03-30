package main

import (
	"reflect"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pchchv/golog"
)

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

		// Check that a command message was received from the user
		if reflect.TypeOf(update.Message.Command()).Kind() == reflect.String && update.Message.Command() != "" {
			switch update.Message.Command() {
				// Sending a message
			case "start":
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, `Hi, I am a bot for creating one-time secret notes,
				to create a note send me the command /send.
				To retrieve a note, send /get.`)
				bot.Send(msg)
			case "send":
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Send me your secret.")
				bot.Send(msg)
			case "get":
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Send me a password.")
				bot.Send(msg)
			}
		}
	}
}
