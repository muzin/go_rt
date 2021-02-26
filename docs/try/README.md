## try

`try`包对 panic/recover 进行封装，更接近与面向对象的try/catch

使用步骤：
1. 声明异常
2. 在`try.Throw()`前调用`try.Catch()`
3. 抛异常

注意事项：
1. 调用`try.Catch()`时，在前面必须声明`defer`
2. `try.Catch()`返回的是函数，前面声明了`defer`, 该函数必须被调用
   
   示例：`defer try.Catch()()`
   
3. `try.Catch()`函数调用必须写到`try.Throw()`函数后面
   
   其中原理请了解`panic/recover`工作机制

> 声明异常推荐全局声明


#### 示例：

```go

// 声明 异常
var PointException = try.DeclareException("PointException")

// 捕获异常  可以将 捕获异常放到 抛异常之前，可以在抛出异常之前 定义 异常的捕获，防止异常没有处理
defer try.Catch(try.PointException, func(err Throwable){
    err.PrintStackTrace()
})()
// 抛异常
try.Throw(try.PointException.NewThrow("指针"))


// or


// 捕获异常  可以将 捕获异常放到 抛异常之前，可以在抛出异常之前 定义 异常的捕获，防止异常没有处理
defer try.CatchUncaughtException(func(err Throwable){
    err.PrintStackTrace()
})()
// 抛异常
try.Throw(try.PointException.NewThrow("指针"))

```
