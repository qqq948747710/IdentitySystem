package database

import (
	"database/sql"
	"fmt"
	_"github.com/Go-SQL-Driver/MySQL"
)

func Connesql()(*sql.DB,error){
	db,err:=sql.Open("mysql","root:awdawd123@tcp(127.0.0.1:3306)/identity?charset=utf8")
	fmt.Println(db)
	if err!=nil{
		return nil,fmt.Errorf("mysql信息有误！",err)
	}
	err=db.Ping()
	if err!=nil{
		return nil,fmt.Errorf("mysql链接失败",err)
	}
	return db,nil
}
