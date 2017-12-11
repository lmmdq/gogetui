package igetui

import "github.com/wayhood/gogetui/igetui/template"

type IGtMessage struct {
	IsOffline         bool
	OfflineExpireTime int32
	Data              template.ITemplate
	PushNetWorkType  byte
	Priority     int32
	
}

type IGtSingleMessage struct {
	IGtMessage
	
}

func NewIGtSingleMessage(isoffline bool, offlineexpiretime int32, templatee template.ITemplate) *IGtSingleMessage {
	return &IGtSingleMessage{
		IGtMessage: IGtMessage{
			IsOffline:         isoffline,
			OfflineExpireTime: offlineexpiretime,
			Data:              templatee,	
		},
	}
}

type IGtListMessage struct {
	IGtMessage
}

func NewIGtListMessage(isoffline bool, offlineexpiretime int32, templatee template.ITemplate) *IGtListMessage {
	return &IGtListMessage{
		IGtMessage: IGtMessage{
			IsOffline:         isoffline,
			OfflineExpireTime: offlineexpiretime,
			Data:              templatee,
		},
	}
}

type IGtAppMessage struct {
	IGtMessage	
	Speed  int32
	TagList		[]string
	AppIdList     []string
	PhoneTypeList []string
	ProvinceList  []string
	Conditions	[]interface{}
}

func NewIGtAppMessage(isoffline bool, offlineexpiretime int32, templatee template.ITemplate) *IGtAppMessage {
	return &IGtAppMessage{
		IGtMessage: IGtMessage{
			IsOffline:         isoffline,
			OfflineExpireTime: offlineexpiretime,
			Data:              templatee,
		},
		Speed:	0,
	}
}






