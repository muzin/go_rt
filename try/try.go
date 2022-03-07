package try

import (
	"bytes"
	"fmt"
	"reflect"
	"runtime"
)

// 抛异常
//
// @param err *Exception 异常 对象
func Throw(err Throwable) {
	panic(err)
}

// 捕获异常
// @param err Throwable 异常 对象
// @param cb Throwable 异常对象 的 处理函数
func Catch(err Throwable, cb func(Throwable)) func() {
	return func() {
		if e := recover(); e != nil {
			typeOfE := reflect.TypeOf(e)
			typeOfErr := reflect.TypeOf(err)
			if typeOfE == typeOfErr {
				errname := (err).GetName()
				ename := (e.(Throwable)).GetName()
				if errname == ename {
					cb(e.(Throwable))
				} else {
					// 如果 没有 匹配到 异常 继续向上 抛 异常
					panic(e)
				}
			} else {
				// 如果 没有 匹配到 异常 继续向上 抛 异常
				panic(e)
			}
		}
	}
}

func CatchUncaughtException(cb func(Throwable)) func() {
	return func() {
		if e := recover(); e != nil {
			switch e.(type) {
			case Throwable:
				cb(e.(Throwable))
			default:
				// 如果 没有 匹配到 异常 继续向上 抛 异常
				err := e.(error)
				cb(UnhandledError.NewThrow(err.Error()))
			}
		}
	}
}

// 尝试捕获作用域中的异常，指定处理方法
//
// catches 使用 `Catch(err Throwable, cb func(Throwable))`
//
func Try(scope func(), catches ...func()) {
	for i := (len(catches) - 1); i >= 0; i-- {
		defer (catches[i])()
	}
	scope()
}

// 打印堆栈信息
func getStackTrace(err interface{}) string {
	buf := new(bytes.Buffer)
	fmt.Fprintf(buf, "%v\n", err)
	for i := 1; ; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		funcForPc := runtime.FuncForPC(pc)
		funcName := funcForPc.Name()
		fmt.Fprintf(buf, "at %v ( %s:%d ) (0x%x)\n", funcName, file, line, pc)
	}

	return buf.String()
}

// declare Error Exception
var Error = DeclareException("Error")

// declare UnhandledError Exception
var UnhandledError = DeclareException("UnhandledError")
