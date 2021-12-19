package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
)

var DBHelper *sql.DB

// 初始化数据连接
func initDBHelper() {
	db, err := sql.Open("mysql", "root:123123123@tcp(127.0.0.1:3306)/asknodes?charset=utf8")
	if err != nil {
		log.Fatalf("Open database error: %s\n", err)
	}
	err = db.Ping()
	if err != nil {
		db.Close()
		log.Fatal(err)
	}
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	DBHelper = db
}

// 根据id获取用户昵称
func getNicknameById(id int) (string, error) {
	var nickname string
	sql := "select nickname from user_profile where id = ?"
	err := DBHelper.QueryRow(sql, nickname).Scan(&nickname)
	return nickname, errors.Wrap(err, "dao: "+fmt.Sprint(sql, id))
}

// 结论：大多数场景sql.ErrNoRows应该被warp，然后往上抛 （如果业务不太关心sql.ErrNoRows则可以不用warp，用降级处理）；
// 原因：1.dao层属于application级别的代码，而非第三方库那样可直接返回根因；
//      2.包装带着调用堆栈信息或sql语句再统一记录日志，在查看线上错误或调试时可以避免尴尬
func main() {
	initDBHelper()
	nickname, err := getNicknameById(100)
	if err != nil {
		// 这里记录日志
		fmt.Printf("original error: %T %v\n", errors.Cause(err), errors.Cause(err))
		fmt.Printf("stack trace: \n%+v\n", err)
		return
	}
	fmt.Printf("该用户的昵称是:%s\n", nickname)

}
