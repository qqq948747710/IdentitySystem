package controller

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
)

func showView(w http.ResponseWriter,r *http.Request,templateName string,data interface{}){
	path:=filepath.Join("web","tpl",templateName)
	page,err:=template.ParseFiles(path)
	if err!=nil{
		fmt.Println("模板创建失败",err)
		return
	}
	page.Execute(w,data)
	if err != nil {
		fmt.Println("模板融合失败",err)
		return
	}
}

//构建返回消息回信模板
func message(w http.ResponseWriter,r *http.Request,message string)  {
	data:=&struct {
		Message string
	}{Message:message}
	showView(w,r,"message.html",data)
}