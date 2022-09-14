package db

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func NewDb(dsn string) *sql.DB {
	var db *sql.DB
	for count := 0; ; {
		var err error
		db, err = sql.Open("mysql", dsn)
		if err != nil {
			log.Printf("%v\n", err)
			count++
		} else {
			err = db.Ping()
			if err == nil {
				log.Printf("Db OK\n")
				break
			} else {
				log.Printf("%v\n", err)
				count++
			}
		}
		if count > 5 {
			panic("db fail")
		} else {
			count++
			log.Printf("[%d]Wait Db...\n", count)
			time.Sleep(time.Duration(count*count+1) * time.Second)
		}
	}
	return db
}
