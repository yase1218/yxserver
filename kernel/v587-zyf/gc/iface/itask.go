package iface

type ITask interface {
	Do()
}

type ITaskEx interface {
	ITask
	// 返回分片键（如用户ID），空字符串表示普通任务
	ShardKey() string
	// 是否要求严格串行（当ShardKey非空时有效）
	RequireSerialized() bool
}
