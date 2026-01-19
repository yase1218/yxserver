package iface

type IState interface {
	Enter()         // 触发状态
	ExecuteBefore() // 状态执行前
	Execute()       // 状态执行
	ExecuteAfter()  // 状态执行后
	End()           // 状态结束
}

type StateTrigger int
