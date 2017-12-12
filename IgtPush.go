package gogetui

import (
	"fmt"
	"time"
	"github.com/golang/protobuf/proto"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"
	"strconv"
	"io/ioutil"
	"crypto/md5"
	"io"
	"encoding/hex"
	"github.com/wayhood/gogetui/igetui"
	"github.com/snluu/uuid"
)

type IGeTui struct{
	serviceMap string
	Host string
	AppKey string
	MasterSecret string
	flag bool
}


func (this *IGeTui) Init(host string, appKey string, masterSecret string) {
	this.AppKey = appKey
	this.MasterSecret = masterSecret

	//if host is not None:
	//	host = host.strip()
	//
	//if ssl is None and host is not None and host != '' and host.lower().startswith('https:'):
	//ssl = True
	//
	//self.useSSL = (ssl if ssl is not None else False);
	//
	//if host is None or len(host) <= 0:
	//self.hosts = GtConfig.getDefaultDomainUrl(self.useSSL)
	//else:else
	//self.hosts = list()
	//self.hosts.append(host)
	this.initOSDomain()
}

func (this *IGeTui) initOSDomain() {
	//hosts = IGeTui.serviceMap.get(this.appKey)
	// 第一次创建时要获取域名列表同时启动检测线程
	//if hosts is None or len(hosts) == 0:
	//hosts = self.getOSPushDomainUrlList()
	//IGeTui.serviceMap[self.appKey] = hosts
	this.FastUrl()
}

func (this *IGeTui) FastUrl() {
	if this.Host ==""{
		ln := this.ConnOSServerHostList()
		if len(ln)==1{
			this.Host=ln[0]
		}
		if len(ln)>1{
			this.gethost(ln)
			for ;this.Host=="";{
				time.Sleep(time.Millisecond)
			}
			go this.timerhost(ln)
		}
	}
}


func (this *IGeTui) test(url string) {
	_,err := http.Get(url)
	if err !=nil{
		return
	}else if this.flag == false{
		this.Host=url
		this.flag = true
	}
}

func (this *IGeTui) gethost(ln []string){
	this.flag =  false
	for i:=0;i <len(ln);i++ {
		go this.test(ln[i])
	}
}

func (this *IGeTui) timerhost(ln []string){
	for {
		time.Sleep(600*1000*time.Millisecond)
		go this.gethost(ln)
	}
}

func (this *IGeTui) ConnOSServerHostList() []string {
	l := this.ConfigOsServerHostList()
	if l == nil || len(l) == 0 {
		ln := []string{"http://sdk.open.api.igexin.com/serviceex",
			"http://sdk.open.api.gepush.com/serviceex",
			"http://sdk.open.api.getui.net/serviceex",
		}
		return ln
	}
	return l
}

func (this *IGeTui) ConfigOsServerHostList() []string {
	url := "http://sdk.open.apilist.igexin.com/os_list"
	response, err := http.Get(url)
	if err != nil {
		return nil
	}
	defer response.Body.Close()
	body, _ := ioutil.ReadAll(response.Body)
	l := strings.Split(string(body), "\r\n")
	var ll []string
	for i := 0; i < len(l); i++ {
		if strings.HasPrefix(l[i], "http") {
			ll = append(ll, l[i])
		}
	}
	return l
}



func (this *IGeTui) connect() bool {
	sign := this.Sign(this.AppKey, this.CurrentTime(), this.MasterSecret)
	params := map[string]interface{}{}
	params["action"] = "connect"
	params["appkey"] = this.AppKey
	params["timeStamp"] = this.CurrentTime()
	params["sign"] = sign

	rep := this.HttpPost(params)
	fmt.Println("rep")
	fmt.Println(rep)
	if "success" == rep["result"] {
		return true
	} else {
		fmt.Println("connect failed")
		panic("connect failed")
	}
	return false
}

func (this *IGeTui) PushMessageToSingle(message igetui.IGtSingleMessage, tartget igetui.Target) map[string]interface{} {
	params := map[string]interface{}{}

	var id uuid.UUID = uuid.Rand()
	params["requestId"] = id.Hex()
	params["action"] = "pushMessageToSingleAction"
	params["appkey"] = this.AppKey
	transparent := message.Data.GetTransparent()
	// fmt.Println(transparent)
	byteArray, _ := proto.Marshal(transparent)
	params["clientData"] = base64.StdEncoding.EncodeToString(byteArray)
	params["transmissionContent"] = message.Data.GetTransmissionContent()
	params["isOffline"] = message.IsOffline
	params["offlineExpireTime"] = message.OfflineExpireTime
	params["appId"] = tartget.AppId
	params["clientId"] = tartget.ClientId
	params["type"] = 2
	params["pushType"] = message.Data.GetPushType()
	//增加pushNetWorkType参数(0:不限;1:wifi;)
	params["pushNetWorkType"] = message.PushNetWorkType
	return this.HttpPostJson(params)

}

func (this *IGeTui) PushMessageToSingleWithRequestId(message igetui.IGtSingleMessage, tartget igetui.Target, requestId string) map[string]interface{} {
	params := map[string]interface{}{}
	
	params["requestId"] = requestId
	params["action"] = "pushMessageToSingleAction"
	params["appkey"] = this.AppKey
	transparent := message.Data.GetTransparent()
	// fmt.Println(transparent)
	byteArray, _ := proto.Marshal(transparent)
	params["clientData"] = base64.StdEncoding.EncodeToString(byteArray)
	params["transmissionContent"] = message.Data.GetTransmissionContent()
	params["isOffline"] = message.IsOffline
	params["offlineExpireTime"] = message.OfflineExpireTime
	params["appId"] = tartget.AppId
	params["clientId"] = tartget.ClientId
	params["type"] = 2
	params["pushType"] = message.Data.GetPushType()
	//增加pushNetWorkType参数(0:不限;1:wifi;)
	params["pushNetWorkType"] = message.PushNetWorkType
	return this.HttpPostJson(params)

}

//appMessage
/*
        params = dict()
        contentId = self.getContentId(message, taskGroupName)
        params['action'] = "pushMessageToAppAction"
        params['appkey'] = self.appKey
        params['contentId'] = contentId
        params['type'] = 2
        return self.httpPostJson(self.host, params)
 */
func (this *IGeTui) PushMessageToApp(message igetui.IGtAppMessage) map[string]interface{} {
	params := map[string]interface{}{}
	contentId := this.ContentIdApp(message)
	params["action"] = "pushMessageToAppAction"
	params["appkey"] = this.AppKey
	params["contentId"] = contentId
	params["type"] = 2
	return this.HttpPostJson(params)
}

func (this *IGeTui) PushMessageToAppWithTaskGroupName(message igetui.IGtAppMessage, taskGroupName string) map[string]interface{} {
	params := map[string]interface{}{}
	contentId := this.ContentIdAppWtihTaskGroupName(message, taskGroupName)
	params["action"] = "pushMessageToAppAction"
	params["appkey"] = this.AppKey
	params["contentId"] = contentId
	params["type"] = 2
	return this.HttpPostJson(params)
}

func (this *IGeTui) PushMessageToList(contentId string, targets []igetui.Target) map[string]interface{} {
	params := map[string]interface{}{}
	params["action"] = "pushMessageToListAction"
	params["appkey"] = this.AppKey
	params["contentId"] = contentId

	targetList := []interface{}{}
	for _, target := range targets {
		appId := target.AppId
		clientId := target.ClientId
		targetTmp := map[string]string{"appId": appId, "clientId": clientId}
		targetList = append(targetList, targetTmp)
	}

	params["targetList"] = targetList
	params["type"] = 2

	return this.HttpPostJson(params)
}

func (this *IGeTui) ContentIdApp(message igetui.IGtAppMessage) interface{} {
	params := map[string]interface{}{}

	params["action"] = "getContentIdAction"
	params["appkey"] = this.AppKey
	transparent := message.Data.GetTransparent()
	byteArray, _ := proto.Marshal(transparent)
	params["clientData"] = base64.StdEncoding.EncodeToString(byteArray)
	params["transmissionContent"] = message.Data.GetTransmissionContent()
	params["isOffline"] = message.IsOffline
	params["offlineExpireTime"] = message.OfflineExpireTime
	params["pushType"] = message.Data.GetPushType()
	// 增加pushNetWorkType参数(0:不限;1:wifi;2:4G/3G/2G)
	params["pushNetWorkType"] = message.PushNetWorkType
	params["tpye"] = 2

	if len(message.Conditions) == 0 {
		params["phoneTypeList"] = message.PhoneTypeList
		params["provinceList"] = message.ProvinceList
		params["tagList"] = message.TagList
	} else {
		conditions := message.Conditions
		params["conditions"] = conditions
	}
	params["speed"] = message.Speed
	params["contentType"] = 2
	params["appIdList"] = message.AppIdList

	ret := this.HttpPostJson(params)

	if ret["result"] == "ok" {
		return ret["contentId"]
	} else {
		panic("获取 contentId 失败：")
	}
}

func (this *IGeTui) ContentIdAppWtihTaskGroupName(message igetui.IGtAppMessage, taskGroupName string) interface{} {
	params := map[string]interface{}{}

	if len(taskGroupName) > 40 {
		panic("TaskGroupName is OverLimit 40")
	}

	params["taskGroupName"] = taskGroupName
	params["action"] = "getContentIdAction"
	params["appkey"] = this.AppKey
	transparent := message.Data.GetTransparent()
	byteArray, _ := proto.Marshal(transparent)
	params["clientData"] = base64.StdEncoding.EncodeToString(byteArray)
	params["transmissionContent"] = message.Data.GetTransmissionContent()
	params["isOffline"] = message.IsOffline
	params["offlineExpireTime"] = message.OfflineExpireTime
	params["pushType"] = message.Data.GetPushType()
	// 增加pushNetWorkType参数(0:不限;1:wifi;2:4G/3G/2G)
	params["pushNetWorkType"] = message.PushNetWorkType
	params["tpye"] = 2

	if len(message.Conditions) == 0 {
		params["phoneTypeList"] = message.PhoneTypeList
		params["provinceList"] = message.ProvinceList
		params["tagList"] = message.TagList
	} else {
		conditions := message.Conditions
		params["conditions"] = conditions
	}
	params["speed"] = message.Speed
	params["contentType"] = 2
	params["appIdList"] = message.AppIdList

	ret := this.HttpPostJson(params)

	if ret["result"] == "ok" {
		return ret["contentId"]
	} else {
		panic("获取 contentId 失败：")
	}
}

func (this *IGeTui) ContentIdListTaskGroupName(message igetui.IGtListMessage, taskGroupName string) interface{} {
	params := map[string]interface{}{}

	if len(taskGroupName) > 40 {
		panic("TaskGroupName is OverLimit 40")
	}
	params["taskGroupName"] = taskGroupName
	params["action"] = "getContentIdAction"
	params["appkey"] = this.AppKey
	transparent := message.Data.GetTransparent()
	byteArray, _ := proto.Marshal(transparent)
	params["clientData"] = base64.StdEncoding.EncodeToString(byteArray)
	params["transmissionContent"] = message.Data.GetTransmissionContent()
	params["isOffline"] = message.IsOffline
	params["offlineExpireTime"] = message.OfflineExpireTime
	params["pushType"] = message.Data.GetPushType()
	// 增加pushNetWorkType参数(0:不限;1:wifi;2:4G/3G/2G)
	params["pushNetWorkType"] = message.PushNetWorkType
	params["tpye"] = 2
	params["contentType"] = 1

	ret := this.HttpPostJson(params)

	if ret["result"] == "ok" {
		return ret["contentId"]
	} else {
		panic("获取 contentId 失败：")
	}
}

func (this *IGeTui) ContentIdListWith(message igetui.IGtListMessage) interface{} {
	params := map[string]interface{}{}

	params["action"] = "getContentIdAction"
	params["appkey"] = this.AppKey
	transparent := message.Data.GetTransparent()
	byteArray, _ := proto.Marshal(transparent)
	params["clientData"] = base64.StdEncoding.EncodeToString(byteArray)
	params["transmissionContent"] = message.Data.GetTransmissionContent()
	params["isOffline"] = message.IsOffline
	params["offlineExpireTime"] = message.OfflineExpireTime
	params["pushType"] = message.Data.GetPushType()
	// 增加pushNetWorkType参数(0:不限;1:wifi;2:4G/3G/2G)
	params["pushNetWorkType"] = message.PushNetWorkType
	params["tpye"] = 2
	params["contentType"] = 1

	ret := this.HttpPostJson(params)

	if ret["result"] == "ok" {
		return ret["contentId"]
	} else {
		panic("获取 contentId 失败：")
	}
}

func (this *IGeTui) Sign(appKey string, timeStamp int64, masterSecret string) string {
	rawValue := appKey + strconv.FormatInt(timeStamp, 10) + masterSecret
	h := md5.New()
	io.WriteString(h, rawValue)
	return hex.EncodeToString(h.Sum(nil))
}

func (this *IGeTui) CurrentTime() int64 {
	t := time.Now().Unix() * 1000
	return t
}

func (this *IGeTui) HttpPostJson(params map[string]interface{}) map[string]interface{} {
	ret := this.HttpPost(params)
	if ret["result"] == "sign_error" {
		this.connect()
		ret = this.HttpPost(params)
	}
	return ret
}

func (this *IGeTui) HttpPost(params map[string]interface{}) map[string]interface{} {
	data, _ := json.Marshal(params)
	//fmt.Printf("%s\n", data)
	tryTime := 1
tryAgain:
	fmt.Println(this.Host)
	res, err := http.Post(this.Host, "application/json", strings.NewReader(string(data)))
	if err != nil {
		fmt.Println("第"+strconv.Itoa(tryTime)+"次", "请求失败")
		tryTime += 1
		if tryTime < 4 {
			goto tryAgain
		}
		return map[string]interface{}{"result": "post error"}
	}
	body, _ := ioutil.ReadAll(res.Body)
	var ret map[string]interface{}
	json.Unmarshal(body, &ret)
	return ret
}


func NewIGeTui(host string, appKey string, masterSecret string) *IGeTui {
	getui := new(IGeTui)
	getui.serviceMap = "Hello"
	getui.Host = host
	getui.AppKey = appKey
	getui.MasterSecret = masterSecret
	//fmt.Println("init")
	//fmt.Println(getui.serviceMap)
	return getui
}

//host, appKey, masterSecret, ssl = None):
