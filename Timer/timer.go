package timer

import (
	"buy-list/product"
	"buy-list/storage/postgresql"
	"buy-list/tgbot"
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	litime, _       = time.Parse("02-01-2006", "31-12-2030")
	frtime, _       = time.Parse("02-01-2006", "31-12-2030")
	optime, _       = time.Parse("02-01-2006", "31-12-2030")
	liid      int64 = 0
	frid      int64 = 0
	opid      int64 = 0
	li              = true
	fr              = true
	op              = true
	lich            = make(chan bool, 1)
	frch            = make(chan bool, 1)
	opch            = make(chan bool, 1)
)

func SetTimerList(user_id int64, chat_id int64, bot *tgbotapi.BotAPI, db *postgresql.Connection, update tgbotapi.Update) {
	list := db.GetListLess(user_id, chat_id)
	fridge := db.GetFridgeLess(user_id, chat_id)
	open := db.GetOpenLess(user_id, chat_id)
	if list.Finished_at.Before(litime) && list.Id != liid {
		li = true
		litime = list.Finished_at
		liid = list.Id
	}
	if fridge.Finished_at.Before(frtime) && fridge.Id != frid {
		fr = true
		frtime = fridge.Finished_at
		frid = fridge.Id
	}
	if open.Finished_at.Before(optime) && open.Id != opid {
		op = true
		optime = open.Finished_at
		opid = open.Id
	}
	if li {
		if list.Rest_time > 0 {
			li = false
			go TimerList(list, bot, update, db)
		}
	}
	if fr {
		if fridge.Rest_time > 0 {
			fr = false
			go TimerFridge(fridge, bot, update, db)
		}
	}
	if op {
		if open.Rest_time > 0 {
			op = false
			go TimerOpen(open, bot, update, db)
		}
	}
}

func TimerList(p_less product.Product, bot *tgbotapi.BotAPI, update tgbotapi.Update, db *postgresql.Connection) {
	c := time.Tick(3 * time.Second)
	for next := range c {
		if next.After(p_less.Finished_at) {
			msg := "Сработал таймер! Пора купить " + p_less.Name
			log.Println(msg)
			tgbot.SendMessageToUser(bot, update, msg)
			lich <- true
			li = <-lich
			db.UpdateTimer(p_less.User_id, p_less.Chat_id)
			time.Sleep(3 * time.Second)
			SetTimerList(p_less.User_id, p_less.Chat_id, bot, db, update)
			return
		}
	}
}

func TimerFridge(p_less product.Product, bot *tgbotapi.BotAPI, update tgbotapi.Update, db *postgresql.Connection) {
	c := time.Tick(3 * time.Second)
	for next := range c {
		if next.Add(24 * time.Hour).After(p_less.Finished_at) {
			msg := "Сработал таймер! " + p_less.Name + " испортится в холодильнике течение суток!"
			log.Println(msg)
			tgbot.SendMessageToUser(bot, update, msg)
			frch <- true
			fr = <-frch
			db.UpdateTimer(p_less.User_id, p_less.Chat_id)
			time.Sleep(3 * time.Second)
			SetTimerList(p_less.User_id, p_less.Chat_id, bot, db, update)
			return
		}
	}
}

func TimerOpen(p_less product.Product, bot *tgbotapi.BotAPI, update tgbotapi.Update, db *postgresql.Connection) {
	c := time.Tick(3 * time.Second)
	for next := range c {
		if next.Add(24 * time.Hour).After(p_less.Finished_at) {
			msg := "Сработал таймер! " + p_less.Name + " испортится в течение суток!"
			log.Println(msg)
			tgbot.SendMessageToUser(bot, update, msg)
			opch <- true
			op = <-opch
			db.UpdateTimer(p_less.User_id, p_less.Chat_id)
			time.Sleep(3 * time.Second)
			SetTimerList(p_less.User_id, p_less.Chat_id, bot, db, update)
			return
		}
	}
}

func TimerListHour(finished_at time.Time, bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	c := time.Tick(3 * time.Second)
	for next := range c {
		if next.Add(1 * time.Hour).After(finished_at) {
			msg := "До истечения срока меньше часа!"
			log.Println(msg)
			tgbot.SendMessageToUser(bot, update, msg)
			return
		}
	}
}
