package try

// 异常需要实现 Throwable 接口
type Throwable interface {

	// 创建新的异常， 创建 名称相同的异常
	//New(msg string) *Throwable
	NewThrow(msg string) Throwable

	// 获取 消息
	GetMsg() string

	// 获取 name
	GetName() string

	// 获取 栈信息
	GetStackTrace() string

	Error() string

	// 打印 栈信息
	PrintStackTrace()
}
