package igetui

type Target struct {
	AppId    string
	ClientId string
	Alias string
}

func NewTarget(appid string, clientid string ,alias string) *Target {
	return &Target{
		AppId:    appid,
		ClientId: clientid,
		Alias :alias,
	}
}
