package test

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"container/list"
)

var db *sql.DB

func init() {
	db, _ = sql.Open("mysql", "root:123.com@tcp(10.0.30.120:3306)/iov?charset=utf8")
	db.SetMaxOpenConns(2000)
	db.SetMaxIdleConns(1000)
	db.Ping()
}
func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}
func main() {
	l := list.New()

	rows, err := db.Query("SELECT * FROM area")
	defer rows.Close()
	checkErr(err)

	columns, _ := rows.Columns()
	fmt.Println(len(columns))
	scanArgs := make([]interface{}, len(columns))
	values := make([]interface{}, len(columns))
	for j := range values {
		scanArgs[j] = &values[j]
	}

	record := make(map[string]string)
	for rows.Next() {
		//将行数据保存到record字典
		err = rows.Scan(scanArgs...)
		for i, col := range values {
			if col != nil {
				record[columns[i]] = string(col.([]byte))
			}
		}
		//fmt.Println(record)
		l.PushBack(record)
	}

	//fmt.Println(record)
	fmt.Println(l.Len())
}
