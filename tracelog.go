// Copyright 2013 Ardan Studios. All rights reserved.
// Use of traceLog source code is governed by a BSD-style
// license that can be found in the LICENSE handle.

/*
	Package TraceLog implements a logging system to trace all aspect of your code. This is great for task oriented programs.
	Based on the Go log standard library. It provides 4 destinations with logging levels plus you can attach a file for persistent
	writes. A log clean process is provided to maintain disk space. There is also email support to send email alerts.

		Read the following blog post for more information:
		http://www.goinggo.net/2013/11/using-log-package-in-go.html

	Example Code:

		package main

		import (
		    "fmt"
		    "github.com/goinggo/tracelog"
		)

		func main() {
		    //tracelog.StartFile(tracelog.LEVEL_TRACE, "/Users/bill/Temp/logs", 1)
		    tracelog.Start(tracelog.LEVEL_TRACE)

		    tracelog.TRACE("main", "main", "Hello Trace")
		    tracelog.INFO("main", "main", "Hello Info")
		    tracelog.WARN("main", "main", "Hello Warn")
		    tracelog.ERRORf(fmt.Errorf("Exception At..."), "main", "main", "Hello Error")

		    Example()

		    tracelog.Stop()
		}

		func Example() {
		    tracelog.STARTED("main", "Example")

		    err := foo()
		    if err != nil {
		        tracelog.COMPLETED_ERROR(err, "main", "Example")
		        return
		    }

		    tracelog.COMPLETED("main", "Example")
		}

	Output:

		TRACE: 2013/11/07 08:24:32 main.go:12: main : main : Info : Hello Trace
		INFO: 2013/11/07 08:24:32 main.go:13: main : main : Info : Hello Info
		WARNING: 2013/11/07 08:24:32 main.go:14: main : main : Info : Hello Warn
		ERROR: 2013/11/07 08:24:32 main.go:15: main : main : Info : Hello Error : Exception At...

		TRACE: 2013/11/07 08:24:32 main.go:23: main : Example : Started
		TRACE: 2013/11/07 08:24:32 main.go:31: main : Example : Completed

		TRACE: 2013/11/07 08:24:32 tracelog.go:149: main : Stop : Started
		TRACE: 2013/11/07 08:24:32 tracelog.go:156: main : Stop : Completed
*/
package tracelog

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/smtp"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"text/template"
	"time"
)

//** CONSTANTS

const (
	_SYSTEM_ALERT_SUBJECT = "TraceLog Exception"
)

const (
	LEVEL_TRACE int32 = 1 // Log everything
	LEVEL_INFO  int32 = 2 // Log Info, Warnings and Errors
	LEVEL_WARN  int32 = 4 // Log Warning and Errors
	LEVEL_ERROR int32 = 8 // Log just Errors
)

//** PACKAGE VARIABLES

var (
	_This *traceLog // A reference to the singleton
)

//** TYPES

type (
	// emailConfiguration contains configuration information required by the ConfigureEmailAlerts function
	emailConfiguration struct {
		Host     string
		Port     int
		UserName string
		Password string
		To       []string
		Auth     smtp.Auth
		Template *template.Template
	}

	// traceLog provides support to write to log files
	traceLog struct {
		LogLevel           int32
		Serialize          sync.Mutex
		EmailConfiguration *emailConfiguration
		TRACE              *log.Logger
		INFO               *log.Logger
		WARN               *log.Logger
		ERROR              *log.Logger
		FILE               *log.Logger
		LogFile            *os.File
	}
)

//** INIT FUNCTION

// Called to init the logging system
func init() {
	log.SetPrefix("TRACE: ")
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

//** PUBLIC FUNCTIONS

// Start initializes tracelog and only displays the specified logging level
func Start(logLevel int32) {
	turnOnLogging(logLevel, nil)
}

// StartFile initializes tracelog and only displays the specified logging level
// and creates a file to capture writes
func StartFile(logLevel int32, baseFilePath string, daysToKeep int) {
	baseFilePath = strings.TrimRight(baseFilePath, "/")
	currentDate := time.Now().UTC()
	dateDirectory := time.Now().UTC().Format("2006-01-02")
	dateFile := currentDate.Format("2006-01-02T15-04-05")

	filePath := fmt.Sprintf("%s/%s/", baseFilePath, dateDirectory)
	fileName := strings.Replace(fmt.Sprintf("%s.txt", dateFile), " ", "-", -1)

	err := os.MkdirAll(filePath, os.ModePerm)
	if err != nil {
		log.Fatalf("main : Start : Failed to Create log directory : %s : %s\n", filePath, err)
	}

	logf, err := os.Create(fmt.Sprintf("%s%s", filePath, fileName))
	if err != nil {
		log.Fatalf("main : Start : Failed to Create log file : %s : %s\n", fileName, err)
	}

	// Turn the logging on
	turnOnLogging(logLevel, logf)

	// Cleanup any existing directories
	_This.LogDirectoryCleanup(baseFilePath, daysToKeep)
}

// Stop will release resources and shutdown all processing
func Stop() (err error) {
	STARTED("main", "Stop")

	if _This.LogFile != nil {
		TRACE("main", "Stop", "Closing File")
		err = _This.LogFile.Close()
	}

	COMPLETED("main", "Stop")

	return err
}

// ConfigureEmail configures the email system for use
func ConfigureEmail(host string, port int, userName string, password string, to []string) {
	emailConfiguration := &emailConfiguration{
		Host:     host,
		Port:     port,
		UserName: userName,
		Password: password,
		To:       to,
		Auth:     smtp.PlainAuth("", userName, password, host),
		Template: template.Must(template.New("emailTemplate").Parse(_This.EmailScript())),
	}

	_This.EmailConfiguration = emailConfiguration
}

// SendEmailException will send an email along with the exception
func SendEmailException(subject string, message string, a ...interface{}) (err error) {
	defer _This.CatchPanic(&err, "SendEmailException")

	if _This.EmailConfiguration == nil {
		return
	}

	parameters := &struct {
		From    string
		To      string
		Subject string
		Message string
	}{
		_This.EmailConfiguration.UserName,
		strings.Join([]string(_This.EmailConfiguration.To), ","),
		subject,
		fmt.Sprintf(message, a...),
	}

	emailMessage := new(bytes.Buffer)
	_This.EmailConfiguration.Template.Execute(emailMessage, parameters)

	err = smtp.SendMail(fmt.Sprintf("%s:%d", _This.EmailConfiguration.Host, _This.EmailConfiguration.Port), _This.EmailConfiguration.Auth, _This.EmailConfiguration.UserName, _This.EmailConfiguration.To, emailMessage.Bytes())

	return err
}

// LogLevel returns the configured logging level
func LogLevel() int32 {
	return atomic.LoadInt32(&_This.LogLevel)
}

//** PRIVATE FUNCTIONS

// turnOnLogging configures the logging writers
func turnOnLogging(logLevel int32, fileHandle io.Writer) {
	traceHandle := ioutil.Discard
	infoHandle := ioutil.Discard
	warnHandle := ioutil.Discard
	errorHandle := ioutil.Discard

	if logLevel&LEVEL_TRACE != 0 {
		traceHandle = os.Stdout
		infoHandle = os.Stdout
		warnHandle = os.Stdout
		errorHandle = os.Stderr
	}

	if logLevel&LEVEL_INFO != 0 {
		infoHandle = os.Stdout
		warnHandle = os.Stdout
		errorHandle = os.Stderr
	}

	if logLevel&LEVEL_WARN != 0 {
		warnHandle = os.Stdout
		errorHandle = os.Stderr
	}

	if logLevel&LEVEL_ERROR != 0 {
		errorHandle = os.Stderr
	}

	if fileHandle != nil {
		if traceHandle == os.Stdout {
			traceHandle = io.MultiWriter(fileHandle, traceHandle)
		}

		if infoHandle == os.Stdout {
			infoHandle = io.MultiWriter(fileHandle, infoHandle)
		}

		if warnHandle == os.Stdout {
			warnHandle = io.MultiWriter(fileHandle, warnHandle)
		}

		if errorHandle == os.Stderr {
			errorHandle = io.MultiWriter(fileHandle, errorHandle)
		}
	}

	_This = &traceLog{
		TRACE: log.New(traceHandle, "TRACE: ", log.Ldate|log.Ltime|log.Lshortfile),
		INFO:  log.New(infoHandle, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile),
		WARN:  log.New(warnHandle, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile),
		ERROR: log.New(errorHandle, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
	}

	atomic.StoreInt32(&_This.LogLevel, logLevel)
}

//** PRIVATE MEMBER FUNCTIONS

// LogDirectoryCleanup performs all the directory cleanup and maintenance
func (traceLog *traceLog) LogDirectoryCleanup(baseFilePath string, daysToKeep int) {
	defer traceLog.CatchPanic(nil, "LogDirectoryCleanup")

	STARTEDf("main", "LogDirectoryCleanup", "BaseFilePath[%s] DaysToKeep[%d]", baseFilePath, daysToKeep)

	// Get a list of existing directories
	fileInfos, err := ioutil.ReadDir(baseFilePath)
	if err != nil {
		COMPLETED_ERROR(err, "main", "LogDirectoryCleanup")
		return
	}

	// Create the date to compare for directories to remove
	currentDate := time.Now().UTC()
	compareDate := time.Date(currentDate.Year(), currentDate.Month(), currentDate.Day()-daysToKeep, 0, 0, 0, 0, time.UTC)

	TRACE("main", "LogDirectoryCleanup", "CompareDate[%v]", compareDate)

	for _, fileInfo := range fileInfos {
		if fileInfo.IsDir() == false {
			continue
		}

		// The file name look like: YYYY-MM-DD
		parts := strings.Split(fileInfo.Name(), "-")

		year, err := strconv.Atoi(parts[0])
		if err != nil {
			ERRORf(err, "main", "LogDirectoryCleanup", "Attempting To Convert Directory [%s]", fileInfo.Name())
			continue
		}

		month, err := strconv.Atoi(parts[1])
		if err != nil {
			ERRORf(err, "main", "LogDirectoryCleanup", "Attempting To Convert Directory [%s]", fileInfo.Name())
			continue
		}

		day, err := strconv.Atoi(parts[2])
		if err != nil {
			ERRORf(err, "main", "LogDirectoryCleanup", "Attempting To Convert Directory [%s]", fileInfo.Name())
			continue
		}

		// The directory to check
		fullFileName := fmt.Sprintf("%s/%s", baseFilePath, fileInfo.Name())

		// Create a time type from the directory name
		directoryDate := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)

		// Compare the dates and convert to days
		daysOld := int(compareDate.Sub(directoryDate).Hours() / 24)

		TRACE("main", "LogDirectoryCleanup", "Checking Directory[%s] DaysOld[%d]", fullFileName, daysOld)

		if daysOld >= 0 {
			TRACE("main", "LogDirectoryCleanup", "Removing Directory[%s]", fullFileName)

			err = os.RemoveAll(fullFileName)

			if err != nil {
				TRACE("main", "LogDirectoryCleanup", "Attempting To Remove Directory [%s]", fullFileName)
				continue
			}

			TRACE("main", "LogDirectoryCleanup", "Directory Removed [%s]", fullFileName)
		}
	}

	// We don't need the catch handler to log any errors
	err = nil

	COMPLETED("main", "LogDirectoryCleanup")
	return
}

// CatchPanic is used to catch any Panic and log exceptions to Stdout. It will also write the stack trace
func (traceLog *traceLog) CatchPanic(err *error, functionName string) {
	if r := recover(); r != nil {

		// Capture the stack trace
		buf := make([]byte, 10000)
		runtime.Stack(buf, false)

		SendEmailException(_SYSTEM_ALERT_SUBJECT, "%s : PANIC Defered [%s] : Stack Trace : %s", functionName, r, string(buf))

		if err != nil {
			*err = fmt.Errorf("%v", r)
		}
	}
}

// EmailScript returns a template for the email message to be sent
func (traceLog *traceLog) EmailScript() (script string) {
	return `From: {{.From}}
To: {{.To}}
Subject: {{.Subject}}
MIME-version: 1.0
Content-Type: text/html; charset="UTF-8"

<html><body>{{.Message}}</body></html>`
}
