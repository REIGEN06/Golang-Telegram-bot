package product

import (
	"testing"
	"time"
)

var (
	user_id        int64 = 1000000
	chat_id        int64 = 1000000
	name                 = "tester"
	weight               = 2.2
	inList               = true
	inFridge             = false
	created_at           = time.Now()
	finished_at, _       = time.Parse("02-01-2006", "31-12-2025")
	timer                = true
)

func TestProduct(t *testing.T) {
	p := Product{}
	p.CreateProduct(user_id, chat_id, name, weight, inList, inFridge, created_at, finished_at, timer)
}

func BenchmarkProduct(b *testing.B) {
	for i := 0; i < b.N; i++ {
		p := Product{}
		p.CreateProduct(user_id, chat_id, name, weight, inList, inFridge, created_at, finished_at, timer)
	}
}
