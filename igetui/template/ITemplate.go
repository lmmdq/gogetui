package template

import "github.com/wayhood/gogetui/protobuf"

type ITemplate interface {
	GetTransparent() *protobuf.Transparent
	GetActionChains() []*protobuf.ActionChain
	GetPushInfo() *protobuf.PushInfo
	GetTransmissionContent() string
	GetPushType() string
}
