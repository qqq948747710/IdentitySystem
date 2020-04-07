package main

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

type EducatuonCC struct {

}

func (t *EducatuonCC)Init(stub shim.ChaincodeStubInterface)peer.Response{
	return shim.Success(nil)
}

func(t *EducatuonCC)Invoke(stub shim.ChaincodeStubInterface)peer.Response{
	fun,args:=stub.GetFunctionAndParameters()
	if fun =="addEdu"{
		return t.addEdu(stub,args)
	}else if fun=="queryEduByCerNoAndName"{
		return t.queryEduByCerNoAndName(stub,args)
	}else if fun=="queryEduInfoByEntityID"{
		return t.queryEduInfoByEntityID(stub,args)
	}else if fun=="updateEdu"{
		return t.updateEdu(stub,args)
	}else if fun=="delEdu"{
		return t.delEdu(stub,args)
	}
	return shim.Error("chaincoe error:is not this func,you can try addEdu...")
}

//arg Education obj
//key 身份证号 val obj
func(t *EducatuonCC)addEdu(stub shim.ChaincodeStubInterface,args []string)peer.Response{
	if len(args)!=1{
		return shim.Error("给定参数不对,需要1个")
	}
	var item *Education
	err:=json.Unmarshal([]byte(args[0]),&item)
	if err!=nil{
		return shim.Error("chaincode main err:this json is error:"+err.Error())
	}
	info,err:=GetEduInfo(stub,item.EntityID)
	if info!=nil{
		return shim.Error("chaincode main err:this identity is has")
	}
	err=PutEdu(stub,item)
	if err!=nil{
		return shim.Error("chaincode main putdata error:"+err.Error())
	}
	return shim.Success([]byte(item.EntityID))
}

func(t *EducatuonCC)queryEduByCerNoAndName(stub shim.ChaincodeStubInterface,args []string)peer.Response{
	if len(args)!=2{
		return shim.Error("给定参数不对,需要2个")
	}
	CertNo:=args[0]
	Name:=args[1]

	querString:=fmt.Sprintf("{\"selector\":{\"docType\":\"%s\",\"CertNo\":\"%s\",\"Name\":\"%s\"}}",DOC_TYPE,CertNo,Name)
	result,err:=GetEduByQueryString(stub,querString)
	if err!=nil{
		return shim.Error(err.Error())
	}
	if len(result)==0{
		return shim.Error("is not this entityid,pules addEdu")
	}
	bytes,err:=json.Marshal(result)
	if err!=nil{
		return shim.Error("json is err"+err.Error())
	}
	return shim.Success(bytes)
}

func(t *EducatuonCC)queryEduInfoByEntityID(stub shim.ChaincodeStubInterface,args []string)peer.Response{
	if len(args)!=1{
		return shim.Error("给定参数不对,需要1个")
	}
	edu,err:=GetEduInfo(stub,args[0])
	if err != nil {
		return shim.Error(err.Error())
	}
	if edu==nil{
		return shim.Error("is not this entityid,pules addEdu")
	}
	iterator,err:=stub.GetHistoryForKey(edu.EntityID)
	defer iterator.Close()
	if err != nil {
		return shim.Error("指定身份证号获取历史记录失败"+err.Error())
	}
	if iterator==nil{
		return shim.Error("no history iterator"+err.Error())
	}
	var hisEdus []HistoryItem
	for iterator.HasNext(){
		hisData,err:=iterator.Next()
		if err != nil {
			return shim.Error("身份对象获取失败"+err.Error())
		}
		var hisEdu HistoryItem
		hisEdu.TxId=hisData.TxId
		var Edu Education
		err=json.Unmarshal(hisData.Value,&Edu)
		if err!=nil{
			return shim.Error("json对象结构错误"+err.Error())
		}
		hisEdu.Education=Edu
		hisEdus=append(hisEdus,hisEdu)
	}
	edu.Historys=hisEdus
	edubyte,err:=json.Marshal(edu)
	if err!=nil{
		return shim.Error("json对象结构错误"+err.Error())
	}
	return shim.Success(edubyte)
}

func(t *EducatuonCC)updateEdu(stub shim.ChaincodeStubInterface,args []string)peer.Response{
	if len(args)!=2{
		return shim.Error("给定参数不对,需要2个")
	}
	var item *Education
	err:=json.Unmarshal([]byte(args[1]),&item)
	if err!=nil{
		return shim.Error("json对象结构错误"+err.Error())
	}
	bytes,err:=GetEduInfo(stub,args[0])
	if bytes==nil{
		return shim.Error("并没有该身份的用户信息请先addEdu添加")
	}
	err=PutEdu(stub,item)
	if err!=nil{
		return shim.Error("chaincode main putdata error:"+err.Error())
	}
	return shim.Success([]byte(args[0]))
}

func (t *EducatuonCC)delEdu(stub shim.ChaincodeStubInterface,args []string)peer.Response{
	if len(args) != 1 {
		return shim.Error("给定参数不对,需要1个")
	}
	err:=stub.DelState(args[0])
	if err != nil {
		return shim.Error("删除身份信息失败"+err.Error())
	}
	return shim.Success([]byte("删除成功"))
}

func main()  {
	err:=shim.Start(new(EducatuonCC))
	if err!=nil{
		fmt.Println("启动EducationChaincode时发生错误")
	}
}