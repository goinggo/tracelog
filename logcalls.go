// Copyright 2013 Ardan Studios. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE handle.

package tracelog

import (
	"fmt"
)

//** Started and Completed

// STARTED uses the TRACE destination and adds a Started tag to the log line
func STARTED(routineName string, functionName string) {
	_This.Serialize.Lock()
	defer _This.Serialize.Unlock()
	_This.TRACE.Output(2, fmt.Sprintf("%s : %s : Started\n", routineName, functionName))
}

// STARTEDf uses the TRACE destination and writes a Started tag to the log line
func STARTEDf(routineName string, functionName string, format string, a ...interface{}) {
	_This.Serialize.Lock()
	defer _This.Serialize.Unlock()
	_This.TRACE.Output(2, fmt.Sprintf("%s : %s : Started : %s\n", routineName, functionName, fmt.Sprintf(format, a...)))
}

// COMPLETED uses the TRACE destination and writes a Completed tag to the log line
func COMPLETED(routineName string, functionName string) {
	_This.Serialize.Lock()
	defer _This.Serialize.Unlock()
	_This.TRACE.Output(2, fmt.Sprintf("%s : %s : Completed\n", routineName, functionName))
}

// COMPLETEDf uses the TRACE destination and writes a Completed tag to the log line
func COMPLETEDf(routineName string, functionName string, format string, a ...interface{}) {
	_This.Serialize.Lock()
	defer _This.Serialize.Unlock()
	_This.TRACE.Output(2, fmt.Sprintf("%s : %s : Completed : %s\n", routineName, functionName, fmt.Sprintf(format, a...)))
}

// COMPLETED_ERROR uses the ERROR destination and writes a Completed tag to the log line
func COMPLETED_ERROR(err error, routineName string, functionName string) {
	_This.Serialize.Lock()
	defer _This.Serialize.Unlock()
	_This.ERROR.Output(2, fmt.Sprintf("%s : %s : Completed : %s\n", routineName, functionName, err))
}

// COMPLETED_ERRORf uses the ERROR destination and writes a Completed tag to the log line
func COMPLETED_ERRORf(err error, routineName string, functionName string, format string, a ...interface{}) {
	_This.Serialize.Lock()
	defer _This.Serialize.Unlock()
	_This.ERROR.Output(2, fmt.Sprintf("%s : %s : Completed : %s : %s\n", routineName, functionName, fmt.Sprintf(format, a...), err))
}

//** TRACE

// TRACE writes to the TRACE destination
func TRACE(routineName string, functionName string, format string, a ...interface{}) {
	_This.Serialize.Lock()
	defer _This.Serialize.Unlock()
	_This.TRACE.Output(2, fmt.Sprintf("%s : %s : Info : %s\n", routineName, functionName, fmt.Sprintf(format, a...)))
}

//** INFO

// INFO writes to the INFO destination
func INFO(routineName string, functionName string, format string, a ...interface{}) {
	_This.Serialize.Lock()
	defer _This.Serialize.Unlock()
	_This.INFO.Output(2, fmt.Sprintf("%s : %s : Info : %s\n", routineName, functionName, fmt.Sprintf(format, a...)))
}

//** WARN

// WARN writes to the WARN destination
func WARN(routineName string, functionName string, format string, a ...interface{}) {
	_This.Serialize.Lock()
	defer _This.Serialize.Unlock()
	_This.WARN.Output(2, fmt.Sprintf("%s : %s : Info : %s\n", routineName, functionName, fmt.Sprintf(format, a...)))
}

//** ERROR

// ERROR writes to the ERROR destination and accepts an err
func ERROR(err error, routineName string, functionName string) {
	_This.Serialize.Lock()
	defer _This.Serialize.Unlock()
	_This.ERROR.Output(2, fmt.Sprintf("%s : %s : Info : %s\n", routineName, functionName, err))
}

// ERROR writes to the ERROR destination and accepts an err
func ERRORf(err error, routineName string, functionName string, format string, a ...interface{}) {
	_This.Serialize.Lock()
	defer _This.Serialize.Unlock()
	_This.ERROR.Output(2, fmt.Sprintf("%s : %s : Info : %s : %s\n", routineName, functionName, fmt.Sprintf(format, a...), err))
}

//** ALERT

// ALERT write to the ERROR destination and sends email alert
func ALERT(subject string, routineName string, functionName string, format string, a ...interface{}) {
	message := fmt.Sprintf("%s : %s : Info : %s\n", routineName, functionName, fmt.Sprintf(format, a...))

	_This.Serialize.Lock()
	defer _This.Serialize.Unlock()
	_This.ERROR.Output(2, message)

	SendEmailException(subject, message)
}

// COMPLETED_ALERT write to the ERROR destination, writes a Completed tag to the log line and sends email alert
func COMPLETED_ALERT(subject string, routineName string, functionName string, format string, a ...interface{}) {
	message := fmt.Sprintf("%s : %s : Completed : %s\n", routineName, functionName, fmt.Sprintf(format, a...))

	_This.Serialize.Lock()
	defer _This.Serialize.Unlock()
	_This.ERROR.Output(2, message)

	SendEmailException(subject, message)
}
