package postgresql

import (
	"buy-list/product"
	"log"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
)

var (
	nickname       string = " "
	name           string = "Чипсы"
	user_id        int64  = 463372101
	chat_id        int64  = 463372101
	created_at            = time.Now()
	finished_at, _        = time.Parse("02-01-2006", "31-12-2025")
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}

}

func TestCreateUser(t *testing.T) {
	db := Connect(os.Getenv("TELEGRAM_APITOKEN"))
	db.CreateUser(nickname, name, user_id, chat_id)
}
func TestAddIn(t *testing.T) {
	db := Connect(os.Getenv("TELEGRAM_APITOKEN"))
	p := product.Product{User_id: -999, Chat_id: -999}
	db.AddIn(&p)
}
func TestGetStatus(t *testing.T) {
	db := Connect(os.Getenv("TELEGRAM_APITOKEN"))
	db.GetStatus(user_id, chat_id)
}
func TestGetList(t *testing.T) {
	db := Connect(os.Getenv("TELEGRAM_APITOKEN"))
	db.GetList(user_id, chat_id, 3)
}
func TestStatusTime(t *testing.T) {
	db := Connect(os.Getenv("TELEGRAM_APITOKEN"))
	db.StatusTime(user_id, chat_id, created_at, finished_at)
}
func TestSetStatus(t *testing.T) {
	db := Connect(os.Getenv("TELEGRAM_APITOKEN"))
	db.SetStatus(1, user_id, chat_id)
}
func TestSetFridge(t *testing.T) {
	db := Connect(os.Getenv("TELEGRAM_APITOKEN"))
	db.SetFridge(user_id, chat_id, created_at, finished_at, name)
}
func TestSetTrash(t *testing.T) {
	db := Connect(os.Getenv("TELEGRAM_APITOKEN"))
	db.SetTrash(user_id, chat_id, name)
}
func TestSetUsed(t *testing.T) {
	db := Connect(os.Getenv("TELEGRAM_APITOKEN"))
	db.SetUsed(user_id, chat_id, created_at, finished_at, name)
}
func TestUpdateTimer(t *testing.T) {
	db := Connect(os.Getenv("TELEGRAM_APITOKEN"))
	db.UpdateTimer(user_id, chat_id)
}
