package template

import (
	"github.com/wayhood/gogetui/protobuf"
	"github.com/golang/protobuf/proto"
	"fmt"
	"time"
	"strconv"
	"strings"
	apnPayload "github.com/wayhood/gogetui/payload"
)

type BaseTemplate struct {
	AppKey   string
	AppId    string
	PushInfo interface{} //protobuf.PushInfo
	Duration string
}

func (bt *BaseTemplate) GetPushInfo() *protobuf.PushInfo {

	if bt.PushInfo == nil {
		pushInfo := &protobuf.PushInfo{
			Message:   proto.String(""),
			ActionKey: proto.String(""),
			Sound:     proto.String(""),
			Badge:     proto.String(""),
		}
		return pushInfo
	}
	pushInfo := new(protobuf.PushInfo)
	*pushInfo = bt.PushInfo.(protobuf.PushInfo)
	return pushInfo

}


func (bt *BaseTemplate) SetPushInfo(actionLocKey string, badge string, message string,
	sound string, payload string, locKey string, locArgs string, launchImage string, contentAvailable int32) {

	PushInfo := protobuf.PushInfo{
		ActionLocKey :proto.String(actionLocKey),
		Badge:     proto.String(badge),
		Message:   proto.String(message),
		Sound:     proto.String(sound),
		Payload:  proto.String(payload),
		LocKey:  proto.String(locKey),
		LocArgs:  proto.String(locArgs),
		LaunchImage:  proto.String(launchImage),
		ContentAvailable:  proto.Int32(contentAvailable),
	}

	bt.PushInfo = PushInfo
	//fmt.Println("%s",bt.PushInfo)
	l := len(bt.ProcessPayload(actionLocKey, badge, message, sound,
		payload, locKey, locArgs, launchImage, contentAvailable))
	if l > 512 {
		fmt.Println("PushInfo length over limit: " + "%d" + ". Allowed: 256.", l)
		panic("")
	}

}


func (bt *BaseTemplate) SetDuration(begin string, end string) {
	t1, _ := time.Parse("2006-01-02 15:04:05", begin)
	t2, _ := time.Parse("2006-01-02 15:04:05", end)
	s1 := (t1.Unix() - 28800) * 1000
	s2 := (t2.Unix() - 28800) * 1000
	if s1 > 0 && s2 > 0 && s2 >= s1 {
		bt.Duration = strconv.FormatInt(s1, 10) + "-" + strconv.FormatInt(s2, 10)
	}else if s1 > s2 {
		panic("startTime should be smaller than endTime")
	}else {
		panic("DateFormat: yyyy-MM-dd HH:mm:ss")
	}
}

func (bt *BaseTemplate) GetDuration() string {
	return bt.Duration
}

func (bt *BaseTemplate) GetDurCondition() []string {
	Du := []string{"duration=" + bt.GetDuration()}
	// fmt.Println(Du)
	return Du
}


func (bt *BaseTemplate) ProcessPayload(actionLocKey string, badge string, message string, sound string,
	payload string, locKey string, locArgs string, launchImage string, contentAvailable int32) string {
	isValid := false
	pb := new(apnPayload.APNPayload)
	if len(locKey) > 0 {
		pb.AlertLocKey = locKey
		if len(locArgs) > 0 {
			pb.AlertLocArgs = strings.Split(locArgs, ",")
		}
		isValid = true
	}
	if len(message) > 0 {
		pb.AlertBody = message
		isValid = true
	}
	if len(actionLocKey) > 0 {
		pb.AlertActionLocKey = actionLocKey
	}
	if len(launchImage) > 0 {
		pb.AlertLaunchImage = launchImage
	}
	badgeNum, _ := strconv.Atoi(badge)
	if badgeNum >= 0 {
		pb.Badge = strconv.Itoa(badgeNum)
		isValid = true
	}
	if len(sound) > 0 {
		pb.Sound = sound
	}else {
		pb.Sound = "default"
	}
	if len(payload) > 0 {
		pb.AddParam("payload", payload)
	}
	if contentAvailable == 1 {
		pb.ContentAvailable = 1
		isValid = true
	}
	if isValid == false {
		fmt.Println("one of the params(locKey,message,badge) must not be null or contentAvailable must be 1")
		panic("")
	}
	jsons := pb.ToString()
	return jsons

}
