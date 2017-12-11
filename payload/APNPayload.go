package payload

import (
	"strings"
	"fmt"
	"encoding/json"
)

const APS string = "aps"

type APNPayload struct {
	Params            map[string]interface{}
	Alert             string
	Badge             string
	Sound             string
	AlertBody         string
	AlertActionLocKey string
	AlertLocKey       string
	AlertLocArgs      []string
	AlertLaunchImage  string
	ContentAvailable  int32
}

func (p *APNPayload) AddParam(key string, obj interface{}) {
	if p.Params == nil {
		p.Params = map[string]interface{}{}
	}
	if strings.EqualFold(APS, key) {
		fmt.Printf("the key can't be aps")
	}else {
		p.Params[key] = obj
	}
}

func (p *APNPayload) PutIntoJson(key string, value interface{}, obj map[string]interface{}) {
	if value != nil {
		obj[key] = value
	}
}

func (p *APNPayload) ToString() string {
	objectt := map[string]interface{}{}
	ApsObj := map[string]interface{}{}

	if len(p.Alert) > 0 {
		ApsObj["alert"] = p.Alert
	}else {
		if len(p.AlertBody) > 0 || len(p.AlertLocKey) > 0 {
			alertObj := map[string]interface{}{}
			p.PutIntoJson("body", p.AlertBody, alertObj)
			p.PutIntoJson("action-loc-key", p.AlertActionLocKey, alertObj)
			p.PutIntoJson("loc-key", p.AlertLocKey, alertObj)
			p.PutIntoJson("launch-image", p.AlertLaunchImage, alertObj)
			if p.AlertLocArgs != nil {
				alertObj["loc-args"] = p.AlertLocArgs
			}
			ApsObj["alert"] = alertObj
		}
	}
	if len(p.Badge) >= 0 {
		ApsObj["badge"] = p.Badge
	}
	if p.Sound != "com.gexin.ios.silence" {
		p.PutIntoJson("sound", p.Sound, ApsObj)
	}
	if p.ContentAvailable == 1 {
		ApsObj["content-available"] = 1
		objectt[APS] = ApsObj
	}
	if p.Params != nil {
		for k, v := range p.Params {
			objectt[k] = v
		}
	}
	datajson, err := json.Marshal(objectt)
	if err != nil {
		fmt.Println("error:", err)
		panic("")
	}
	return string(datajson)
}

