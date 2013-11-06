// Copyright 2013 Ardan Studios. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE handle.

/*
	Package TraceLog implements a logging system to trace all aspect of your code.
	The programmer should feel free to trace log as much of the code as possbile. This is our eyes to the running application.
	Logging System Based On The Log/Logger Standard Library

	Startup Options

	There are two options for starting the TraceLog:
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
	"text/template"
	"time"
)

//** CONSTANTS

// Support constants
const (
	_SYSTEM_ALERT_SUBJECT = "TraceLog Exception"
)

const (
	LEVEL_TRACE = 1
	LEVEL_INFO  = 2
	LEVEL_WARN  = 4
	LEVEL_ERROR = 8
)

//** NEW TYPES

// emailConfiguration contains configuration information required by the ConfigureEmailAlerts function
type emailConfiguration struct {
	Host     string
	Port     int
	UserName string
	Password string
	To       []string
	Auth     smtp.Auth
	Template *template.Template
}

// traceLog provides support to write to log files
type traceLog struct {
	Serialize          sync.Mutex
	EmailConfiguration *emailConfiguration
	TRACE              *log.Logger
	INFO               *log.Logger
	WARN               *log.Logger
	ERROR              *log.Logger
	FILE               *log.Logger
	LogFile            *os.File
}

//** SINGLETON REFERENCE

var _This *traceLog // A reference to the singleton

//** PUBLIC FUNCTIONS

// Called to init the logging system
func init() {
	log.SetPrefix("TRACE: ")
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
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

// Start initializes tracelog for all logging levels
func Start(logLevel int) {

	turnOnLogging(logLevel, nil)
}

// Start initializes tracelog for all logging levels plus file
func StartFile(logLevel int, baseFilePath string, daysToKeep int) {
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

// SendEmailException will send an email along with the exception
func SendEmailException(subject string, message string, a ...interface{}) (err error) {
	defer _This.CatchPanic(&err, "SendEmailException")

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

//** PRIVATE FUNCTIONS

// turnOnLogging configures the logging writers
func turnOnLogging(logLevel int, fileHandle io.Writer) {
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
}

//** PRIVATE MEMBER FUNCTIONS

// LogDirectoryCleanup performs all the directory cleanup and maintenance
func (this *traceLog) LogDirectoryCleanup(baseFilePath string, daysToKeep int) {
	defer this.CatchPanic(nil, "LogDirectoryCleanup")

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
func (this *traceLog) CatchPanic(err *error, functionName string) {
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
func (this *traceLog) EmailScript() (script string) {
	return `From: {{.From}}
To: {{.To}}
Subject: {{.Subject}}
MIME-version: 1.0
Content-Type: text/html; charset="UTF-8"

<html><body>{{.Message}}</body></html>`
}
