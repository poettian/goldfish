package types

type Link struct {
	Link string
	Text string
}

type Fields struct {
	ID              string   `json:"需求ID"`
	Name            Link     `json:"需求"`
	Created         uint64   `json:"创建时间"`
	Modified        uint64   `json:"修改时间"`
	Status          string   `json:"状态"`
	Begin           uint64   `json:"预计开始"`
	Due             uint64   `json:"预计结束"`
	Developer       []string `json:"开发人员"`
	Iteration       Link     `json:"迭代"`
	Release         Link     `json:"发布计划"`
	Effort          float64  `json:"预估工时"`
	EffortCompleted float64  `json:"完成工时"`
	Progress        float64  `json:"进度"`
	CustomFieldFour string   `json:"优先级"`
	ParentId        string   `json:"父需求ID"`
	ParentRecordId  []string `json:"父记录"`
	//CustomFieldFive string `json:"开发端"`
	//CustomField9    string `json:"业务线"`
}

type Record struct {
	RecordId string `json:"record_id"`
	Fields   Fields `json:"fields"`
}
