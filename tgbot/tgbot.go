package tgbot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func SendMessageToUser(bot *tgbotapi.BotAPI, update tgbotapi.Update, msg string) {
	bot_msg := tgbotapi.NewMessage(update.Message.From.ID, msg)
	bot_msg.ReplyToMessageID = update.Message.MessageID
	if _, err := bot.Send(bot_msg); err != nil {
		panic(err)
	}
}
