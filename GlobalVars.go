/*
Created on Fri Sep 16 17:04:36 2024
@author:liwancai

	QQ:248411282
	Tel:13199701121
*/
package EQUseApi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/li-wancai/GoScripts/DirsFile"

	"github.com/li-wancai/GoScripts/Formulae"
)

// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
var (
	EQDataUseHTTPS  bool
	EQDataApiToken  string
	EQDataHTTPAPI   []string
	SendToGroupList []string
	EQUseApi_CFG    map[string]interface{}
)

func SetUseapi(config map[string]interface{}) {
	EQUseApi_CFG, _ = DirsFile.ReadToml(config["EQUseApi_FileName"].(string), config["EQUseApi_TomlPath"].(string)) //读取配置文件
	EQDataApiToken = EQUseApi_CFG["EQDataApiToken"].(string)                                                        //数据接口的密匙
	EQDataHTTPAPI = Formulae.ToStringList(EQUseApi_CFG["EQDataHTTPAPI"].([]interface{}))                            //数据接口的IP地址
	SendToGroupList = Formulae.ToStringList(EQUseApi_CFG["SendToGroupList"].([]interface{}))                        //需要推送消息到群的列表
	EQDataUseHTTPS = EQUseApi_CFG["EQDataUseHTTPS"].(bool)                                                          //是否采用https传输
}

type EQApiN struct {
	Cookies []*http.Cookie
	Session *http.Client
	BaseUrl []string
	Token   string
	HTTPS   bool
	Headers map[string]string
}

func EQApi() *EQApiN {
	return &EQApiN{
		HTTPS:   EQDataUseHTTPS,
		Cookies: []*http.Cookie{},
		Session: &http.Client{},
		BaseUrl: EQDataHTTPAPI,
		Token:   EQDataApiToken,
		Headers: map[string]string{
			"content-type":  "application/json;charset=utf-8",
			"Authorization": EQDataApiToken,
			"Connection":    "keep-alive",
		},
	}
}

func (api *EQApiN) SetApiToken(token string) {
	api.Token = token
	api.Headers["Authorization"] = "Bearer " + token
}

func (api *EQApiN) Login(ApiEndUrl string, username string, password string) {
	loginbody := map[string]interface{}{
		"TELorEmail": username,
		"Password":   password,
	}
	_, err := api.Request(loginbody, ApiEndUrl)
	if err != nil {
		log.Errorf("【失败】%s登录错误:%s", username, err)
		return
	}
}

// Request 方法用于向指定的 API 发送 POST 请求，可以处理有或无请求体的情况。
func (api *EQApiN) Request(body interface{}, ApiEndUrl string) (*http.Response, error) {
	jsonData, err := json.Marshal(body) // 将请求体 body 转换为 JSON 格式的字节流，以便在网络上传输。
	if err != nil {                     // 如果转换过程中出现错误，返回错误信息。
		return nil, fmt.Errorf("json转换过程中出现错误: %v", err)
	}
	for _, baseUrl := range api.BaseUrl { // 遍历 api.BaseUrl 列表中的所有基础 URL，尝试构建和发送请求。
		url := fmt.Sprintf("%s/%s", baseUrl, ApiEndUrl) // 构建完整的 URL，包括基础 URL 和 API 的结尾 URL。
		url = strings.Replace(url, "//", "/", -1)       // 确保 URL 中不会有连续的两个斜杠，这在 URL 构造中是不必要的。
		urlHeader := "http://"                          // 根据 api.HTTPS 字段判断是否使用 HTTPS 协议进行安全连接。
		if api.HTTPS {
			urlHeader = "https://"
		}
		// 创建一个新的 HTTP POST 请求，包含构建好的 URL 和 JSON 格式的请求体。
		req, err := http.NewRequest("POST", urlHeader+url, bytes.NewBuffer(jsonData))
		if err != nil { // 如果创建请求时发生错误，返回错误信息。
			return nil, fmt.Errorf("创建请求时发生错误: %v", err)
		}
		// 遍历 Headers 字典，将每一个键值对设置到请求的头部，确保请求携带正确的元数据。
		for key, value := range api.Headers {
			req.Header.Set(key, value)
		}
		// 遍历 Cookies 列表，将每一个 Cookie 添加到请求中，以维持会话状态。
		for _, cookie := range api.Cookies {
			req.AddCookie(cookie)
		}
		// 使用 api.Session 发送 HTTP 请求，Session 可能包含了重试逻辑或连接池管理。
		resp, err := api.Session.Do(req)
		if err != nil {
			continue // 如果请求发送失败，跳过本次循环，尝试下一个基础 URL。
		}
		if resp.StatusCode == 200 { // 检查响应状态码是否为 200，即 HTTP 成功状态。
			return resp, nil // 如果状态码为 200，表示请求成功，返回响应对象。
		}
	}
	return nil, fmt.Errorf("请求失败") // 如果所有的基础 URL 都尝试过但没有成功，返回一个通用的请求失败错误。
}
func (api *EQApiN) PostData(body map[string]interface{}, ApiEndUrl string) (map[string]interface{}, error) {
	resp, err := api.Request(body, ApiEndUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
