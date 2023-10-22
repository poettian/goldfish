package tapd

import (
	"fmt"
	"goldfish/internal/tapd/api"
	"goldfish/internal/tapd/config"
	"goldfish/internal/tapd/types"
	"goldfish/pkg"
	"log"
	"strings"
)

type Sync struct {
	config   *config.Config
	storyApi *api.StoryApi
	tableApi *api.TableApi
	records  map[string][]string
}

func New() *Sync {
	c := config.GetTapdConfig()
	s := &Sync{
		config:   c,
		storyApi: api.NewStoryApi(&c.Tapd),
		tableApi: api.NewTableApi(&c.Feishu),
		records:  make(map[string][]string),
	}

	return s
}

func (sync *Sync) Run(force bool) {
	// 如果是强制模式，清空表格
	if force {
		success, err := sync.tableApi.TruncateTable()
		if err != nil || !success {
			panic(err)
		}
		log.Println("清空表格成功")
	}
	// 获取需求分页数量
	storiesCount, err := sync.storyApi.GetStoriesCount()
	if err != nil {
		panic(err)
	}
	pageTotal := storiesCount / sync.config.Tapd.PageSize
	if storiesCount%sync.config.Tapd.PageSize != 0 {
		pageTotal++
	}
	log.Printf("需求总数: %d 页数：%d\n", storiesCount, pageTotal)
	// 获取需求并插入表格
	for i := 1; i <= pageTotal; i++ {
		stories, err := sync.storyApi.GetStories(i)
		if err != nil {
			panic(err)
		}
		insertRecordsRequest := sync.getInsertRecords(stories)
		records, err := sync.tableApi.InsertRecords(insertRecordsRequest)
		if err != nil {
			panic(err)
		}
		for k, v := range records {
			sync.records[k] = v
		}
		log.Printf("第 %d 页需求插入表格成功\n", i)
	}
	// 更新表格
	updateRecordRequest := api.UpdateRecordsRequest{
		Records: make([]api.UpdateRecord, 0, len(sync.records)),
	}
	for _, v := range sync.records {
		if v[1] == "0" {
			continue
		}
		if _, ok := sync.records[v[1]]; !ok {
			continue
		}
		parentRecordIds := []string{sync.records[v[1]][0]}
		updateRecordRequest.Records = append(updateRecordRequest.Records, api.UpdateRecord{
			RecordId: v[0],
			Fields: api.UpdateRecordFields{
				ParentRecordId: parentRecordIds,
			},
		})
	}
	success, err := sync.tableApi.UpdateRecords(updateRecordRequest)
	if err != nil || !success {
		panic(err)
	}
	log.Println("表格更新数据成功")
}

func (sync *Sync) getInsertRecords(stories []types.Story) api.InsertRecordsRequest {
	insertRecordsRequest := api.InsertRecordsRequest{
		Records: make([]api.InsertRecord, 0, len(stories)),
	}
	for _, story := range stories {
		fields := types.Fields{
			ID: story.ID,
			Name: &types.Link{
				Link: fmt.Sprintf("https://www.tapd.cn/%s/prong/stories/view/%s", sync.config.Tapd.WorkspaceId, story.ID),
				Text: story.Name,
			},
			Created:         pkg.TimeStrToUnixMilli(story.Created),
			Modified:        pkg.TimeStrToUnixMilli(story.Modified),
			Status:          types.StoryStatus[story.Status],
			EffortCompleted: pkg.StrToFloat64(story.EffortCompleted),
			Progress:        pkg.StrToFloat64(story.Progress) / 100,
			CustomFieldFour: story.CustomFieldFour,
		}
		if story.Begin != "" {
			fields.Begin = pkg.DateStrToUnixMilli(story.Begin)
		}
		if story.Due != "" {
			fields.Due = pkg.DateStrToUnixMilli(story.Due)
		}
		if story.Developer != "" {
			fields.Developer = strings.Split(strings.TrimRight(story.Developer, ";"), ";")
		}
		if story.IterationID != "0" {
			fields.Iteration = &types.Link{
				Link: fmt.Sprintf("https://www.tapd.cn/%s/prong/iterations/view/%s", sync.config.Tapd.WorkspaceId, story.IterationID),
				Text: story.IterationID,
			}
		}
		if story.ReleaseID != "0" {
			fields.Release = &types.Link{
				Link: fmt.Sprintf("https://www.tapd.cn/%s/releases/view/%s", sync.config.Tapd.WorkspaceId, story.ReleaseID),
				Text: story.ReleaseID,
			}
		}
		if story.Effort != "" {
			fields.Effort = pkg.StrToFloat64(story.Effort)
		}
		if story.ParentID != "0" {
			fields.ParentId = story.ParentID
		}

		insertRecordsRequest.Records = append(insertRecordsRequest.Records, api.InsertRecord{Fields: fields})
	}

	return insertRecordsRequest
}
