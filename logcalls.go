package tracelog

import (
	"fmt"
)

//** Started and Completed

func STARTED(routineName string, functionName string) {
	_This.Serialize.Lock()
	defer _This.Serialize.Unlock()
	_This.TRACE.Output(2, fmt.Sprintf("%s : %s : Started\n", routineName, functionName))
}

func STARTEDf(routineName string, functionName string, format string, a ...interface{}) {
	_This.Serialize.Lock()
	defer _This.Serialize.Unlock()
	_This.TRACE.Output(2, fmt.Sprintf("%s : %s : Started : %s\n", routineName, functionName, fmt.Sprintf(format, a...)))
}

func COMPLETED(routineName string, functionName string) {
	_This.Serialize.Lock()
	defer _This.Serialize.Unlock()
	_This.TRACE.Output(2, fmt.Sprintf("%s : %s : Completed\n", routineName, functionName))
}

func COMPLETEDf(routineName string, functionName string, format string, a ...interface{}) {
	_This.Serialize.Lock()
	defer _This.Serialize.Unlock()
	_This.TRACE.Output(2, fmt.Sprintf("%s : %s : Completed : %s\n", routineName, functionName, fmt.Sprintf(format, a...)))
}

func COMPLETED_ERROR(err error, routineName string, functionName string) {
	_This.Serialize.Lock()
	defer _This.Serialize.Unlock()
	_This.ERROR.Output(2, fmt.Sprintf("%s : %s : Completed : %s\n", err))
}

func COMPLETED_ERRORf(err error, routineName string, functionName string, format string, a ...interface{}) {
	_This.Serialize.Lock()
	defer _This.Serialize.Unlock()
	_This.ERROR.Output(2, fmt.Sprintf("%s : %s : Completed : %s : %s\n", routineName, functionName, fmt.Sprintf(format, a...), err))
}

//** TRACE

func TRACE(routineName string, functionName string, format string, a ...interface{}) {
	_This.Serialize.Lock()
	defer _This.Serialize.Unlock()
	_This.TRACE.Output(2, fmt.Sprintf("%s : %s : Info : %s\n", routineName, functionName, fmt.Sprintf(format, a...)))
}

//** INFO

func INFO(routineName string, functionName string, format string, a ...interface{}) {
	_This.Serialize.Lock()
	defer _This.Serialize.Unlock()
	_This.INFO.Output(2, fmt.Sprintf("%s : %s : Info : %s\n", routineName, functionName, fmt.Sprintf(format, a...)))
}

//** WARN

func WARN(routineName string, functionName string, format string, a ...interface{}) {
	_This.Serialize.Lock()
	defer _This.Serialize.Unlock()
	_This.WARN.Output(2, fmt.Sprintf("%s : %s : Info : %s\n", routineName, functionName, fmt.Sprintf(format, a...)))
}

//** ERROR

func ERROR(err error, routineName string, functionName string) {
	_This.Serialize.Lock()
	defer _This.Serialize.Unlock()
	_This.ERROR.Output(2, fmt.Sprintf("%s : %s : Info : %s\n", err))
}

func ERRORf(err error, routineName string, functionName string, format string, a ...interface{}) {
	_This.Serialize.Lock()
	defer _This.Serialize.Unlock()
	_This.ERROR.Output(2, fmt.Sprintf("%s : %s : Info : %s : %s\n", routineName, functionName, fmt.Sprintf(format, a...), err))
}
