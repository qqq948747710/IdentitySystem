package sdkinit

import "github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"

type Sdkinfo struct {
	ChaincodeID string
	Chaincodepath string
	ChaincodeGopath string

	ChannelID string
	ChannelPath string

	OrgName string
	User string
	Admin string
	Resclient *resmgmt.Client

	OrderName string
}