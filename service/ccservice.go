package service

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
)

const(
	addEdu string="addEdu"
	queryEduByCerNoAndName string="queryEduByCerNoAndName"
	queryEduInfoByEntityID string="queryEduInfoByEntityID"
	updateEdu string="updateEdu"
	delEdu string="delEdu"
)

type SrviceSetup struct {
	ChainCodeID string
	Client *channel.Client
}

func(s *SrviceSetup)AddEdu(edu *Education)error{
	bytes,err:=json.Marshal(edu)
	if err!=nil{
		return fmt.Errorf("json is err!",err)
	}
	channelreq:=channel.Request{ChaincodeID:s.ChainCodeID,Fcn:addEdu,Args:[][]byte{[]byte(bytes)}}
	respone,err:=s.Client.Execute(channelreq,channel.WithRetry(retry.DefaultResMgmtOpts))
	if err!=nil{
		return fmt.Errorf("插入身份证信息错误",err)
	}
	fmt.Println("身份证信息插入成功!交易id:"+respone.TransactionID)
	return nil
}

func(s *SrviceSetup)QueryEduByCertNoAndName(certno string,name string)([]Education,error){
	channelreq:=channel.Request{ChaincodeID:s.ChainCodeID,Fcn:queryEduByCerNoAndName,Args:[][]byte{[]byte(certno),[]byte(name)}}
	respone, err := s.Client.Execute(channelreq)
	if err != nil {
		return nil,fmt.Errorf("QueryEduByCertNoAndName查询失败",err)
	}
	var edus []Education
	err=json.Unmarshal(respone.Payload,&edus)
	if err!=nil{
		return nil,fmt.Errorf("json解析失败!",err)
	}
	return edus,nil
}

func (s *SrviceSetup)QueryEduInfoByEntityID(entityid string)(*Education,error){
	channelreq:=channel.Request{ChaincodeID:s.ChainCodeID,Fcn:queryEduInfoByEntityID,Args:[][]byte{[]byte(entityid)}}
	respone,err:=s.Client.Execute(channelreq)
	if err != nil {
		return nil,fmt.Errorf("QueryEduInfoByEntityID查询失败",err)
	}
	var edu *Education
	err=json.Unmarshal(respone.Payload,&edu)
	if err != nil {
		return nil,fmt.Errorf("json解析失败!",err)
	}
	return edu,err
}

func (s *SrviceSetup)UpdateEdu(entityid string,edu *Education)error{
	bytes,err:=json.Marshal(edu)
	if err != nil {
		return fmt.Errorf("json解析失败!",err)
	}
	channelreq:=channel.Request{ChaincodeID:s.ChainCodeID,Fcn:updateEdu,Args:[][]byte{[]byte(entityid),bytes}}
	_,err=s.Client.Execute(channelreq)
	if err != nil {
		return fmt.Errorf("update失败",err)
	}
	return nil
}

func (s *SrviceSetup)DelEdu(entityid string)error{
	channelreq:=channel.Request{ChaincodeID:s.ChainCodeID,Fcn:delEdu,Args:[][]byte{[]byte(entityid)}}
	_,err:=s.Client.Execute(channelreq)
	if err != nil {
		return fmt.Errorf("delete失败",err)
	}
	return nil
}