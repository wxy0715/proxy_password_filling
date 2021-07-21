package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)
// go get github.com/go-sql-driver/mysql
func main1() {
	db,err := sql.Open("mysql","education:education@(47.101.130.87:3306)/education?charset=utf8")
	defer db.Close()
	if err!=nil {
		fmt.Println(err)
	}else{
		rows, err := db.Query("SELECT nickname FROM ucenter_member")
		if err != nil {
			fmt.Println(err)
		}
		for rows.Next() {
			var nickname string
			err = rows.Scan(&nickname)
			if err != nil {
				panic(err)
			}
			fmt.Println(nickname)
		}
	}
}
