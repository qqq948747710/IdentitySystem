package memory

import (
	"container/list"
	"fmt"
	"github.com/IdentitySystem/web/session"

	"sync"
	"time"
)
var pdef=&Provider{list:list.New()}
type SessionStore struct {
	sid string
	timeAccessed time.Time
	value map[interface{}]interface{}
}

func(st *SessionStore)Set(key,value interface{})error{
	st.value[key]=value
	pdef.SessionUpdate(st.sid)
	return nil
}

func(st *SessionStore)Get(key interface{})interface{}{
	pdef.SessionUpdate(st.sid)
	if v,ok:=st.value[key];ok{
		return v
	}else {
		return nil
	}
	return nil
}

func(st *SessionStore)Delete(key interface{})error{
	if _,ok:=st.value[key];ok{
		delete(st.value,key)
	}else {
		return fmt.Errorf(st.sid+"session is not key")
	}
	return nil
}
func (st *SessionStore)SessionID()string{
	return st.sid
}

type Provider struct {
	lock sync.Mutex
	sessions map[string]*list.Element
	list *list.List
}

func (pder *Provider)SessionInit(sid string)(session.Session,error){
	pder.lock.Lock()
	defer pder.lock.Unlock()
	v:=make(map[interface{}]interface{},0)
	newsess:=&SessionStore{sid:sid,timeAccessed:time.Now(),value:v}
	element:=pder.list.PushBack(newsess)
	pder.sessions[sid]=element
	return newsess,nil
}
func (pder *Provider)SessionRead(sid string)(session.Session,error){
   if elememt,ok:=pder.sessions[sid];ok{
   	 return elememt.Value.(*SessionStore),nil
   }else {
   	 sess,err:=pder.SessionInit(sid)
   	 return sess,err
   }
   return nil,nil
}
func (pder *Provider)SessionDestroy(sid string)(error){
	if element,ok:=pder.sessions[sid];ok{
		delete(pder.sessions,sid)
		pder.list.Remove(element)
		return nil
	}
	return nil
}
func (pder *Provider)SessionGC(maxLifeTime int64){
	pder.lock.Lock()
	defer pder.lock.Unlock()
	for{
		element:=pder.list.Back()
		if element==nil{
			break
		}
		if(element.Value.(*SessionStore).timeAccessed.Unix()+maxLifeTime)<time.Now().Unix(){
			pder.list.Remove(element)
			delete(pder.sessions,element.Value.(*SessionStore).sid)
		}else {
			break
		}
	}
}
func(pder *Provider)SessionUpdate(sid string)error{
	pder.lock.Lock()
	defer pder.lock.Unlock()
	if element,ok:=pder.sessions[sid];ok{
		element.Value.(*SessionStore).timeAccessed=time.Now()
		pder.list.MoveToFront(element)
		return nil
	}
	return nil
}
func init() {
	pdef.sessions=make(map[string]*list.Element,0)
	session.RegisterProvider("memory",pdef)
	fmt.Println("memory init")
}
