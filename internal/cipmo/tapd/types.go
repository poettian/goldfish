package tapd

import "reflect"

/* 需求状态
"options": {
	"status_6": "需求已内审或跳过",
	"status_7": "新需求",
	"status_10": "需求已评估、内审或审核",
	"status_9": "已确认发布计划（打牌会）",
	"planning": "产品规划设计中",
	"planned": "已完成产品设计",
	"status_5": "UE设计",
	"UI_design": "UI设计",
	"auditing": "待宣讲",
	"audited": "已宣讲",
	"developing": "实现中",
	"status_3": "开发完成",
	"status_4": "产品体验阶段",
	"for_test": "转测试",
	"testing": "测试中",
	"status_2": "测试通过",
	"resolved": "已上线",
	"rejected": "已拒绝",
	"status_11": "已确认迭代计划",
	"status_12": "技术预研中",
	"status_13": "技术规划中",
	"status_16": "待验收",
	"status_17": "已验收",
	"status_18": "已技术方案评审",
	"status_19": "已完成交互优化方案",
	"status_20": "验收中",
	"status_21": "跳过技术方案评审"
},
*/

var StoryStatus = map[string]string{
	"status_6":   "需求已内审或跳过",
	"status_7":   "新需求",
	"status_10":  "需求已评估、内审或审核",
	"status_9":   "已确认发布计划（打牌会）",
	"planning":   "产品规划设计中",
	"planned":    "已完成产品设计",
	"status_5":   "UE设计",
	"UI_design":  "UI设计",
	"auditing":   "待宣讲",
	"audited":    "已宣讲",
	"developing": "实现中",
	"status_3":   "开发完成",
	"status_4":   "产品体验阶段",
	"for_test":   "转测试",
	"testing":    "测试中",
	"status_2":   "测试通过",
	"resolved":   "已上线",
	"rejected":   "已拒绝",
	"status_11":  "已确认迭代计划",
	"status_12":  "技术预研中",
	"status_13":  "技术规划中",
	"status_16":  "待验收",
	"status_17":  "已验收",
	"status_18":  "已技术方案评审",
	"status_19":  "已完成交互优化方案",
	"status_20":  "验收中",
	"status_21":  "跳过技术方案评审",
}

type Story struct {
	ID              string `json:"id"`                // 需求id
	Name            string `json:"name"`              // 需求名称
	WorkspaceID     string `json:"workspace_id"`      // 项目id
	Created         string `json:"created"`           // 创建时间: 2023-09-13 10:41:54
	Modified        string `json:"modified"`          // 修改时间: 2023-09-13 10:41:54
	Status          string `json:"status"`            // 状态: 见注释
	Begin           string `json:"begin"`             // 预计开始：2023-10-19
	Due             string `json:"due"`               // 预计结束：2023-11-02
	Developer       string `json:"developer"`         // 开发者：张三;李四;
	IterationID     string `json:"iteration_id"`      // 迭代id：1145976096001001368
	ReleaseID       string `json:"release_id"`        // 发布计划：1145976096001000245
	ParentID        string `json:"parent_id"`         // 父需求id
	Effort          string `json:"effort"`            // 预估工时：11.5
	EffortCompleted string `json:"effort_completed"`  // 完成工时：3
	Progress        string `json:"progress"`          // 进度：26
	CustomFieldFour string `json:"custom_field_four"` // 优先级别
	CustomFieldFive string `json:"custom_field_five"` // 开发端
	CustomField9    string `json:"custom_field_9"`    // 业务线
}

func (s *Story) GetFields() []string {
	t := reflect.TypeOf(s)
}
