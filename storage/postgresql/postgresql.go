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
	_, err := c.conn.NamedExec(`INSERT INTO products (telegram_user_id, telegram_chat_id, name, weight, inlist, infridge,  timerenable, created_at, finished_at, rest_time)
	VALUES (:telegram_user_id, :telegram_chat_id, :name, :weight, :inlist, :infridge, :timerenable, :created_at, :finished_at, :rest_time)`, p)
	if err != nil {
		log.Printf("AddInfail: %s", err)
		return err
	}
	c.UpdateTimer(p.User_id, p.Chat_id)
	return err
}

func (c *Connection) GetStatus(user_id int64, chat_id int64) int {
	var status int
	c.conn.QueryRow(`SELECT user_status FROM users WHERE telegram_user_id = $1 AND telegram_chat_id = $2`, user_id, chat_id).Scan(&status)
	return status
}

func (c *Connection) GetList(user_id int64, chat_id int64, param int64) ([]product.Product, error) {
	products := make([]product.Product, 0)
	var (
		id          int64
		name        string
		weight      float64
		timerenable bool
		rest_time   int64
	)
	c.UpdateTimer(user_id, chat_id)
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
		c.UpdateTimer(user_id, chat_id)
		sqlStatement := `UPDATE products SET alreadyused = $1, inlist = $2, infridge = $3, intrash = $4, created_at = $5, finished_at = $6, timerenable = $7 WHERE telegram_user_id = $8 AND telegram_chat_id = $9 AND name = $10`
		_, err := c.conn.Exec(sqlStatement, false, false, true, false, from, to, true, user_id, chat_id, name)
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
		c.UpdateTimer(user_id, chat_id)
		sqlStatement := `UPDATE products SET alreadyused = $1, inlist = $2, infridge = $3, intrash = $4, timerenable = $5 WHERE telegram_user_id = $6 AND telegram_chat_id = $7 AND name = $8`
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
		c.UpdateTimer(user_id, chat_id)
		sqlStatement := `UPDATE products SET alreadyused = $1, inlist = $2, infridge = $3, intrash = $4, created_at = $5, finished_at = $6, timerenable = $7 WHERE telegram_user_id = $8 AND telegram_chat_id = $9 AND name = $10`
		_, err := c.conn.Exec(sqlStatement, true, false, false, false, from, to, true, user_id, chat_id, name)
		if err != nil {
			log.Println(err)
		}
		return err
	}
	return existErr
}

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
		sqlStatement := `UPDATE products SET rest_time = $1 WHERE telegram_user_id = $2 AND telegram_chat_id = $3 AND id = $4`
		_, err := c.conn.Exec(sqlStatement, Timer[i].Rest_time, user_id, chat_id, Timer[i].Id)
		if err != nil {
			log.Println(err)
		}
	}
	return err
}
