package product

import (
	"math"
	"time"
)

type Product struct {
	Id          int64         `db:"id"`
	User_id     int64         `db:"telegram_user_id"`
	Chat_id     int64         `db:"telegram_chat_id"`
	Name        string        `db:"name"`
	Weight      float64       `db:"weight"`
	AlreadyUsed bool          `db:"alreadyused"`
	InList      bool          `db:"inlist"`
	InFridge    bool          `db:"infridge"`
	InTrash     bool          `db:"intrash"`
	Created_at  time.Time     `db:"created_at"`
	Finished_at time.Time     `db:"finished_at"`
	Last_update time.Time     `db:"last_update"`
	Rest_time   time.Duration `db:"rest_time"`
	TimerEnable bool          `db:"timerenable"`
}

func (p *Product) CreateProduct(user_id int64, chat_id int64, name string, weight float64, inlist bool, infridge bool, created_at time.Time, finished_at time.Time, timer bool) *Product {
	p.User_id = user_id
	p.Chat_id = chat_id
	p.Name = name
	p.Weight = (math.Round(weight*100) / 100)
	p.InList = inlist
	p.InFridge = infridge
	p.Created_at = created_at
	p.Finished_at = finished_at
	p.Last_update = time.Now()
	p.Rest_time = p.Finished_at.Sub(time.Now())
	p.TimerEnable = timer
	return p
}
