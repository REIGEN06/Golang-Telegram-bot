package statuses

import (
	"buy-list/product"
	"buy-list/storage/postgresql"
	"log"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func StatusHandler(status int, user_msg string, user_id int64, chat_id int64, db *postgresql.Connection) string {
	//Add in List
	if status == 1 {
		msg := "Продукт указан не верно"
		text := strings.Fields(user_msg)
		if len(text) == 4 {
			weight, errFloat := strconv.ParseFloat(text[1], 32)
			date, errDate := time.Parse("02-01-2006 15:04", text[2]+" "+text[3])
			date = date.Add(-3 * time.Hour)
			if date.Before(time.Now()) {
				msg = "Ошибка: Дата уже истекла"
			}
			if (reflect.TypeOf(text[0]) == reflect.TypeOf(user_msg)) && (errFloat == nil) && (errDate == nil) && date.After(time.Now()) {
				p := product.Product{}
				p.CreateProduct(user_id, chat_id, text[0], weight, true, false, time.Now(), date, true)
				err := db.AddIn(&p)

				msg = "Продукт добавлен в список!\nТаймер заведён."
				if err != nil {
					msg = "Ошибка при добавлении."
				}
			}
		}
		if len(text) == 1 {
			p := product.Product{}
			p.CreateProduct(user_id, chat_id, text[0], 0, true, false, time.Now(), time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC), false)
			err := db.AddIn(&p)

			msg = "Продукт добавлен в список!\nБез таймера."
			if err != nil {
				msg = "Ошибка при добавлении."
			}
		}
		if len(text) == 3 && text[1] == "с" {
			today := time.Now().Format("02-01-2006")
			today += " " + text[2]
			log.Println(today)
			date, dateerr := time.Parse("02-01-2006 15:04", today)
			log.Println(date)
			date = date.Add(-3 * time.Hour)
			if date.Before(time.Now()) {
				msg = "Ошибка: Дата уже истекла"
			}
			if dateerr == nil && date.After(time.Now()) {
				p := product.Product{}
				p.CreateProduct(user_id, chat_id, text[0], 0, true, false, time.Now(), date, true)
				err := db.AddIn(&p)
				msg = "Продукт добавлен в список!\nТаймер заведён."
				if err != nil {
					msg = "Ошибка при добавлении в базу данных."
				}
			} else {
				msg = "С датой что-то не так.."
			}
		}
		db.SetStatus(0, user_id, chat_id)
		return msg
	}
	//Add in Fridge
	if status == 2 {
		msg := "Продукт указан не верно"
		text := strings.Fields(user_msg)
		if len(text) == 3 {
			from, fromerr := time.Parse("02-01-2006", text[1])
			to, toerr := time.Parse("02-01-2006", text[2])
			if to.Before(from) {
				msg = "Ошибка: Правая дата раньше левой"
			}
			if to.Before(time.Now()) {
				msg = "Ошибка: Правая дата уже истекла"
			}
			if (fromerr == nil) && (toerr == nil) && to.After(from) && to.After(time.Now()) {
				p := product.Product{}
				p.CreateProduct(user_id, chat_id, text[0], 0, false, true, from, to, true)
				err := db.AddIn(&p)

				msg = text[0] + " добавлен в холодильник!"
				if err != nil {
					msg = "Ошибка при добавлении."
				}
			}
		}
		//купил - перенос из листа
		if len(text) == 4 {
			from, fromerr := time.Parse("02-01-2006", text[2])
			to, toerr := time.Parse("02-01-2006", text[3])
			if to.Before(from) && to.Before(time.Now()) {
				msg = "Некорректное время"
			}
			if (fromerr == nil) && (toerr == nil) && to.After(from) && to.After(time.Now()) {
				err := db.SetFridge(user_id, chat_id, from, to, text[1])

				msg = text[1] + " перенесён в холодильник!"
				if err != nil {
					msg = "Ошибка при переносе."
					if err.Error() == "NOT_EXISTS" {
						msg = text[1] + " не существует"
					}
				}
			}
		}
		db.SetStatus(0, user_id, chat_id)
		return msg
	}
	//Open product
	if status == 3 {
		msg := "Продукт указан не верно"
		text := strings.Fields(user_msg)
		if len(text) == 3 {
			from, fromerr := time.Parse("02-01-2006", text[1])
			to, toerr := time.Parse("02-01-2006", text[2])
			if to.Before(from) {
				msg = "Ошибка: Правая дата раньше левой"
			}
			if to.Before(time.Now()) {
				msg = "Ошибка: Правая дата уже истекла"
			}
			if (fromerr == nil) && (toerr == nil) && to.After(from) && to.After(time.Now()) {
				err := db.SetUsed(user_id, chat_id, from, to, text[0])

				msg = text[0] + " открыт, срок хранения обновлен!"
				if err != nil {
					msg = "Ошибка при открытии продукта."
					if err.Error() == "NOT_EXISTS" {
						msg = text[0] + " не существует"
					}
				}
			}

		}
		db.SetStatus(0, user_id, chat_id)
		return msg
	}
	//Finish product
	if status == 4 {
		msg := "Продукт указан не верно"
		text := strings.Fields(user_msg)
		if len(text) == 1 {
			msg = text[0] + " успешно приготовлен/выброшен!"
			err := db.SetTrash(user_id, chat_id, text[0])
			if err != nil {
				msg = "Ошибка при приготовлении/выбрасывании."
				if err.Error() == "NOT_EXISTS" {
					msg = text[0] + " не существует"
				}
			}
		}
		db.SetStatus(0, user_id, chat_id)
		return msg
	}
	//List
	if status == 5 {
		msg := "Неправильно выбран список"
		text := strings.Fields(user_msg)
		if len(text) == 1 {
			if text[0] == "1" || text[0] == "2" || text[0] == "3" {
				param, _ := strconv.ParseInt(text[0], 10, 64)
				msg = "Список продуктов"
				if text[0] == "1" {
					msg += " по алфавиту:"
				} else if text[0] == "2" {
					msg += " в холодильнике по сроку годности:"
				} else if text[0] == "3" {
					msg += " (полный):"
				}
				products, _ := db.GetList(user_id, chat_id, param)
				for i := 0; i < len(products); i++ {
					msg += "\n" + strconv.Itoa(i+1) + ". " + products[i].Name
					if products[i].Weight != 0 {
						msg += ", " + strconv.FormatFloat(products[i].Weight, 'f', 0, 64) + "гр."
					}
					if products[i].TimerEnable && products[i].Rest_time > 0 {
						msg += ", таймер включен"
					}
					if (text[0] == "2" || text[0] == "3") && products[i].TimerEnable {
						if products[i].Rest_time > 0 { //&& products[i].Rest_time < 20000
							rtime := products[i].Rest_time.String()
							rtime = strings.ReplaceAll(rtime, "h", " часов, ")
							rtime = strings.ReplaceAll(rtime, "m", " минут, ")
							rtime = strings.ReplaceAll(rtime, "s", " секунд ")
							msg += ", испортится через: " + rtime
						} else {
							msg += ", срок годности вышел"
						}
					}
					if text[0] == "3" && products[i].InTrash {
						msg += " [Приготовлен/Выкинут]"
					} else if text[0] == "3" && products[i].AlreadyUsed {
						msg += " [Открыт]"
					} else if text[0] == "3" && products[i].InList {
						msg += " [Список продуктов]"
					} else if text[0] == "3" && products[i].InFridge {
						msg += " [Холодильник]"
					}
				}

			}
		}
		db.SetStatus(0, user_id, chat_id)
		return msg
	}
	//Status time
	if status == 6 {
		msg := "Слишком много параметров"
		text := strings.Fields(user_msg)
		if len(text) == 2 {
			from, fromerr := time.Parse("02-01-2006", text[0])
			to, toerr := time.Parse("02-01-2006", text[1])
			if to.Before(from) {
				msg = "Ошибка: Правая дата раньше левой"
			}
			if to.Before(time.Now()) {
				msg = "Ошибка: Правая дата уже истекла"
			}
			if (fromerr == nil) && (toerr == nil) && to.After(from) && to.After(time.Now()) {
				alreadyused, trashcount, _ := db.StatusTime(user_id, chat_id, from, to)
				msg = "В этот промежуток времени нет использованных продуктов"
				if len(alreadyused) != 0 {
					msg = "Список использованных продуктов за промежуток с " + from.Format("2006-01-02") + " по " + to.Format("2006-01-02") + ":"
					for i := 0; i < len(alreadyused); i++ {
						msg += "\n" + strconv.Itoa(i+1) + ". " + alreadyused[i].Name
					}
				}
				msg += "\n\nКоличество приготовленных/выкинутых продуктов: " + strconv.Itoa(trashcount)
			}
		}
		db.SetStatus(0, user_id, chat_id)
		return msg
	}
	return "Сообщение не распознано\nДля справки воспользуйтесь командой /help"
}
