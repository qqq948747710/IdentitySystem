package sdkinit

import (
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	mspclient "github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab/ccpackager/gopackager"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/common/cauthdsl"
)

const ChainCVersion  ="1.0"
func SetupSdk(Configpath string,inited bool)(*fabsdk.FabricSDK,error) {
	if inited {
		return nil,fmt.Errorf("sdk已经实例化")
	}
	sdk,err:=fabsdk.New(config.FromFile(Configpath))
	if err!=nil{
		return nil,fmt.Errorf("sdk实例化失败:",err)
	}
	fmt.Println("sdk实例化成功")
	return sdk,nil
}

func Getresmag(sdk *fabsdk.FabricSDK,info *Sdkinfo)(*resmgmt.Client,error){
	clientContext:=sdk.Context(fabsdk.WithUser(info.Admin),fabsdk.WithOrg(info.OrgName))
	resmgmtclient,err:=resmgmt.New(clientContext)
	if err!=nil{
		return nil,fmt.Errorf("获取管理客户端失败",err)
	}
	return resmgmtclient,nil
}

func Getchannel(sdk *fabsdk.FabricSDK,info *Sdkinfo)(*channel.Client,error){
	channelClient,err:=channel.New(sdk.ChannelContext(info.ChannelID,fabsdk.WithUser(info.User),fabsdk.WithOrg(info.OrgName)))
	if err!=nil{
		return nil,fmt.Errorf("通道操作客户端创建失败!",err)
	}
	return channelClient,nil
}
func CreateChannel(sdk *fabsdk.FabricSDK,info *Sdkinfo)error{
	fmt.Println("开始创建通道")
	clientContext:=sdk.Context(fabsdk.WithUser(info.Admin),fabsdk.WithOrg(info.OrgName))
	resmgmtclient,err:=resmgmt.New(clientContext)
	if err!=nil{
		return fmt.Errorf("实例化管理客户端出错",err)
	}
	mspclient,err:=mspclient.New(sdk.Context(),mspclient.WithOrg(info.OrgName))
	if err!=nil{
		return fmt.Errorf("用户管理客户端初始化失败",err)
	}
	adminidentity,err:=mspclient.GetSigningIdentity(info.Admin)
	if err!=nil{
		return fmt.Errorf("获取管理员签名标识失败",err)
	}
	reqchannel:=resmgmt.SaveChannelRequest{ChannelID:info.ChannelID,ChannelConfigPath:info.ChannelPath,SigningIdentities:[]msp.SigningIdentity{adminidentity}}
	_,err=resmgmtclient.SaveChannel(reqchannel,resmgmt.WithRetry(retry.DefaultResMgmtOpts),resmgmt.WithOrdererEndpoint(info.OrderName))
	if err!=nil{
		return fmt.Errorf("新建通道失败",err)
	}
	fmt.Println("创建通道成功")
	info.Resclient=resmgmtclient
	fmt.Println("把节点加入通道...")


	err=info.Resclient.JoinChannel(info.ChannelID,resmgmt.WithRetry(retry.DefaultResMgmtOpts),resmgmt.WithOrdererEndpoint(info.OrderName))
	if err!=nil{
		return fmt.Errorf("通道加入失败",err)
	}
	fmt.Println("节点加入成功")
	return nil

}

func InstallAndInstantiateCC(sdk *fabsdk.FabricSDK,info *Sdkinfo)(*channel.Client,error){
	fmt.Println("开始安装链码！！！")
	ccPkg,err:=gopackager.NewCCPackage(info.Chaincodepath,info.ChaincodeGopath)
	if err!=nil{
		return nil,fmt.Errorf("ccpkg初始化失败",err)
	}
	req:=resmgmt.InstallCCRequest{Name:info.ChaincodeID,Path:info.Chaincodepath,Version:ChainCVersion,Package:ccPkg}
	_,err=info.Resclient.InstallCC(req,resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	if err!=nil{
		return nil,fmt.Errorf("链码安装失败",err)
	}
	fmt.Println("链码安装成功！")
	fmt.Println("链码开始实例化!")
	policy:=cauthdsl.SignedByAnyMember([]string{"OrgIdentityMSP"})
	instantiateCCReq:=resmgmt.InstantiateCCRequest{Name:info.ChaincodeID,Path:info.Chaincodepath,Version:ChainCVersion,Args:[][]byte{[]byte("init")},Policy:policy}
	_,err=info.Resclient.InstantiateCC(info.ChannelID,instantiateCCReq,resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	if err!=nil{
		return nil,fmt.Errorf("链码实例化失败",err)
	}
	fmt.Println("链码实例化成功！")
	channelClient,err:=channel.New(sdk.ChannelContext(info.ChannelID,fabsdk.WithUser(info.User),fabsdk.WithOrg(info.OrgName)))
	if err!=nil{
		return nil,fmt.Errorf("通道操作客户端创建失败!",err)
	}
	return channelClient,nil
}