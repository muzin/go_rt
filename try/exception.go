package try

import "fmt"

// 异常 对象
//
// 使用示例：
//		var PointException = try.DeclareException("PointException")
//
//		// 捕获异常  可以将 捕获异常放到 抛异常之前，可以在抛出异常之前 定义 异常的捕获，防止异常没有处理
//		defer try.Catch(try.PointException, func(err Throwable){
//			err.PrintStackTrace()
//		})()
//		// 抛异常
//		try.Throw(try.PointException.NewThrow("指针"))
//
type Exception struct {
	name  string
	msg   string
	stack string
}

// 声明的异常不能修改，
// 抛异常时，需要根据 声明的异常 调用New方法重新创建新的异常
func DeclareException(name string) *Exception {
	return newException(name)
}

func newException(name string) *Exception {
	return &Exception{name: name, msg: "", stack: ""}
}

func (this *Exception) New(msg string) *Exception {
	exception := newException(this.name)
	exception.msg = msg
	return exception
}

func (this Exception) NewThrow(msg string) Throwable {
	exception := newException(this.name)
	exception.msg = msg
	exception.stack = getStackTrace(exception)
	return exception
}

func (this Exception) GetName() string {
	return this.name
}

func (this Exception) GetMsg() string {
	return this.msg
}

func (this Exception) GetStackTrace() string {
	return this.stack
}

func (this Exception) Error() string {
	return fmt.Sprintf("%v: %v", this.name, this.msg)
}

func (this Exception) PrintStackTrace() {
	fmt.Printf("%v", this.stack)
}
