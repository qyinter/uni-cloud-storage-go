package unicloud

// StringObject 请求 AccessToken 的参数
type StringObject struct {
	Method    string `json:"method,omitempty"`
	Params    string `json:"params,omitempty"`
	SpaceId   string `json:"spaceId,omitempty"`
	Timestamp int64  `json:"timestamp,omitempty"`
}

// AccessToken 获取 AccessToken 的结构体
type AccessToken struct {
	Success bool `json:"success"`
	Data    struct {
		AccessToken     string `json:"accessToken"`
		ExpiresInSecond int    `json:"expiresInSecond"`
	} `json:"data"`
}

// CreatFileNameParams 请求 CreatFileName 的参数
type CreatFileNameParams struct {
	Method    string `json:"method,omitempty"`
	Params    string `json:"params,omitempty"`
	SpaceId   string `json:"spaceId,omitempty"`
	Timestamp int64  `json:"timestamp,omitempty"`
	Token     string `json:"token,omitempty"`
}

// FileInfoData 请求 CreatFileName 的返回值
type FileInfoData struct {
	Success bool `json:"success"`
	Data    struct {
		Id          string `json:"id"`
		CdnDomain   string `json:"cdnDomain"`
		Signature   string `json:"signature"`
		Policy      string `json:"policy"`
		AccessKeyId string `json:"accessKeyId"`
		OssPath     string `json:"ossPath"`
		Host        string `json:"host"`
	} `json:"data"`
}
