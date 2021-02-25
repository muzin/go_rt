package try

import (
	"fmt"
	"testing"
	"time"
)

type LogicException struct {
	super *Exception
}

func DeclareLogicException(name string) *LogicException {
	return newLogicException(name)
}

func newLogicException(name string) *LogicException {
	exception := DeclareException(name)
	return &LogicException{super: exception}
}

func (this LogicException) New(msg string) *LogicException {
	exception := newLogicException(this.super.GetName())
	exception.super.msg = msg
	return exception
}

func (this LogicException) NewThrow(msg string) Throwable {
	exception := newLogicException(this.super.GetName())
	exception.super.msg = msg
	exception.super.stack = getStackTrace(exception)
	return exception
}

func (this LogicException) GetName() string {
	return this.super.GetName()
}

func (this LogicException) GetMsg() string {
	return this.super.GetMsg()
}

func (this LogicException) GetStackTrace() string {
	return this.super.GetStackTrace()
}

func (this LogicException) Error() string {
	return fmt.Sprintf("%v: %v", this.super.GetName(), this.super.GetMsg())
}

func (this LogicException) PrintStackTrace() {
	fmt.Printf("%v", this.super.GetStackTrace())
}

func TestThrowAndTryCatch(t *testing.T) {
	t.Run("TestThrowAndTryCatch", func(t *testing.T) {

		var ALogicException = DeclareException("ALogicException")

		now := time.Now()

		catched := false

		defer Catch(ALogicException, func(err Throwable) {
			err.PrintStackTrace()
			catched = true

			since := time.Since(now)
			t.Logf("捕获异常 耗时：%v", since)

			if catched == false {
				t.Logf("catched not")
			} else {
				t.Logf("catched")
			}

		})()
		Throw(ALogicException.NewThrow("Exception one"))

	})
}

func TestThrowAndTryCatch2(t *testing.T) {
	t.Run("TestThrowAndTryCatch2", func(t *testing.T) {

		var BLogicException = DeclareLogicException("BLogicException")

		now := time.Now()

		catched := false

		defer Catch(BLogicException, func(err Throwable) {
			err.PrintStackTrace()
			catched = true

			since := time.Since(now)
			t.Logf("捕获异常 耗时：%v", since)

			if catched == false {
				t.Logf("catched not")
			} else {
				t.Logf("catched")
			}

		})()
		Throw(BLogicException.NewThrow("Exception two"))

	})
}

func TestThrowAndTryCatchUncaughtException(t *testing.T) {
	t.Run("TestThrowAndTryCatchUncaughtException", func(t *testing.T) {

		var BLogicException = DeclareLogicException("BLogicException")

		now := time.Now()

		catched := false

		defer CatchUncaughtException(func(err Throwable) {
			err.PrintStackTrace()
			catched = true

			since := time.Since(now)
			t.Logf("捕获异常 耗时：%v", since)

			if catched == false {
				t.Logf("catched not")
			} else {
				t.Logf("catched")
			}

		})()
		Throw(BLogicException.NewThrow("Exception three"))

	})
}
