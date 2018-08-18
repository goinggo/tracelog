# Tracelog

Tracelog sits on top of [Go's log library](https://golang.org/pkg/log/). See the following example for an implementation (it's super easy to use)

### Example
```go
package main

import (
    "fmt"
    "errors"
    "github.com/goinggo/tracelog"
)

func main() {
    //tracelog.StartFile(tracelog.LevelTrace, "/Users/bill/Temp/logs", 1)
    
    tracelog.Start(tracelog.LevelTrace)
    tracelog.Trace("main", "main", "Hello Trace")
    tracelog.Info("main", "main", "Hello Info")
    tracelog.Warning("main", "main", "Hello Warn")
    tracelog.Errorf(fmt.Errorf("Exception At..."), "main", "main", "Hello Error")
    
    Example()
    tracelog.Stop()
}

func foo() error {
        return errors.New("test")
}

func Example() {
    tracelog.Started("main", "Example")
    
    if err := foo(); err != nil {
        tracelog.CompletedError(err, "main", "Example")
        return
    }
    
    tracelog.Completed("main", "Example")
}
```

#### Output

```sh
TRACE: 2013/11/07 08:24:32 main.go:12: main : main : Info : Hello Trace
INFO: 2013/11/07 08:24:32 main.go:13: main : main : Info : Hello Info
WARNING: 2013/11/07 08:24:32 main.go:14: main : main : Info : Hello Warn
ERROR: 2013/11/07 08:24:32 main.go:15: main : main : Info : Hello Error : Exception At...
TRACE: 2013/11/07 08:24:32 main.go:23: main : Example : Started
TRACE: 2013/11/07 08:24:32 main.go:31: main : Example : Completed
TRACE: 2013/11/07 08:24:32 tracelog.go:149: main : Stop : Started
TRACE: 2013/11/07 08:24:32 tracelog.go:156: main : Stop : Completed
```


For more details on how to use this package (or to see how it works), see the source. Also, [docs](http://godoc.org/github.com/goinggo/tracelog)

Copyright 2013 Ardan Studios. All rights reserved.  
Use of this source code is governed by a BSD-style license that can be found in the LICENSE handle.

Package TraceLog implements a logging system to trace all aspect of your code. This is great for task oriented programs.	Based on the Go log standard library. It provides 4 destinations with logging levels plus you can attach a file for persistent writes. A log clean process is provided to maintain disk space. There is also email support to send email alerts.

Ardan Studios  
12973 SW 112 ST, Suite 153  
Miami, FL 33186  
bill@ardanstudios.com

[Click To View Documentation](http://godoc.org/github.com/goinggo/tracelog)


