package error

import (
	"fmt"
	"reflect"
	"runtime"
)

func GetFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

func (err Error) Error() string {
	return err.Title
}

func (err Error) GetMessage() string {
	return err.Message
}

func (err Error) GetTrace() string {
	return err.Trace
}

func AddTrace(err error, functionName string) Error {
	var (
		trace string
	)

	if newError, ok := err.(Error); ok {
		_, file, ln, ok := runtime.Caller(1)
		if ok {
			trace = fmt.Sprintf("%s(%d) %s\n%s", file, ln, functionName, newError.GetTrace())
		}

		return NewError(newError.Error(), newError.GetMessage(), trace)
	}

	if newErr, ok := err.(error); ok {
		_, file, ln, ok := runtime.Caller(1)
		if ok {

			trace = fmt.Sprintf("%s(%d) %s\n", file, ln, functionName)
		}

		return NewError(newErr.Error(), newErr.Error(), trace)
	}

	return NewError("unrecognized error type", "", "")
}

func AssignErr(from error, to Error) Error {
	if fromError, ok := from.(Error); ok {
		return NewError(to.Error(), fromError.GetMessage(), fromError.GetTrace())
	}

	if _, ok := from.(error); ok {
		return NewError(to.Error(), from.Error(), "")
	}

	return NewError("unrecognized error type", "", "")
}

func logged(err error) Error {
	if newErr, ok := err.(Error); ok {
		fmt.Println(fmt.Sprintf("Title: %s\nMessage: %s\nTrace: \n%s\n", newErr.Error(), newErr.GetMessage(), newErr.GetTrace()))
		res := Error{
			Title:   newErr.Title,
			Message: newErr.Message,
			Trace:   newErr.Trace,
		}
		return res
	}

	return InternalServerErr
}

// func func2() error {
// 	errFunc1 := func1()
// 	funcName := GetFunctionName(func2)
// 	return AddTrace(errFunc1, funcName)
// }

// func func1() error {
// 	funcName := GetFunctionName(func1)
// 	return AddTrace(errors.New("source error message"), funcName)
// }

// func main() {
// 	fmt.Println("--- Log Traced Error ---")
// 	errFunc2 := func2()
// 	fmt.Println(errFunc2.Error(), "\n")
// 	logged(errFunc2)
// 	log.Println(logged(errFunc2))

// 	fmt.Println("--- Assign New Error and Log Assigned Error ---")
// 	assignedErr := AssignErr(errFunc2, internalServerErr)
// 	fmt.Println(assignedErr.Error(), "\n")
// 	logged(assignedErr)
// }
