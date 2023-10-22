package types

type Link struct {
	Link string `json:"link"`
	Text string `json:"text"`
}

type Fields struct {
	ID              string   `json:"需求ID,omitempty"`
	Name            *Link    `json:"需求,omitempty"`
	Created         uint64   `json:"创建时间,omitempty"`
	Modified        uint64   `json:"修改时间,omitempty"`
	Status          string   `json:"状态,omitempty"`
	Begin           uint64   `json:"预计开始,omitempty"`
	Due             uint64   `json:"预计结束,omitempty"`
	Developer       []string `json:"开发人员,omitempty"`
	Iteration       *Link    `json:"迭代,omitempty"`
	Release         *Link    `json:"发布计划,omitempty"`
	Effort          float64  `json:"预估工时,omitempty"`
	EffortCompleted float64  `json:"完成工时,omitempty"`
	Progress        float64  `json:"进度,omitempty"`
	CustomFieldFour string   `json:"优先级,omitempty"`
	ParentId        string   `json:"父需求ID,omitempty"`
	ParentRecordId  []string `json:"父记录,omitempty"`
	//CustomFieldFive string `json:"开发端,omitempty"`
	//CustomField9    string `json:"业务线,omitempty"`
}

type Record struct {
	RecordId string `json:"record_id"`
	Fields   Fields `json:"fields"`
}
