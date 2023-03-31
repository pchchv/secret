package main

import (
	"fmt"
	"reflect"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pchchv/golog"
)

func tgbot() {
	bot, err := tgbotapi.NewBotAPI(getEnvValue("TG_BOT_TOKEN"))
	if err != nil {
		golog.Fatal("Failed to create bot: %s", err.Error())
	}
	bot.Debug = false

	golog.Info("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if reflect.TypeOf(update.Message.Text).Kind() != reflect.String || update.Message.Text == "" {
			continue
		}

		var msg tgbotapi.MessageConfig
		userMessage := update.Message.Text

		switch userMessage {
		case "/start":
			msg = tgbotapi.NewMessage(update.Message.Chat.ID,
				"Hi, I am a bot for creating one-time secret notes.\nTo create a note send me the command /send.\nTo retrieve a note, send /get.")
		case "/send":
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Send me your secret.")
		case "/get":
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Send me a password.")
		default:
			if strings.Contains(userMessage, " ") {
				pass, err := encryptor(userMessage)
				if err != nil {
					msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Something went wrong. Try again.")
					golog.Info("Encryption error: %s", err.Error())
				} else {
					msg = tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Your password: '%v'", pass))
				}
			} else {
				text, err := decryptor(userMessage)
				if err != nil {
					msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Something went wrong. Try again.")
					golog.Info("Decryption error: %s", err.Error())
				} else {
					msg = tgbotapi.NewMessage(update.Message.Chat.ID, text)
				}
			}
		}

		if _, err := bot.Send(msg); err != nil {
			golog.Info("Failed to send message: %s", err.Error())
		}
	}
}
