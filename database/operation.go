package database

import (
	"database/sql"
	"fmt"
)

type DBlogic struct {
	DB *sql.DB
	Name string
}

func(t *DBlogic)Insert(table string,data map[string]string)(sql.Result, error){
	statement:=""
	i:=0
	var values []interface{}
	statement="insert into "+table+"("
	for k,v:=range data{
		if i==0{
			statement+=k
		}else {
			statement+=","+k
		}
		values = append(values, v)
		i++
	}
	statement+=") values("
	for i:=0;i<len(data);i++{
		if i==0{
			statement+="?"
		}else {
			statement+=","+"?"
		}
	}
	statement+=")"
	exc,err:=t.DB.Prepare(statement)
	res,err:=exc.Exec(values...)
	return res,err
}

func(t *DBlogic)Update(table string,data map[string]string,where string)(sql.Result,error){
	statement:="update "+table+" set "
	for k,v:=range data{
		statement+=k+"="+"'"+v+"'"+" "
	}
	statement+=" "+where
	result,err:=t.DB.Exec(statement)
	if err!=nil{
		return nil,fmt.Errorf("更新失败",err)
	}
	return result,nil
}

func (t *DBlogic)Delete(table string,where string)(sql.Result,error){
	statement:="delete from "+table+" "+where
	res,err:=t.DB.Exec(statement)
	if err!=nil{
		return nil,fmt.Errorf("删除失败",err)
	}
	return res,nil
}

func(t *DBlogic) Query(table string,field string,where string)(map[int]map[string]string,error) {
	statement:=""
	if field!=""{
		statement="select "+field+" from "+table+" "+where
	}else {
		statement="select * from "+table+" "+where
	}
	exec,err:=t.DB.Prepare(statement)
	rows,err:=exec.Query()
	if err!=nil{
		return nil,fmt.Errorf("查询失败",err)
	}
	colums,err:=rows.Columns()
	if err!=nil{
		return nil,fmt.Errorf("查询获取字段失败",err)
	}
	datas:=make(map[int]map[string]string)
	values:=make([][]byte,len(colums))
	scans:=make([]interface{},len(colums))
	for i:=range values{
		scans[i]=&values[i]
	}
	for i:=0;rows.Next();i++{
		rows.Scan(scans...)
		data:=make(map[string]string)
		for k,v:=range values{
			data[colums[k]]=string(v)
		}
		datas[i]=data
	}
	return datas,nil
}
