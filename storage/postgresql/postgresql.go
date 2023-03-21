package postgresql

import (
	"buy-list/product"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Connection struct {
	conn *sqlx.DB
}

func Connect(connString string) *Connection {
	conn, err := sqlx.Connect("pgx", connString)
	if err != nil {
		panic(err)
	}

	driver, err := postgres.WithInstance(conn.DB, &postgres.Config{})
	if err != nil {
		panic(err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres", driver)

	if err != nil {
		panic(err)
	}
	//не создает переменную err
	if err := m.Up(); err != nil {
		if err.Error() == "no change" {
			println("tables are already migrated!")
		} else {
			panic(err)
		}
	} else {
		println("successfully migrated!")
	}

	return &Connection{
		conn: conn,
	}
}

func (c *Connection) CreateUser(nickname string, name string, user_id int64, chat_id int64) string {
	msg := "Теперь вы можете пользоваться ботом, "
	msg += name
	msg += "!😊"
	sqlStatement := `INSERT INTO users  (telegram_user_nickname, telegram_user_name, telegram_user_id, telegram_chat_id) VALUES ($1, $2, $3, $4)`
	_, err := c.conn.Exec(sqlStatement, nickname, name, user_id, chat_id)

	if err != nil {
		msg := "Неизвестная ошибка. Обратитесь к разработчику"
		if strings.Contains(err.Error(), "повторяющееся") {
			msg = name
			msg += ", Вы уже зарегистрированы!"
		}
		log.Printf("CreateUser fail: %s", err)
		return msg
	}
	return msg
}

func (c *Connection) AddIn(p *product.Product) error {
	_, err := c.conn.NamedExec(`INSERT INTO products (telegram_user_id, telegram_chat_id, name, weight, inlist, infridge,  timerenable, created_at, finished_at, rest_time, last_update)
	VALUES (:telegram_user_id, :telegram_chat_id, :name, :weight, :inlist, :infridge, :timerenable, :created_at, :finished_at, :rest_time, :last_update)`, p)
	if err != nil {
		log.Printf("AddInfail: %s", err)
		return err
	}
	return err
}

func (c *Connection) GetStatus(user_id int64, chat_id int64) int {
	var status int
	c.conn.QueryRow(`SELECT user_status FROM users WHERE telegram_user_id = $1 AND telegram_chat_id = $2`, user_id, chat_id).Scan(&status)
	return status
}

func (c *Connection) GetListLess(user_id int64, chat_id int64) product.Product {
	var (
		p           product.Product
		id          int64
		finished_at time.Time
		name        string
	)
	c.conn.QueryRow(`SELECT id, name, finished_at FROM products WHERE telegram_user_id = $1 AND telegram_chat_id = $2 AND inlist = $3 AND timerenable = $4 ORDER BY rest_time limit 1`, user_id, chat_id, true, true).Scan(&id, &name, &finished_at)
	p.CreateProduct(user_id, chat_id, name, 0, true, false, time.Now(), finished_at, true)
	p.Id = id
	log.Println("List: ", p.Id, p.User_id, p.Finished_at, p.Rest_time)
	return p
}

func (c *Connection) GetFridgeLess(user_id int64, chat_id int64) product.Product {
	var (
		p           product.Product
		id          int64
		finished_at time.Time
		name        string
	)
	c.conn.QueryRow(`SELECT id, name, finished_at FROM products WHERE telegram_user_id = $1 AND telegram_chat_id = $2 AND infridge = $3 AND timerenable = $4 ORDER BY rest_time limit 1`, user_id, chat_id, true, true).Scan(&id, &name, &finished_at)
	p.CreateProduct(user_id, chat_id, name, 0, true, false, time.Now(), finished_at, true)
	p.Id = id
	log.Println("Fridge: ", p.Id, p.User_id, p.Finished_at, p.Rest_time)
	return p
}

func (c *Connection) GetOpenLess(user_id int64, chat_id int64) product.Product {
	var (
		p           product.Product
		id          int64
		finished_at time.Time
		name        string
	)
	c.conn.QueryRow(`SELECT id, name, finished_at FROM products WHERE telegram_user_id = $1 AND telegram_chat_id = $2 AND alreadyused = $3 AND timerenable = $4 ORDER BY rest_time limit 1`, user_id, chat_id, true, true).Scan(&id, &name, &finished_at)
	p.CreateProduct(user_id, chat_id, name, 0, true, false, time.Now(), finished_at, true)
	p.Id = id
	log.Println("Open: ", p.Id, p.User_id, p.Finished_at, p.Rest_time)
	return p
}

func (c *Connection) GetList(user_id int64, chat_id int64, param int64) ([]product.Product, error) {
	products := make([]product.Product, 0)
	var (
		id          int64
		name        string
		weight      float64
		timerenable bool
		rest_time   int64
		alreadyused bool
		inlist      bool
		infridge    bool
		intrash     bool
	)
	if param == 1 {
		rows, err := c.conn.Query(`SELECT id, name, weight, timerenable FROM products WHERE telegram_user_id = $1 AND telegram_chat_id = $2 AND inlist = $3 ORDER BY name ASC`, user_id, chat_id, true)
		if err != nil {
			return nil, fmt.Errorf("GetList NameASC error: %w", err)
		}
		defer rows.Close()
		for rows.Next() {
			if err := rows.Scan(&id, &name, &weight, &timerenable); err != nil {
				fmt.Fprintf(os.Stderr, "GetList scan failed: %v\n", err)
				return nil, nil
			}
			products = append(products, product.Product{Id: id, Name: name, Weight: weight, TimerEnable: timerenable})
		}
		return products, err
	} else if param == 2 {
		rows, err := c.conn.Query(`SELECT id, name, weight, timerenable, rest_time FROM products WHERE telegram_user_id = $1 AND telegram_chat_id = $2 AND infridge = $3 ORDER BY rest_time ASC`, user_id, chat_id, true)
		if err != nil {
			return nil, fmt.Errorf("GetList LastTime error: %w", err)
		}
		defer rows.Close()
		for rows.Next() {
			if err := rows.Scan(&id, &name, &weight, &timerenable, &rest_time); err != nil {
				fmt.Fprintf(os.Stderr, "GetList scan failed: %v\n", err)
				return nil, nil
			}
			products = append(products, product.Product{Id: id, Name: name, Weight: weight, TimerEnable: timerenable, Rest_time: time.Duration(rest_time)})
		}
		return products, err
	} else if param == 3 {
		rows, err := c.conn.Query(`SELECT id, name, weight, rest_time, timerenable, alreadyused, inlist, infridge, intrash FROM products WHERE telegram_user_id = $1 AND telegram_chat_id = $2 AND intrash =$3`, user_id, chat_id, false)
		if err != nil {
			return nil, fmt.Errorf("GetListAll error: %w", err)
		}
		defer rows.Close()
		for rows.Next() {
			if err := rows.Scan(&id, &name, &weight, &rest_time, &timerenable, &alreadyused, &inlist, &infridge, &intrash); err != nil {
				return nil, fmt.Errorf("GetListAll scan failed: %w\n", err)
			}
			products = append(products, product.Product{Id: id, Name: name, Weight: weight, TimerEnable: timerenable, Rest_time: time.Duration(rest_time), AlreadyUsed: alreadyused, InList: inlist, InFridge: infridge, InTrash: intrash})
		}
		return products, err
	} else if param == 4 {
		rows, err := c.conn.Query(`SELECT id, name, weight, rest_time FROM products WHERE telegram_user_id = $1 AND telegram_chat_id = $2 AND alreadyused = $3 ORDER BY rest_time ASC`, user_id, chat_id, true)
		if err != nil {
			return nil, fmt.Errorf("GetList LastTime error: %w", err)
		}
		defer rows.Close()
		for rows.Next() {
			if err := rows.Scan(&id, &name, &weight, &rest_time); err != nil {
				fmt.Fprintf(os.Stderr, "GetList scan failed: %v\n", err)
				return nil, nil
			}
			products = append(products, product.Product{Id: id, Name: name, Weight: weight, TimerEnable: timerenable, Rest_time: time.Duration(rest_time)})
		}
		return products, err
	}
	return products, nil
}

func (c *Connection) StatusTime(user_id int64, chat_id int64, from time.Time, to time.Time) ([]product.Product, int, error) {
	used_prod := make([]product.Product, 0) //список всех использованных
	trash_prod := 0                         //количество всех выкинутых\приготовленных
	var (
		name        string
		last_update time.Time
	)
	//использованные
	rows, err := c.conn.Query(`SELECT name, last_update FROM products WHERE telegram_user_id = $1 AND telegram_chat_id = $2 AND alreadyused = $3 ORDER BY name ASC`, user_id, chat_id, true)
	if err != nil {
		return nil, 0, fmt.Errorf("StatusTime alreadyused error: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&name, &last_update); err != nil {
			return nil, 0, fmt.Errorf("StatusTime alreadyused scan failed: %w", err)
		}
		if last_update.After(from) && last_update.Before(to) {
			used_prod = append(used_prod, product.Product{Name: name, Last_update: last_update})
		}
	}

	rows, err = c.conn.Query(`SELECT last_update FROM products WHERE telegram_user_id = $1 AND telegram_chat_id = $2 AND intrash = $3 ORDER BY name ASC`, user_id, chat_id, true)
	if err != nil {
		return nil, 0, fmt.Errorf("StatusTime intrash error: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&last_update); err != nil {
			return nil, 0, fmt.Errorf("StatusTime intrash scan failed: %w", err)
		}
		if last_update.After(from) && last_update.Before(to) {
			trash_prod += 1
		}
	}
	return used_prod, trash_prod, err
}

func (c *Connection) SetStatus(newStatus int, user_id int64, chat_id int64) {
	sqlStatement := `UPDATE users SET user_status = $1 WHERE telegram_user_id = $2 AND telegram_chat_id = $3`
	_, err := c.conn.Exec(sqlStatement, newStatus, user_id, chat_id)

	if err != nil {
		log.Println(err)
	}
}

func (c *Connection) SetFridge(user_id int64, chat_id int64, from time.Time, to time.Time, name string) error {
	var exists int
	c.conn.QueryRow(`SELECT id FROM products WHERE telegram_user_id = $1 AND telegram_chat_id = $2 AND name = $3`, user_id, chat_id, name).Scan(&exists)
	existErr := errors.New("NOT_EXISTS")
	if exists != 0 {
		sqlStatement := `UPDATE products SET alreadyused = $1, inlist = $2, infridge = $3, intrash = $4, created_at = $5, finished_at = $6, timerenable = $7, last_update = NOW() WHERE telegram_user_id = $8 AND telegram_chat_id = $9 AND name = $10`
		_, err := c.conn.Exec(sqlStatement, false, false, true, false, from, to.Add(time.Duration(exists)*time.Millisecond*100), true, user_id, chat_id, name)
		if err != nil {
			log.Println(err)
		}
		return err
	}
	return existErr
}

func (c *Connection) SetTrash(user_id int64, chat_id int64, name string) error {
	var exists int
	c.conn.QueryRow(`SELECT id FROM products WHERE telegram_user_id = $1 AND telegram_chat_id = $2 AND name = $3`, user_id, chat_id, name).Scan(&exists)
	existErr := errors.New("NOT_EXISTS")
	if exists != 0 {
		sqlStatement := `UPDATE products SET alreadyused = $1, inlist = $2, infridge = $3, intrash = $4, timerenable = $5, last_update = NOW() WHERE telegram_user_id = $6 AND telegram_chat_id = $7 AND name = $8`
		_, err := c.conn.Exec(sqlStatement, true, false, false, true, false, user_id, chat_id, name)
		if err != nil {
			log.Println(err)
		}
		return err
	}
	return existErr
}

func (c *Connection) SetUsed(user_id int64, chat_id int64, from time.Time, to time.Time, name string) error {
	var exists int
	c.conn.QueryRow(`SELECT id FROM products WHERE telegram_user_id = $1 AND telegram_chat_id = $2 AND name = $3`, user_id, chat_id, name).Scan(&exists)
	existErr := errors.New("NOT_EXISTS")
	if exists != 0 {
		sqlStatement := `UPDATE products SET alreadyused = $1, inlist = $2, infridge = $3, intrash = $4, created_at = $5, finished_at = $6, timerenable = $7, last_update = NOW() WHERE telegram_user_id = $8 AND telegram_chat_id = $9 AND name = $10`
		_, err := c.conn.Exec(sqlStatement, true, false, false, false, from, to.Add(time.Duration(exists)*time.Millisecond*100), true, user_id, chat_id, name)
		if err != nil {
			log.Println(err)
		}
		return err
	}
	return existErr
}

// Знаю, что при большом количестве продуктов нужно использовать горутины, но нет времени на рефакторинг
func (c *Connection) UpdateTimer(user_id int64, chat_id int64) error {
	var finished_at time.Time
	var id int64
	Timer := make([]product.Product, 0)
	rows, err := c.conn.Query(`SELECT id, finished_at FROM products WHERE telegram_user_id = $1 AND telegram_chat_id = $2`, user_id, chat_id)
	if err != nil {
		return fmt.Errorf("GetList Timer error: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&id, &finished_at); err != nil {
			fmt.Fprintf(os.Stderr, "GetList TimerScan failed: %v\n", err)
			return nil
		}
		Timer = append(Timer, product.Product{Id: id, Finished_at: finished_at})
	}
	for i := 0; i < len(Timer); i++ {
		Timer[i].Rest_time = (Timer[i].Finished_at).Sub(time.Now())
	}
	for i := 0; i < len(Timer); i++ {
		if Timer[i].Rest_time < 0 {
			sqlStatement := `UPDATE products SET rest_time = $1, timerenable = $2 WHERE telegram_user_id = $3 AND telegram_chat_id = $4 AND id = $5`
			_, err := c.conn.Exec(sqlStatement, Timer[i].Rest_time, false, user_id, chat_id, Timer[i].Id)
			if err != nil {
				log.Println(err)
			}
		} else {
			sqlStatement := `UPDATE products SET rest_time = $1 WHERE telegram_user_id = $2 AND telegram_chat_id = $3 AND id = $4`
			_, err := c.conn.Exec(sqlStatement, Timer[i].Rest_time, user_id, chat_id, Timer[i].Id)
			if err != nil {
				log.Println(err)
			}
		}
	}
	return err
}
