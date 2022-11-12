package unicloud

import (
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"path"
	"time"
)

const baseHost = "https://api.bspapp.com"
const clientSecret = "你的clientSecret"
const spaceId = "你的spaceId"

func Upload(localPath string) string {
	base := path.Base(localPath)

	token := GetAccessToken()
	fileInfoData := CreatFileName(base, token)
	code := uploadFile(localPath, fileInfoData)
	if code == 200 {
		flag := CheckFile(token, fileInfoData.Data.Id)
		if flag == true {
			uniCloudUrl := fmt.Sprintf("https://%v/%v", fileInfoData.Data.CdnDomain, fileInfoData.Data.OssPath)
			return uniCloudUrl
		} else {
			log.Error("文件上传异常")
		}
	}
	return ""
}

func GetAccessToken() string {
	data := StringObject{
		Method:    "serverless.auth.user.anonymousAuthorize",
		Params:    "{}",
		SpaceId:   spaceId,
		Timestamp: time.Now().Unix(),
	}
	client := resty.New()
	resp, _ := client.R().
		SetHeaders(map[string]string{
			"Content-Type":      "application/json",
			"x-serverless-sign": sign(data, clientSecret),
		}).
		SetBody(data).
		Post(baseHost + "/client")
	var accessToken AccessToken
	json.Unmarshal(resp.Body(), &accessToken)
	return accessToken.Data.AccessToken
}

func CreatFileName(filename, accessToken string) FileInfoData {
	data := CreatFileNameParams{
		Method:    "serverless.file.resource.generateProximalSign",
		Params:    "{\"env\":\"public\",\"filename\":\"" + filename + "\"}",
		SpaceId:   spaceId,
		Timestamp: time.Now().Unix(),
		Token:     accessToken,
	}
	client := resty.New()
	resp, _ := client.R().
		SetHeaders(map[string]string{
			"Content-Type":      "application/json",
			"x-basement-token":  accessToken,
			"x-serverless-sign": sign2(data, clientSecret),
		}).
		SetBody(data).
		Post(baseHost + "/client")
	var fileInfoData FileInfoData
	json.Unmarshal(resp.Body(), &fileInfoData)
	return fileInfoData
}

func uploadFile(path string, vo FileInfoData) int {
	params := vo.Data
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.WriteField("Cache-Control", "max-age=2592000")
	writer.WriteField("Content-Disposition", "attachment")
	writer.WriteField("OSSAccessKeyId", params.AccessKeyId)
	writer.WriteField("Signature", params.Signature)
	writer.WriteField("host", params.Host)
	writer.WriteField("id", params.Id)
	writer.WriteField("key", params.OssPath)
	writer.WriteField("policy", params.Policy)
	writer.WriteField("success_action_status", "200")
	field, _ := writer.CreateFormField("file")
	file, _ := ioutil.ReadFile(path)
	field.Write(file)
	err := writer.Close()
	if err != nil {
		return 0
	}
	req, err := http.NewRequest("POST", "https://"+params.Host, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("X-OSS-server-side-encrpytion", "AES256")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0
	}
	return resp.StatusCode

}
func CheckFile(accessToken, id string) interface{} {
	data := CreatFileNameParams{
		Method:    "serverless.file.resource.report",
		Params:    "{\"id\":\"" + id + "\"}",
		SpaceId:   spaceId,
		Timestamp: time.Now().Unix(),
		Token:     accessToken,
	}
	client := resty.New()
	resp, _ := client.R().
		SetHeaders(map[string]string{
			"Content-Type":      "application/json",
			"x-basement-token":  accessToken,
			"x-serverless-sign": sign2(data, clientSecret),
		}).
		SetBody(data).
		Post(baseHost + "/client")
	var dataMap map[string]interface{}
	err := json.Unmarshal(resp.Body(), &dataMap)
	if err != nil {
		log.Info("序列化uniCloud文件连接失败")
	}
	return dataMap["success"]

}
func sign(data StringObject, clientSecret string) string {
	param := fmt.Sprintf("method=%v&params=%v&spaceId=%v&timestamp=%v",
		data.Method, data.Params, data.SpaceId, data.Timestamp)
	hmacMd5 := HmacMd5(clientSecret, param)
	return hmacMd5
}
func sign2(data CreatFileNameParams, clientSecret string) string {
	param := fmt.Sprintf("method=%v&params=%v&spaceId=%v&timestamp=%v&token=%v",
		data.Method, data.Params, data.SpaceId, data.Timestamp, data.Token)
	hmacMd5 := HmacMd5(clientSecret, param)
	return hmacMd5
}
func HmacMd5(key, data string) string {
	h := hmac.New(md5.New, []byte(key))
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum([]byte("")))
}
