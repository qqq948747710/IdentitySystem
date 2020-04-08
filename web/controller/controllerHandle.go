package controller

import (
	"crypto/rand"
	"fmt"
	"github.com/IdentitySystem/database"
	"github.com/IdentitySystem/service"
	_ "github.com/IdentitySystem/web/memory"
	"github.com/IdentitySystem/web/session"
	"github.com/IdentitySystem/web/stateverifier"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
)
const VERLOGINED string="verlogined"
const VERISADMIN string="verisadmin"
type Application struct {
	Fabric *service.SrviceSetup
	Dblogic *database.DBlogic
	Verifier *stateverifier.Verifier
}

var GlobalSessions *session.Manager
func (app *Application)IndexView(w http.ResponseWriter,r *http.Request){
	islogin,username:=app.Verifier.Verify(VERLOGINED,w,r)
	data:=struct {
		Islogin bool
		Username string
	}{islogin,username.(string)}
	showView(w,r,"index.html",data)
}
func (app *Application)RegisterView(w http.ResponseWriter,r *http.Request){
	islogin,username:=app.Verifier.Verify(VERLOGINED,w,r)
	data:=struct {
		Islogin bool
		Username string
	}{islogin,username.(string)}
	showView(w,r,"register.html",data)
}
func(app *Application)AddEduView(w http.ResponseWriter,r *http.Request){
	islogin,username:=app.Verifier.Verify(VERLOGINED,w,r)
	if !islogin{
		message(w,r,"你还没有登录!")
		return
	}
	data:=struct {
		Islogin bool
		Username string
	}{islogin,username.(string)}
	showView(w,r,"addedu.html",data)
}
func(app *Application)QueryView(w http.ResponseWriter,r *http.Request) {
	islogin, username := app.Verifier.Verify(VERLOGINED, w, r)
	if !islogin {
		message(w, r, "你还没有登录!")
		return
	}
	data := struct {
		Islogin  bool
		Username string
	}{islogin, username.(string)}
	showView(w, r, "query.html", data)
}
func(app *Application)RegisterDo(w http.ResponseWriter,r *http.Request){
	islogin,_:=app.Verifier.Verify(VERLOGINED,w,r)
	if islogin{
		message(w,r,"您已经登录过了")
		return
	}
	r.ParseForm()
	username:=r.FormValue("username")
	password:=r.FormValue("password")
	repassword:=r.FormValue("repassword")
	email:=r.FormValue("email")
	entityid:=r.FormValue("entityid")
	if username==""||password==""||repassword==""||email==""||entityid=="" {
		message(w,r,"发送信息有误！")
		return
	}
	if password!=repassword{
		message(w,r,"两次密码不一致！")
		return
	}
	values,err:=app.Dblogic.Query("user","username","where username='"+username+"' or entityid='"+entityid+"'")
	if err!=nil {
		fmt.Println(err)
		return
	}
	fmt.Println(values)
	if len(values)!=0{
		message(w,r,"帐号已经存在或该身份证已经被注册!")
		return
	}
	sqldata:= map[string]string{"username":username,"password":password,"email":email,"entityid":entityid}
	_,err=app.Dblogic.Insert("user",sqldata);
	if err!=nil{
		fmt.Println("register err",err)
		message(w,r,"注册失败！")
		return
	}
	message(w,r,"注册成功！")
}
func(app *Application)LoginDo(w http.ResponseWriter,r *http.Request){
	islogin,_:=app.Verifier.Verify(VERLOGINED,w,r)
	if islogin{
		message(w,r,"您已经登录过了")
		return
	}
	r.ParseForm()
	username:=r.PostFormValue("username")
	password:=r.PostFormValue("password")
	if username==""||password==""{
		message(w,r,"发送信息有误！")
		return
	}
	values,err:=app.Dblogic.Query("user","password,entityid,isadmin","where username='"+username+"'")
	if err!=nil {
		fmt.Println(err)
		return
	}
	datapassword:=values[0]["password"]
	entityid:=values[0]["entityid"]
	isadmin:=values[0]["isadmin"]
	if password!=datapassword{
		message(w,r,"密码错误")
		return
	}
	message(w,r,"登录成功")
	session:=GlobalSessions.SessionStart(w,r)
	session.Set("username",username)
	session.Set("isadmin",isadmin)
	session.Set("entityid",entityid)
}
func(app *Application)AddEduDo(w http.ResponseWriter,r *http.Request){
	islogin,_:=app.Verifier.Verify(VERLOGINED,w,r)
	if !islogin{
		message(w,r,"你还没有登录!")
		return
	}
	isadmin,_:=app.Verifier.Verify(VERISADMIN,w,r)
	if !isadmin{
		message(w,r,"你不是管理员哦!")
		return
	}
	r.ParseForm()
	values:=r.PostForm
	edu:=service.Education{}
	//寻址找值
	rValue:=reflect.ValueOf(&edu).Elem()
	if len(values)!=rValue.NumField()-2{
		message(w,r,"信息不完整!")
		return
	}
	for k,v:=range values{
		if v[0]==""{
			message(w,r,"信息不完整!")
			return
		}
		value:=rValue.FieldByName(k)
		value.SetString(v[0])
	}
	fmt.Println(edu.Photo)
	err:=app.Fabric.AddEdu(&edu)
	if err !=nil{
		fmt.Println(err)
		return
	}
	message(w,r,"信息添加成功!")
}
func(app *Application)QueryDo(w http.ResponseWriter,r *http.Request) {
	islogin, username := app.Verifier.Verify(VERLOGINED, w, r)
	if !islogin {
		message(w, r, "你还没有登录!")
		return
	}
	r.ParseForm()
	entityid:=r.PostFormValue("entityid")
	if entityid==""{
		message(w, r, "信息不完整!")
		return
	}
	edu,err:=app.Fabric.QueryEduInfoByEntityID(entityid)
	if err != nil {
		fmt.Println(err)
		message(w,r,"查询失败可能信息不存在")
		return
	}
	data := struct {
		Islogin  bool
		Username string
		Edu *service.Education
	}{islogin, username.(string),edu}
	showView(w, r, "edu.html", data)
}
func(app *Application)MyEdu(w http.ResponseWriter,r *http.Request) {
	islogin, username := app.Verifier.Verify(VERLOGINED, w, r)
	if !islogin {
		message(w, r, "你还没有登录!")
		return
	}
	session:=GlobalSessions.SessionStart(w,r)
	entityid:=session.Get("entityid")
	edu,err:=app.Fabric.QueryEduInfoByEntityID(entityid.(string))
	if err != nil {
		fmt.Println(err)
		message(w,r,"查询失败可能信息不存在")
		return
	}
	data := struct {
		Islogin  bool
		Username string
		Edu *service.Education
	}{islogin, username.(string),edu}
	showView(w, r, "edu.html", data)
}

func(app *Application) UploadFile(w http.ResponseWriter, r *http.Request){
	islogin,_:= app.Verifier.Verify(VERLOGINED, w, r)
	if !islogin {
		message(w, r, "你还没有登录!")
		return
	}
	isadmin,_:=app.Verifier.Verify(VERISADMIN,w,r)
	if !isadmin{
		message(w,r,"你不是管理员哦!")
		return
	}
	if r.Method!="POST"{
		message(w,r,"请求失败！")
		return
	}
	start:="{"
	content:=""
	end:="}"
	file,_,err:=r.FormFile("file")
	defer file.Close()
	if err != nil {
		content="\"error\":1,\"result\":{\"msg\":\"文件上传失败!\",\"path\":\"\"}"
		w.Write([]byte(start+content+end))
		return
	}
	fileBytes,err:=ioutil.ReadAll(file)
	if err != nil {
		content="\"error\":1,\"result\":{\"msg\":\"文件读取失败!\",\"path\":\"\"}"
		w.Write([]byte(start+content+end))
		return
	}
	if len(fileBytes)>1000000{
		content="\"error\":1,\"result\":{\"msg\":\"文件太大!\",\"path\":\"\"}"
		w.Write([]byte(start+content+end))
		return
	}
	filetype:=http.DetectContentType(fileBytes)
	switch filetype {
	case "image/jpeg", "image/jpg":
	case "image/gif", "image/png":
	case "application/pdf":
		break
	default:
		content = "\"error\":1,\"result\":{\"msg\":\"文件类型错误\",\"path\":\"\"}"
		w.Write([]byte(start + content + end))
		return
	}
	fileName:=randToken(12)
	fileEndings,err:=mime.ExtensionsByType(filetype)
	newPath:=filepath.Join("web","static","photo",fileName+fileEndings[0])
	newFile,err:=os.Create(newPath)
	if err != nil {
		fmt.Println("创建文件失败：" + err.Error())
		content = "\"error\":1,\"result\":{\"msg\":\"创建文件失败\",\"path\":\"\"}"
		w.Write([]byte(start + content + end))
		return
	}
	defer newFile.Close()
	if _, err := newFile.Write(fileBytes); err != nil || newFile.Close() != nil {
		log.Println("写入文件失败：" + err.Error())
		content = "\"error\":1,\"result\":{\"msg\":\"保存文件内容失败\",\"path\":\"\"}"
		w.Write([]byte(start + content + end))
		return
	}
	path := "/static/photo/" + fileName + fileEndings[0]
	content = "\"error\":0,\"result\":{\"fileType\":\"image/png\",\"path\":\"" + path + "\",\"fileName\":\"ce73ac68d0d93de80d925b5a.png\"}"
	w.Write([]byte(start + content + end))
	return
}

func randToken(len int)string{
	b:=make([]byte,len)
	rand.Read(b)
	return fmt.Sprintf("%x",b)
}

func(app *Application)LoginOut(w http.ResponseWriter,r *http.Request){
	islogin,_:= app.Verifier.Verify(VERLOGINED, w, r)
	if !islogin {
		message(w, r, "你还没有登录!")
		return
	}
	GlobalSessions.SessionDestroy(w,r)
	message(w,r,"退出成功!")
}
func init(){
	GlobalSessions ,_= session.NewManager("memory","gosessionid",3600)
	GlobalSessions.GC()
}
