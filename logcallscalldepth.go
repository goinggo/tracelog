// Copyright 2013 Ardan Studios. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE handle.

package tracelog

import (
	"fmt"
)

//** STARTED AND COMPLETED

// STARTEDcd uses the TRACE destination and adds a Started tag to the log line
func STARTEDcd(callDepth int, routineName string, functionName string) {
	_This.Serialize.Lock()
	defer _This.Serialize.Unlock()
	_This.TRACE.Output(callDepth, fmt.Sprintf("%s : %s : Started\n", routineName, functionName))
}

// STARTEDfcd uses the TRACE destination and writes a Started tag to the log line
func STARTEDfcd(callDepth int, routineName string, functionName string, format string, a ...interface{}) {
	_This.Serialize.Lock()
	defer _This.Serialize.Unlock()
	_This.TRACE.Output(callDepth, fmt.Sprintf("%s : %s : Started : %s\n", routineName, functionName, fmt.Sprintf(format, a...)))
}

// COMPLETEDcd uses the TRACE destination and writes a Completed tag to the log line
func COMPLETEDcd(callDepth int, routineName string, functionName string) {
	_This.Serialize.Lock()
	defer _This.Serialize.Unlock()
	_This.TRACE.Output(callDepth, fmt.Sprintf("%s : %s : Completed\n", routineName, functionName))
}

// COMPLETEDfcd uses the TRACE destination and writes a Completed tag to the log line
func COMPLETEDfcd(callDepth int, routineName string, functionName string, format string, a ...interface{}) {
	_This.Serialize.Lock()
	defer _This.Serialize.Unlock()
	_This.TRACE.Output(callDepth, fmt.Sprintf("%s : %s : Completed : %s\n", routineName, functionName, fmt.Sprintf(format, a...)))
}

// COMPLETED_ERRORcd uses the ERROR destination and writes a Completed tag to the log line
func COMPLETED_ERRORcd(callDepth int, err error, routineName string, functionName string) {
	_This.Serialize.Lock()
	defer _This.Serialize.Unlock()
	_This.ERROR.Output(callDepth, fmt.Sprintf("%s : %s : Completed : ERROR : %s\n", routineName, functionName, err))
}

// COMPLETED_ERRORfcd uses the ERROR destination and writes a Completed tag to the log line
func COMPLETED_ERRORfcd(callDepth int, err error, routineName string, functionName string, format string, a ...interface{}) {
	_This.Serialize.Lock()
	defer _This.Serialize.Unlock()
	_This.ERROR.Output(callDepth, fmt.Sprintf("%s : %s : Completed : ERROR : %s : %s\n", routineName, functionName, fmt.Sprintf(format, a...), err))
}

//** TRACE

// TRACEcd writes to the TRACE destination
func TRACEcd(callDepth int, routineName string, functionName string, format string, a ...interface{}) {
	_This.Serialize.Lock()
	defer _This.Serialize.Unlock()
	_This.TRACE.Output(callDepth, fmt.Sprintf("%s : %s : Info : %s\n", routineName, functionName, fmt.Sprintf(format, a...)))
}

//** INFO

// INFOcd writes to the INFO destination
func INFOcd(callDepth int, routineName string, functionName string, format string, a ...interface{}) {
	_This.Serialize.Lock()
	defer _This.Serialize.Unlock()
	_This.INFO.Output(callDepth, fmt.Sprintf("%s : %s : Info : %s\n", routineName, functionName, fmt.Sprintf(format, a...)))
}

//** WARN

// WARNcd writes to the WARN destination
func WARNcd(callDepth int, routineName string, functionName string, format string, a ...interface{}) {
	_This.Serialize.Lock()
	defer _This.Serialize.Unlock()
	_This.WARN.Output(callDepth, fmt.Sprintf("%s : %s : Info : %s\n", routineName, functionName, fmt.Sprintf(format, a...)))
}

//** ERROR

// ERRORcd writes to the ERROR destination and accepts an err
func ERRORcd(callDepth int, err error, routineName string, functionName string) {
	_This.Serialize.Lock()
	defer _This.Serialize.Unlock()
	_This.ERROR.Output(callDepth, fmt.Sprintf("%s : %s : ERROR : %s\n", routineName, functionName, err))
}

// ERRORfcd writes to the ERROR destination and accepts an err
func ERRORfcd(callDepth int, err error, routineName string, functionName string, format string, a ...interface{}) {
	_This.Serialize.Lock()
	defer _This.Serialize.Unlock()
	_This.ERROR.Output(callDepth, fmt.Sprintf("%s : %s : ERROR : %s : %s\n", routineName, functionName, fmt.Sprintf(format, a...), err))
}

//** ALERT

// ALERTcd write to the ERROR destination and sends email alert
func ALERTcd(callDepth int, subject string, routineName string, functionName string, format string, a ...interface{}) {
	message := fmt.Sprintf("%s : %s : ALERT : %s\n", routineName, functionName, fmt.Sprintf(format, a...))

	_This.Serialize.Lock()
	defer _This.Serialize.Unlock()
	_This.ERROR.Output(callDepth, message)

	SendEmailException(subject, message)
}

// COMPLETED_ALERTcd write to the ERROR destination, writes a Completed tag to the log line and sends email alert
func COMPLETED_ALERTcd(callDepth int, subject string, routineName string, functionName string, format string, a ...interface{}) {
	message := fmt.Sprintf("%s : %s : Completed : ALERT : %s\n", routineName, functionName, fmt.Sprintf(format, a...))

	_This.Serialize.Lock()
	defer _This.Serialize.Unlock()
	_This.ERROR.Output(callDepth, message)

	SendEmailException(subject, message)
}
