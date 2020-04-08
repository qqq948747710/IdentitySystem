package web

import (
	"fmt"
	"github.com/IdentitySystem/web/controller"
	"net/http"
)

func WebStart(app *controller.Application){

	fs:=http.FileServer(http.Dir("web/static"))
	http.Handle("/static/",http.StripPrefix("/static/",fs))

	http.HandleFunc("/",app.IndexView)
	http.HandleFunc("/index.html",app.IndexView)
	http.HandleFunc("/register.html",app.RegisterView)
	http.HandleFunc("/registerdo.html",app.RegisterDo)
	http.HandleFunc("/logindo.html",app.LoginDo)
	http.HandleFunc("/loginout.html",app.LoginOut)
	http.HandleFunc("/myedu.html",app.MyEdu)
	http.HandleFunc("/addedu.html",app.AddEduView)
	http.HandleFunc("/addedudo.html",app.AddEduDo)
	http.HandleFunc("/query.html",app.QueryView)
	http.HandleFunc("/querydo.html",app.QueryDo)
	http.HandleFunc("/upload.html",app.UploadFile)
	err:=http.ListenAndServe(":8000",nil)
	if err != nil {
		fmt.Println("web启动失败",err)
	}
}


//登录验证器
//param1 r http.Request
//param2 w http.Respones
//return logined true;nologin false And Value
func Verlogined(args ...interface{})(bool,interface{}){
	if len(args)!=2{
		fmt.Println("给定验证器验证参数不对")
		return false,""
	}
	w,err:=args[0].(http.ResponseWriter)
	if err == false {
		fmt.Println("类型错误")
		return false,""
	}
	r,err:=args[1].(*http.Request)
	if err == false {
		fmt.Println("类型错误")
		return false,""
	}
	session:=controller.GlobalSessions.SessionStart(w,r)
	value:=session.Get("username")
	if value!=nil{
		return true,value
	}
	return false,""
}

//登录验证器
//param1 r http.Request
//param2 w http.Respones
//return logined true;nologin false And Value
func VerIsAdmin(args ...interface{})(bool,interface{}){
	if len(args)!=2{
		fmt.Println("给定验证器验证参数不对")
		return false,""
	}
	w,err:=args[0].(http.ResponseWriter)
	if err == false {
		fmt.Println("类型错误")
		return false,""
	}
	r,err:=args[1].(*http.Request)
	if err == false {
		fmt.Println("类型错误")
		return false,""
	}
	session:=controller.GlobalSessions.SessionStart(w,r)
	value:=session.Get("isadmin")
	if value=="0"{
		return true,""
	}
	return false,""
}

