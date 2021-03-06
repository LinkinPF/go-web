package main

import (
	"fmt"
	"github.com/jmoiron/sqlx"

	_ "github.com/go-sql-driver/mysql"
)

var db *sqlx.DB

func initDB() (err error) {
	dsn := "root:asdvbn789@tcp(127.0.0.1:3306)/sql_test?charset=utf8mb4&parseTime=True"

	// 如果用MustConnect, 那么连接不成功直接就panic了,不带Must那么就会返回一个错误，然后自己处理
	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		fmt.Printf("connect DB failed, err:%v\n", err)
		return
	}
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(10)

	return
}

func main() {
	err := initDB()
	if err != nil {
		fmt.Printf("init db failed,err:%v\n", err)
		return
	}

	fmt.Println("init db success")
}
