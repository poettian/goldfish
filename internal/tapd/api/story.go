package api

import (
	"encoding/json"
	"fmt"
	"goldfish/internal/tapd/config"
	"goldfish/internal/tapd/types"
	"goldfish/pkg"
	"strings"
)

type StoryApi struct {
	appId     string
	appSecret string
	fields    []string
	filter    map[string]string
	order     string
	pageSize  int
}

type CountResponse struct {
	Status int `json:"status"`
	Data   struct {
		Count int `json:"count"`
	} `json:"data"`
	Info string `json:"info"`
}

type StoryResponse struct {
	Status int `json:"status"`
	Data   []struct {
		Story types.Story `json:"Story"`
	} `json:"data"`
	Info string `json:"info"`
}

func (t *StoryApi) GetStoriesCount() (int, error) {
	// 发起HTTP请求
	request := pkg.NewRequest()
	request.Method = "GET"
	request.Url = "https://api.tapd.cn/stories/count"
	request.BasicAuth = []string{t.appId, t.appSecret}
	request.Query = make(map[string][]string)
	for k, v := range t.filter {
		request.Query[k] = []string{v}
	}
	responseBody, err := request.Do()
	if err != nil {
		return 0, err
	}
	// 解析响应内容为json
	var countResponse CountResponse
	err = json.Unmarshal(responseBody, &countResponse)
	if err != nil {
		return 0, fmt.Errorf("JSON解析错误: %s %w", request.Url, err)
	}
	// 判断状态码
	if countResponse.Status != 1 {
		return 0, fmt.Errorf("API请求出错: %s %s", request.Url, countResponse.Info)
	}

	return countResponse.Data.Count, nil
}

// GetStories 获取故事列表
func (t *StoryApi) GetStories(page int) ([]types.Story, error) {
	// 发起HTTP请求
	request := pkg.NewRequest()
	request.Method = "GET"
	request.Url = "https://api.tapd.cn/stories"
	request.BasicAuth = []string{t.appId, t.appSecret}
	request.Query = map[string][]string{
		"fields": {strings.Join(t.fields, ",")},
		"order":  {t.order},
		"page":   {fmt.Sprintf("%d", page)},
		"limit":  {fmt.Sprintf("%d", t.pageSize)},
	}
	for k, v := range t.filter {
		request.Query[k] = []string{v}
	}
	responseBody, err := request.Do()
	if err != nil {
		return nil, err
	}
	// 解析响应内容为json
	var storyResponse StoryResponse
	err = json.Unmarshal(responseBody, &storyResponse)
	if err != nil {
		return nil, fmt.Errorf("JSON解析错误: %s %w", request.Url, err)
	}
	// 判断状态码
	if storyResponse.Status != 1 {
		return nil, fmt.Errorf("API请求出错: %s %s", request.Url, storyResponse.Info)
	}
	// TODO: 处理迭代和发布计划字段
	storyList := make([]types.Story, 0, len(storyResponse.Data))
	for _, v := range storyResponse.Data {
		storyList = append(storyList, v.Story)
	}
	return storyList, nil
}

func NewStoryApi(c *config.Tapd) *StoryApi {
	// TODO：从 types.Story 的 json tag 中解析
	fields := []string{
		"id",
		"workspace_id",
		"name",
		"developer",
		"status",
		"begin",
		"due",
		"effort",
		"effort_completed",
		"progress",
		"custom_field_four", // 优先级别
		"custom_field_five", // 开发端
		"custom_field_9",    // 业务线
		"parent_id",
		"iteration_id",
		"release_id",
		"created",
		"modified",
	}
	return &StoryApi{
		appId:     c.AppId,
		appSecret: c.AppSecret,
		fields:    fields,
		filter: map[string]string{
			"workspace_id":      c.WorkspaceId,
			"custom_field_five": "后端", // TODO: 从配置文件中读取
			"status":            "<>resolved",
		},
		order:    "id desc",
		pageSize: c.PageSize,
	}
}
