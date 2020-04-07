package session

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type Provider interface {
	SessionInit(sid string)(Session,error)
	SessionRead(sid string)(Session,error)
	SessionDestroy(sid string)(error)
	SessionGC(maxLifeTime int64)
}

type Session interface {
	Set(key,value interface{})error
	Get(key interface{})interface{}
	Delete(key interface{})error
	SessionID()string
}

var providers=make(map[string]Provider)

func RegisterProvider(name string,provider Provider){
	if provider==nil{
		panic("seesion:Register provider is nil")
	}
	if _,e:=providers[name];e{
		panic("session:Register provider is existed")
	}
	providers[name]=provider
}

type Manager struct {
	cookieName string
	lock sync.Mutex
	provider Provider
	maxLifeTime int64
}

func (manager *Manager)sessionId()string {
	b:=make([]byte,32)
	if _,err:=io.ReadFull(rand.Reader,b);err!=nil{
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}
func (manager *Manager)SessionStart(w http.ResponseWriter,r *http.Request)(session Session){
	manager.lock.Lock()
	defer manager.lock.Unlock()
	cookie,err:=r.Cookie(manager.cookieName)
	if err !=nil || cookie.Value==""{
		sid:=manager.sessionId()
		session,_=manager.provider.SessionInit(sid)
		cookie:=http.Cookie{Name:manager.cookieName,Value:url.QueryEscape(sid),Path:"/",HttpOnly:true,MaxAge:int(manager.maxLifeTime)}
		http.SetCookie(w,&cookie)
	}else {
		sid,_:=url.QueryUnescape(cookie.Value)
		session,_=manager.provider.SessionRead(sid)
	}
	return
}

func(manager *Manager)SessionDestroy(w http.ResponseWriter,r *http.Request){
	cookie,err:=r.Cookie(manager.cookieName)
	if err!=nil||cookie.Value==""{
		return
	}else {
		manager.lock.Lock()
		defer manager.lock.Unlock()
		manager.provider.SessionDestroy(cookie.Value)
		expiration:=time.Now()
		cookie:=http.Cookie{Name:manager.cookieName,Path:"/",HttpOnly:true,Expires:expiration}
		http.SetCookie(w,&cookie)
	}
}
func(manager *Manager)GC(){
	manager.lock.Lock()
	defer manager.lock.Unlock()
	manager.provider.SessionGC(manager.maxLifeTime)
	time.AfterFunc(time.Duration(manager.maxLifeTime)*time.Second, func() {
		manager.GC()
	})
}

func NewManager(providerName,cookieName string,maxLifetime int64)(*Manager,error){
	fmt.Println("请务必载入memory包体")
	provider,ok:=providers[providerName]
	if !ok{
		return nil,fmt.Errorf("session:unknow provide %q (forgotten import?)",providerName)
	}
	return &Manager{
		cookieName:cookieName,
		maxLifeTime:maxLifetime,
		provider:provider,
	},nil
}