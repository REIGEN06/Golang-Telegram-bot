package main

import (
	timer "buy-list/Timer"
	"buy-list/commands"
	"buy-list/statuses"
	"buy-list/storage/postgresql"
	"buy-list/tgbot"
	"log"
	"os"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

// как лучше сделать таймер? как вообще делать его
// запустить программу в goroutine как отдельный процесс

func main() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_APITOKEN"))
	if err != nil {
		panic(err)
	}
	db := postgresql.Connect(os.Getenv("POSTGRESQL_TOKEN"))

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 90
	updates := bot.GetUpdatesChan(updateConfig)
	for update := range updates {
		if update.Message == nil {
			continue
		}
		wg := sync.WaitGroup{}
		user_nickname := update.SentFrom().UserName
		user_msg := update.Message.Text
		user_name := update.Message.From.FirstName
		user_id := update.Message.From.ID
		chat_id := update.Message.Chat.ID
		status := db.GetStatus(user_id, chat_id)

		log.Printf("nickname: %s \n", user_nickname)
		log.Printf("user_name: %s \n", user_name)
		log.Printf("msg: %s \n", user_msg)
		log.Println()
		if update.Message.IsCommand() {
			msg := commands.CommandHandler(user_msg, user_id, chat_id, user_nickname, user_name, db)
			tgbot.SendMessageToUser(bot, update, msg)
		} else {
			msg := statuses.StatusHandler(status, user_msg, user_id, chat_id, db)
			tgbot.SendMessageToUser(bot, update, msg)
		}
		wg.Add(1)
		db.UpdateTimer(user_id, chat_id)
		wg.Done()
		timer.SetTimerList(user_id, chat_id, bot, db, update)
	}

}
