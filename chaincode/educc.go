package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

const DOC_TYPE = "eduObj"
func PutEdu(stub shim.ChaincodeStubInterface,item *Education)error{
	item.ObjectType=DOC_TYPE
	EntityID:=item.EntityID
	bytes,err:=json.Marshal(item)
	if err!=nil{
		return fmt.Errorf("chaincode error json error:",err)
	}
	err=stub.PutState(EntityID,bytes)
	if err!=nil{
		return fmt.Errorf("chaincode error 信息录入失败:",err)
	}
	return nil
}


func GetEduInfo(stub shim.ChaincodeStubInterface,entityID string)(*Education,error){
	bytes,err:=stub.GetState(entityID)
	if err != nil {
		return nil,fmt.Errorf(err.Error())
	}
	if bytes==nil{
		return nil,nil
	}
	var item *Education
	err=json.Unmarshal(bytes,&item)
	if err!=nil{
		return nil,fmt.Errorf("chaincode error 获取数据失败",err)
	}
	return item,nil
}

func GetEduByQueryString(stub shim.ChaincodeStubInterface,arg string)([]Education,error){
	resultsIterator,err:=stub.GetQueryResult(arg)
	if err!=nil{
		return nil,fmt.Errorf("chaincode err 富查询失败:",err)
	}
	defer resultsIterator.Close()
	var buffer bytes.Buffer
	buffer.WriteString("[")
	isnotfirst:=false
	for resultsIterator.HasNext(){
		queryResponse,err:=resultsIterator.Next()
		if err!=nil{
			return nil,fmt.Errorf("获取富查询结果失败",err)
		}
		if isnotfirst {
			buffer.WriteString(",")
		}
		buffer.Write(queryResponse.Value)
		isnotfirst=true
	}
	buffer.WriteString("]")
	var items []Education
	err=json.Unmarshal(buffer.Bytes(),&items)
	if err!=nil{
		return nil,fmt.Errorf("富查询结果构建失败",err)
	}
	return items,nil
}