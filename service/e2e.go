package service

import (
	"fmt"
	"github.com/IdentitySystem/sdkinit"
	"os"
)

func E2e()*SrviceSetup{
	const ConfigPath="config.yaml"
	sdk,err:=sdkinit.SetupSdk(ConfigPath,false)
	if err!=nil{
		fmt.Println(err.Error())
	}
	initInfo := &sdkinit.Sdkinfo{

		ChannelID: "identitychannel",
		ChannelPath: "./fixtures/channel-artifacts/channel.tx",

		Admin:"Admin",
		OrgName:"orgidentity",
		OrderName: "orderer.identity.com",
		User:"User1",

		ChaincodeID:"identitycc",
		Chaincodepath:"github.com/IdentitySystem/chaincode",
		ChaincodeGopath:os.Getenv("GOPATH"),
	}
	err=sdkinit.CreateChannel(sdk,initInfo)
	if err!=nil{
		fmt.Println(err.Error())
	}
	//initInfo.Resclient,err=sdkinit.Getresmag(sdk,initInfo)
	client,err:=sdkinit.InstallAndInstantiateCC(sdk,initInfo)
	//client,err:=sdkinit.Getchannel(sdk,initInfo)
	if err != nil {
		fmt.Println(err.Error())
	}
	ssetup:=&SrviceSetup{ChainCodeID:initInfo.ChaincodeID,Client:client}
	edu := &Education{
		Name: "张三",
		Gender: "男",
		Nation: "汉",
		EntityID: "101",
		Place: "北京",
		BirthDay: "1991年01月01日",
		EnrollDate: "2009年9月",
		GraduationDate: "2013年7月",
		SchoolName: "中国政法大学",
		Major: "社会学",
		QuaType: "普通",
		Length: "四年",
		Mode: "普通全日制",
		Level: "本科",
		Graduation: "毕业",
		CertNo: "111",
		Photo: "/static/phone/11.png",
	}
	fmt.Println("AddEdu:")
	ssetup.AddEdu(edu)
	edus, err:= ssetup.QueryEduByCertNoAndName(edu.CertNo, edu.Name)
	if err!=nil{
		fmt.Println(err)
	}
	fmt.Println(edus[0].Name)
	fmt.Println("UpdateEdu:")
	edu.Name="里斯"
	err=ssetup.UpdateEdu(edu.EntityID,edu)
	if err!=nil{
		fmt.Println(err.Error())
	}
	fmt.Println("QueryEduInfoByEntityID:")
	Edu,err:=ssetup.QueryEduInfoByEntityID(edu.EntityID)
	fmt.Println(Edu.Name)
	if err!=nil{
		fmt.Println(err.Error())
	}
	fmt.Println("DelEdu:")
	err=ssetup.DelEdu(edu.EntityID)
	if err!=nil{
		fmt.Println(err.Error())
	}
	edu,err=ssetup.QueryEduInfoByEntityID(edu.EntityID)
	if err!=nil{
		fmt.Println(err.Error())
	}
	return ssetup
}
