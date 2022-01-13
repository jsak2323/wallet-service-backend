package error

import (
	"fmt"
	"runtime"
)

func (err *Error) Error() string {
	return err.Title
}

func (err *Error) GetMessage() string {
	return err.Message
}

func (err *Error) GetTrace() string {
	return err.Trace
}

func AddTrace(err error) *Error {
	var (
		trace string
	)

	if newError, ok := err.(*Error); ok {
		pc, file, ln, ok := runtime.Caller(1)
		if ok {
			trace = fmt.Sprintf("%s(%d) %s\n%s", file, ln, runtime.FuncForPC(pc).Name(), newError.GetTrace())
		}
		return NewError(newError.Error(), newError.GetMessage(), trace)
	}

	if newErr, ok := err.(error); ok {
		pc, file, ln, ok := runtime.Caller(1)
		if ok {

			trace = fmt.Sprintf("%s(%d) %s\n", file, ln, runtime.FuncForPC(pc).Name())
		}

		return NewError(newErr.Error(), newErr.Error(), trace)
	}

	return NewError("unrecognized error type", "", "")
}

func AssignErr(from error, to *Error) *Error {
	if fromError, ok := from.(*Error); ok {
		return NewError(to.Error(), fromError.GetMessage(), fromError.GetTrace())
	}

	if _, ok := from.(error); ok {
		return NewError(to.Error(), from.Error(), "")
	}

	return NewError("unrecognized error type", "", "")
}

func Logged(err error) string {
	if newErr, ok := err.(*Error); ok {
		return fmt.Sprintf("Title: %s\nMessage: %s\nTrace: \n%s\n", newErr.Error(), newErr.GetMessage(), newErr.GetTrace())
	}

	return fmt.Sprintf("Title: %s\nMessage: %s\nTrace: \n%s\n", InternalServerErr.Error(), InternalServerErr.GetMessage(), InternalServerErr.GetTrace())
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
