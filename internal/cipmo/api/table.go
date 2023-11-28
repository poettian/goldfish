package api

import (
	"encoding/json"
	"fmt"
	"goldfish/internal/cipmo/config"
	"goldfish/internal/cipmo/types"
	"goldfish/pkg"
	"log"
	"strings"
)

type TableApi struct {
	appId             string
	appSecret         string
	tenantAccessToken string
	docsToken         string
	storyTableId      string
	pageSize          int
	hasMore           bool
	pageToken         string
}

type InsertRecord struct {
	Fields types.Fields `json:"fields"`
}

type InsertRecordsRequest struct {
	Records []InsertRecord `json:"records"`
}

type UpdateRecordFields struct {
	ParentRecordId []string `json:"父记录"`
}

type UpdateRecord struct {
	RecordId string             `json:"record_id"`
	Fields   UpdateRecordFields `json:"fields"`
}

type UpdateRecordsRequest struct {
	Records []UpdateRecord `json:"records"`
}

type TokenResponse struct {
	Code              int    `json:"code"`
	Msg               string `json:"msg"`
	TenantAccessToken string `json:"tenant_access_token"`
	Expire            int    `json:"expire"`
}

type RecordResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		HasMore   bool           `json:"has_more"`
		PageToken string         `json:"page_token"`
		Total     int            `json:"total"`
		Items     []types.Record `json:"items"`
	} `json:"data"`
}

type DelRecordResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		Records []struct {
			Deleted  bool   `json:"deleted"`
			RecordID string `json:"record_id"`
		} `json:"records"`
	} `json:"data"`
}

type InsertRecordResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		Records []types.Record `json:"records"`
	} `json:"data"`
}

type UpdateRecordResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func (b *TableApi) getRecordIds() (recordIds []string, err error) {
	if !b.hasMore {
		return nil, nil
	}
	// 发起HTTP请求
	request := pkg.NewRequest()
	request.Method = "GET"
	request.Url = fmt.Sprintf("https://open.feishu.cn/open-apis/bitable/v1/apps/%s/tables/%s/records", b.docsToken, b.storyTableId)
	request.Query = map[string][]string{
		"field_names": {`["需求ID"]`},
		"page_size":   {fmt.Sprintf("%d", b.pageSize)},
	}
	if b.pageToken != "" {
		request.Query["page_token"] = []string{b.pageToken}
	}
	request.Headers = map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", b.tenantAccessToken),
	}
	responseBody, err := request.Do()
	if err != nil {
		return nil, err
	}
	// 解析响应内容为json
	var recordResponse RecordResponse
	err = json.Unmarshal(responseBody, &recordResponse)
	if err != nil {
		return nil, fmt.Errorf("JSON解析错误: %s %w", request.Url, err)
	}
	// 判断状态码(1254024: InvalidFieldNames)
	if recordResponse.Code != 0 {
		return nil, fmt.Errorf("API请求出错: %s %s", request.Url, recordResponse.Msg)
	}
	// 获取记录ID
	recordIds = make([]string, 0, len(recordResponse.Data.Items))
	for _, v := range recordResponse.Data.Items {
		recordIds = append(recordIds, v.RecordId)
	}
	b.hasMore = recordResponse.Data.HasMore
	b.pageToken = recordResponse.Data.PageToken

	return recordIds, nil
}

func (b *TableApi) deleteRecords(recordIds []string) (bool, error) {
	// 发起HTTP请求
	request := pkg.NewRequest()
	request.Method = "POST"
	request.Url = fmt.Sprintf("https://open.feishu.cn/open-apis/bitable/v1/apps/%s/tables/%s/records/batch_delete", b.docsToken, b.storyTableId)
	request.Body = fmt.Sprintf(`{"records": ["%s"]}`, strings.Join(recordIds, `","`))
	request.Headers = map[string]string{
		"Content-Type":  "application/json; charset=utf-8",
		"Authorization": fmt.Sprintf("Bearer %s", b.tenantAccessToken),
	}
	responseBody, err := request.Do()
	if err != nil {
		return false, err
	}
	// 解析响应内容为json
	var delRecordResponse DelRecordResponse
	err = json.Unmarshal(responseBody, &delRecordResponse)
	if err != nil {
		return false, fmt.Errorf("JSON解析错误: %s %w", request.Url, err)
	}
	// 判断状态码
	if delRecordResponse.Code != 0 {
		return false, fmt.Errorf("API请求出错: %s %s", request.Url, delRecordResponse.Msg)
	}
	// 判断是否删除成功
	success := true
	for _, v := range delRecordResponse.Data.Records {
		if !v.Deleted {
			log.Printf("删除记录失败: %s\n", v.RecordID)
			success = false
		}
	}

	return success, nil
}

func (b *TableApi) TruncateTable() (bool, error) {
	b.pageToken = ""
	b.hasMore = true
	var recordIds []string
	var err error
	var success bool
	for b.hasMore {
		recordIds, err = b.getRecordIds()
		if err != nil {
			return false, err
		}
		success, err = b.deleteRecords(recordIds)
		if err != nil || !success {
			return false, err
		}
	}

	return true, nil
}

func (b *TableApi) InsertRecords(recordRequest InsertRecordsRequest) (map[string][]string, error) {
	// 发起HTTP请求
	request := pkg.NewRequest()
	request.Method = "POST"
	request.Url = fmt.Sprintf("https://open.feishu.cn/open-apis/bitable/v1/apps/%s/tables/%s/records/batch_create", b.docsToken, b.storyTableId)
	request.Headers = map[string]string{
		"Content-Type":  "application/json; charset=utf-8",
		"Authorization": fmt.Sprintf("Bearer %s", b.tenantAccessToken),
	}
	body, err := json.Marshal(recordRequest)
	if err != nil {
		return nil, fmt.Errorf("JSON编码错误: %s %w", request.Url, err)
	}
	request.Body = string(body)
	responseBody, err := request.Do()
	if err != nil {
		return nil, err
	}
	// 解析响应内容为json
	var insertRecordResponse InsertRecordResponse
	err = json.Unmarshal(responseBody, &insertRecordResponse)
	if err != nil {
		return nil, fmt.Errorf("JSON解析错误: %s %w", request.Url, err)
	}
	// 判断状态码
	if insertRecordResponse.Code != 0 {
		return nil, fmt.Errorf("API请求出错: %s %s", request.Url, insertRecordResponse.Msg)
	}
	// 处理返回结果
	var records = make(map[string][]string)
	for _, v := range insertRecordResponse.Data.Records {
		records[v.Fields.ID] = []string{v.RecordId, v.Fields.ParentId}
	}

	return records, nil
}

func (b *TableApi) UpdateRecords(updateRecordsRequest UpdateRecordsRequest) (bool, error) {
	// 发起HTTP请求
	request := pkg.NewRequest()
	request.Method = "POST"
	request.Url = fmt.Sprintf("https://open.feishu.cn/open-apis/bitable/v1/apps/%s/tables/%s/records/batch_update", b.docsToken, b.storyTableId)
	request.Headers = map[string]string{
		"Content-Type":  "application/json; charset=utf-8",
		"Authorization": fmt.Sprintf("Bearer %s", b.tenantAccessToken),
	}
	body, err := json.Marshal(updateRecordsRequest)
	if err != nil {
		return false, fmt.Errorf("JSON编码错误: %s %w", request.Url, err)
	}
	request.Body = string(body)
	responseBody, err := request.Do()
	if err != nil {
		return false, err
	}
	// 解析响应内容为json
	var updateRecordResponse UpdateRecordResponse
	err = json.Unmarshal(responseBody, &updateRecordResponse)
	if err != nil {
		return false, fmt.Errorf("JSON解析错误: %s %w", request.Url, err)
	}
	// 判断状态码
	if updateRecordResponse.Code != 0 {
		return false, fmt.Errorf("API请求出错: %s %s", request.Url, updateRecordResponse.Msg)
	}

	return true, nil
}

func (b *TableApi) getAccessToken() (string, error) {
	// 发起HTTP请求
	request := pkg.NewRequest()
	request.Method = "POST"
	request.Url = "https://open.feishu.cn/open-apis/auth/v3/tenant_access_token/internal"
	request.Body = fmt.Sprintf(`{"app_id": "%s","app_secret": "%s"}`, b.appId, b.appSecret)
	request.Headers = map[string]string{
		"Content-Type":  "application/json; charset=utf-8",
		"Authorization": fmt.Sprintf("Bearer %s", b.tenantAccessToken),
	}
	responseBody, err := request.Do()
	if err != nil {
		return "", err
	}
	// 解析响应内容为json
	var tokenResponse TokenResponse
	err = json.Unmarshal(responseBody, &tokenResponse)
	if err != nil {
		return "", fmt.Errorf("JSON解析错误: %s %w", request.Url, err)
	}
	// 判断状态码
	if tokenResponse.Code != 0 {
		return "", fmt.Errorf("API请求出错: %s %s", request.Url, tokenResponse.Msg)
	}

	return tokenResponse.TenantAccessToken, nil
}

func NewTableApi(c *config.Feishu) *TableApi {
	api := &TableApi{
		appId:             c.AppId,
		appSecret:         c.AppSecret,
		tenantAccessToken: c.TenantAccessToken,
		docsToken:         c.DocsToken,
		storyTableId:      c.StoryTableId,
		pageSize:          c.PageSize,
	}
	token, err := api.getAccessToken()
	if err != nil {
		log.Fatalln(err)
	}
	api.tenantAccessToken = token

	return api
}
