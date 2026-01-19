package db_model

type BaseTable struct {
}

func (bt *BaseTable) TableName() string {
	return "defName"
}

func (bt *BaseTable) TableEngine() string {
	return "Innodb"
}

func (bt *BaseTable) TableEncode() string {
	return "utf8"
}

func (bt *BaseTable) TableComment() string {
	return ""
}

func (bt *BaseTable) TableIndex() [][]string {
	return [][]string{}
}

func (bt *BaseTable) TableUnique() [][]string {
	return [][]string{}
}
