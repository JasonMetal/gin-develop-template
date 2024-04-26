package testDbEntity

type TestData struct {
	TestId    int32  `gorm:"column:test_id"                 json:"test_id"`
	TestName  string `gorm:"column:test_name"              json:"test_name"`
	TestUnion string `gorm:"column:test_union"            json:"test_union"`
	//TestTag    int32  `gorm:"column:test_tag"             json:"test_tag"`
	TestStatus int32 `gorm:"column:test_status"              json:"test_status"`
}

type OnlineList struct {
	List []*TestData `json:"list"`
}

func (TestData) TableName() string {
	return "test"
}

type Condition struct {
	Id     uint32
	Page   int
	Limit  int
	Offset int
}

type ConditionTest struct {
	TestId     int32  `gorm:"column:test_id"                 json:"test_id"`
	TestName   string `gorm:"column:test_name"              json:"test_name"`
	TestUnion  string `gorm:"column:test_union"            json:"test_union"`
	TestTag    int32  `gorm:"column:test_tag"             json:"test_tag"`
	TestStatus int32  `gorm:"column:test_status"              json:"test_status"`
}
