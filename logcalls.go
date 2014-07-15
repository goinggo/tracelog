// Copyright 2013 Ardan Studios. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE handle.

// Package tracelog : logcalls.go provides formatting functions.
package tracelog

import (
	"fmt"
)

//** STARTED AND COMPLETED

// Started uses the Serialize destination and adds a Started tag to the log line
func Started(title string, functionName string) {
	logger.Trace.Output(2, fmt.Sprintf("%s : %s : Started\n", title, functionName))
}

// Startedf uses the Serialize destination and writes a Started tag to the log line
func Startedf(title string, functionName string, format string, a ...interface{}) {
	logger.Trace.Output(2, fmt.Sprintf("%s : %s : Started : %s\n", title, functionName, fmt.Sprintf(format, a...)))
}

// Completed uses the Serialize destination and writes a Completed tag to the log line
func Completed(title string, functionName string) {
	logger.Trace.Output(2, fmt.Sprintf("%s : %s : Completed\n", title, functionName))
}

// Completedf uses the Serialize destination and writes a Completed tag to the log line
func Completedf(title string, functionName string, format string, a ...interface{}) {
	logger.Trace.Output(2, fmt.Sprintf("%s : %s : Completed : %s\n", title, functionName, fmt.Sprintf(format, a...)))
}

// CompletedError uses the Error destination and writes a Completed tag to the log line
func CompletedError(err error, title string, functionName string) {
	logger.Error.Output(2, fmt.Sprintf("%s : %s : Completed : ERROR : %s\n", title, functionName, err))
}

// CompletedErrorf uses the Error destination and writes a Completed tag to the log line
func CompletedErrorf(err error, title string, functionName string, format string, a ...interface{}) {
	logger.Error.Output(2, fmt.Sprintf("%s : %s : Completed : ERROR : %s : %s\n", title, functionName, fmt.Sprintf(format, a...), err))
}

//** TRACE

// Trace writes to the Trace destination
func Trace(title string, functionName string, format string, a ...interface{}) {
	logger.Trace.Output(2, fmt.Sprintf("%s : %s : Info : %s\n", title, functionName, fmt.Sprintf(format, a...)))
}

//** INFO

// Info writes to the Info destination
func Info(title string, functionName string, format string, a ...interface{}) {
	logger.Info.Output(2, fmt.Sprintf("%s : %s : Info : %s\n", title, functionName, fmt.Sprintf(format, a...)))
}

//** WARNING

// Warning writes to the Warning destination
func Warning(title string, functionName string, format string, a ...interface{}) {
	logger.Warning.Output(2, fmt.Sprintf("%s : %s : Info : %s\n", title, functionName, fmt.Sprintf(format, a...)))
}

//** ERROR

// Error writes to the Error destination and accepts an err
func Error(err error, title string, functionName string) {
	logger.Error.Output(2, fmt.Sprintf("%s : %s : ERROR : %s\n", title, functionName, err))
}

// Errorf writes to the Error destination and accepts an err
func Errorf(err error, title string, functionName string, format string, a ...interface{}) {
	logger.Error.Output(2, fmt.Sprintf("%s : %s : ERROR : %s : %s\n", title, functionName, fmt.Sprintf(format, a...), err))
}

//** ALERT

// Alert write to the Error destination and sends email alert
func Alert(subject string, title string, functionName string, format string, a ...interface{}) {
	message := fmt.Sprintf("%s : %s : ALERT : %s\n", title, functionName, fmt.Sprintf(format, a...))
	logger.Error.Output(2, message)
	SendEmailException(subject, message)
}

// CompletedAlert write to the Error destination, writes a Completed tag to the log line and sends email alert
func CompletedAlert(subject string, title string, functionName string, format string, a ...interface{}) {
	message := fmt.Sprintf("%s : %s : Completed : ALERT : %s\n", title, functionName, fmt.Sprintf(format, a...))
	logger.Error.Output(2, message)
	SendEmailException(subject, message)
}
