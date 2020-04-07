package service

import (
	"fmt"
	"github.com/IdentitySystem/sdkinit"
	"os"
)

func ServicStup()*SrviceSetup{
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
	//err=sdkinit.CreateChannel(sdk,initInfo)
	//if err!=nil{
	//	fmt.Println(err.Error())
	//}
	initInfo.Resclient,err=sdkinit.Getresmag(sdk,initInfo)
	//client,err:=sdkinit.InstallAndInstantiateCC(sdk,initInfo)
	client,err:=sdkinit.Getchannel(sdk,initInfo)
	if err != nil {
		fmt.Println(err.Error())
	}
	ssetup:=&SrviceSetup{ChainCodeID:initInfo.ChaincodeID,Client:client}
	return ssetup
}
