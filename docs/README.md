# Alog v0.6.2

(c) 2021 Gon Y Yi. <https://gonyyi.com>.  
[MIT License](https://raw.githubusercontent.com/gonyyi/alog/master/LICENSE)

Version 0.6.2

[![codecov](https://codecov.io/gh/gonyyi/alog/branch/master/graph/badge.svg?token=Y9RT0VRUQZ)](https://codecov.io/gh/gonyyi/alog)
[![Go Reference](https://pkg.go.dev/badge/github.com/gonyyi/alog.svg)](https://pkg.go.dev/github.com/gonyyi/alog@v0.6.2)
[![License](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://raw.githubusercontent.com/gonyyi/alog/master/LICENSE)

![Alog Screen Shot 1](https://github.com/gonyyi/alog/blob/master/docs/alog_screen_1.png)



## INTRODUCTION

Alog was built with a very simple goal in mind:

- Support Tagging (and also level)
- No memory allocation (or minimum allocation)
- Customizable

If you find any issues, please [create an issue](https://github.com/gonyyi/alog/issues/new).


[^Top](#alog)



## Example Usage

  ~~~go
  package main

  import (
    "github.com/gonyyi/alog"
    "github.com/gonyyi/alog/ext"
    "os"
  )

  func main() {
    al := alog.New(os.Stderr) // if nil is given for io.Writer, io.Discard will be used.
    al.Info(0).Write("info level log message")
    // Output:
    // {"date":20210304,"time":175516,"level":"info","tag":[],"message":"info level log message"}

    IO := al.NewTag("IO")
    DB := al.NewTag("DB")
    SYS := al.NewTag("SYS")
    NET := al.NewTag("NET")
    TEST := al.NewTag("TEST")

    al.Warn(0).Int("id", 1).Write("") // tag 0 means no tag is given.
    // Output: {"date":20210304,"time":180012,"level":"warn","tag":[],"id":1}

    al.Warn(IO).Int("id", 2).Write("this will print IO in the tag") // Use IO tag created above.
    // Output: {"date":20210304,"time":180012,"level":"warn","tag":["IO"],"message":"this will print IO in the tag","id":2}

    al.Warn(IO|DB|SYS|NET).Int("id", 3).Write("this will print IO/DB/SYS/NET to tag") // Use all tags by pipe
    // Output: {"date":20210304,"time":180012,"level":"warn","tag":["IO","DB","SYS","NET"],"message":"this will print IO/DB/SYS/NET to tag","id":3}

    al.Control.Level = alog.FatalLevel // Change loging level to Fatal or above.

    al.Info(IO|SYS).Int("id", 4).Write("this will not print") // No output as below Fatal level
    al.Info(IO).Int("id", 5).Write("this will not print")     // No output as below Fatal level
    al.Error(DB).Int("id", 6).Write("this will not print")    // No output as below Fatal level

    al.Fatal(NET).Int("id", 7).Write("this will print") // This will print because it's a Fatal level log entry
    // Output: {"date":20210304,"time":180150,"level":"fatal","tag":["NET"],"message":"this will print","id":7}

    al.Control.Tags = IO // By adding Tags to control, any log entries with Fatal level or above (as set above),
    // OR any log entries that contains IO tag will show.

    al.Info(IO|SYS). // Tag contains tag IO; will be printed
      Int("id", 4).
      Int("attempt", 2).
      Write("")

    al.Info(IO). // Tag is IO; will be printed
      Int("id", 5).
      Int("attempt", 2).
      Write("")

    al.Error(DB). // Tag IO isn't used. And level is below Fatal. Will NOT be printed.
      Int("id", 6).
      Int("attempt", 2).
      Write("")

    al.Fatal(NET). // Tag is not matching, but Level is, so will be printed.
      Int("id", 7).
      Int("attempt", 2).
      Write("")

    // Output
    // {"date":20210304,"time":180942,"level":"info","tag":["IO","SYS"],"id":4,"attempt":2}
    // {"date":20210304,"time":180942,"level":"info","tag":["IO"],"id":5,"attempt":2}
    // {"date":20210304,"time":180942,"level":"fatal","tag":["NET"],"id":7,"attempt":2}

    // EXTENSIONS: FORMATTER EXTENSION (github.com/gonyyi/alog/ext)
    al.Control.Level = alog.InfoLevel // revert default logging to INFO level.

    al = al.Ext(ext.LogFmt.TXTColor()) // ext.LogFmt.TXTColor() will set the formatter with color terminal output.

    al.Info(TEST).Str("testType", "colorText").Write("") // This will output the log with ANSI colored text format.
    // Output (Color): 2021-0304 18:13:38  INF  [TEST] testType="colorText"

    al = al.Ext(ext.LogFmt.TXT())
    al.Info(TEST).Str("testType", "normalText").Write("") // This will output the log with ANSI colored text format.
    // Output: 2021-0304 18:14:24 INF [TEST] testType="normalText"

    al = al.Ext(ext.LogFmt.NONE())
    al.Info(TEST).Str("testType", "backToJSON").Write("") // This will output the log with default JSON format.
    // Output: {"date":20210304,"time":181615,"level":"info","tag":["TEST"],"testType":"backToJSON"}


    // EXTENSION: Custom Entry
  	myEntry := func(s string) alog.EntryFn {
	  	return func(entry *alog.Entry) *alog.Entry {
		  	return entry.Str("name", s).Str("testStr", "myStr").Int("testInt", 123)
  		}
	  }
  	al.Info(0).Ext(myEntry("GON")).Write("ok")
    // Output: {"date":20210305,"time":81210,"level":"info","tag":[],"message":"ok","name":"GON","testStr":"myStr","testInt":123}
    
  
    // EXTENSION: MACRO LIKE EXTENSION (github.com/gonyyi/alog/ext)
    al = al.Ext(ext.LogMode.PROD("output.log")) // There are also DoMode.DEV(), DoMode.TEST().
    // When used with DoMode.TEST(), although it takes filename, it won't write it to file. It's just to make sure
    // a user can easily switch between TEST, DEV, and PROD mode.
    al.Info(TEST).Str("testType", "PROD").Write("") // This will write log into output.log file using buffered writer (bufio)
    al.Close() // `*Logger.Close() error` will close io.Writer if Close() method is available. Since DoMode.PROD uses buffered writer with
    // a Close method to flush the buffer, al.Close() is necessary.
  }
  ~~~


### Quick Start Example

`github.com/gonyyi/alog/log` _(see /log suffix)_, Alog can be used right away.

Default format is 
  `*Logger`.`[Trace/Debug/Info/Warn/Error/Fatal](Tag)`.`[Str/Int/Int64/Float/Bool/Err](string, VALUE)`.`Write(string)`


Example 1.

  ~~~go
  log.Info(MyTag).
      Str("city", "Gonway").Str("zip", "12345").
      Str("name", "Gon").Int("age", 50).
      Bool("isMarried", false).
      Float("height", 5.8).Write("") // Make sure all log entries are ending with "Write(string)"
      // Any string given to `Write(string)` will written as a "message" in default JSON format.
  ~~~


Example 2.

  ~~~go
  package main

  import (
    "github.com/gonyyi/alog"
    "github.com/gonyyi/alog/log"
  ) // import `.../alog/log` instead of `.../alog`

  func main() {
    // Alog format is
    //    *Logger.`LOGLEVEL(TAG)`.`Str/Int/Int64/Float/Bool/Err(key, value)`.Write(`OPTIONAL MSG`)
    log.Info(0).Str("name", "Alog").Int("buildNo", 6).Int("testID", 1).Write("Starting")

    // Tag can be created by `NewTag(name string)`
    tagDB := log.NewTag("DB")
    tagHTTP := log.NewTag("HTTP")
    tagREQ := log.NewTag("REQ")
    tagRES := log.NewTag("RES")

    // Below Debug will not be printed, because default log level is INFO or higher.
    log.Debug(tagDB).Str("status", "started").Int("buildNo", 6).Int("testID", 2).Write("") // Final message in Write(string) is optional.

    // Set level and tag. Since it is set to DebugLevel (and above) and/or tagHTTP,
    // log entriees with DebugLevel or higher AND also log entries containing tagHTTP will show up.
    log.Control(alog.DebugLevel, tagHTTP)

    log.Debug(tagHTTP|tagREQ).Str("requestFrom", "123.0.1.100").Int("testID", 3).Write("will show")
    log.Trace(tagDB).Int("status", 200).Str("dest", "123.0.1.100").Int("testID", 4).Write("will not show")
    log.Trace(tagHTTP|tagRES).Int("status", 200).Str("dest", "123.0.1.100").Int("testID", 5).Write("will show")

    // Output:
    // {"date":20210304,"time":173510,"level":"info","tag":[],"message":"Starting","name":"Alog","buildNo":6,"testID":1}
    // {"date":20210304,"time":173510,"level":"debug","tag":["HTTP","REQ"],"message":"will show","requestFrom":"123.0.1.100","testID":3}
    // {"date":20210304,"time":173510,"level":"trace","tag":["HTTP","RES"],"message":"will show","status":200,"dest":"123.0.1.100","testID":5}
  }
  ~~~
 


[^Top](#alog)



## Level and Tag 

Alog supports both leveled logging and tagging. Tagging allows minimize total number of logging,
hence can boost the performance, save disk space needed. Tagging and level can be carefully controlled by
control function as well as simply using level and tags alone.

`ControlFn` can be used by `*Logger.Control.Fn = myControlFunc`


[^Top](#alog)



## Extension

Extension

[^Top](#alog)



## Benchmark

Running a benchmark is very tricky. Depend on which system the benchmark is performed, the output can be vary and 
potentially misleading. Please note that this benchmark can be very differ than your own.

Things to consider:

- Alog's sole goal is making a logger that support tagging with zero memory allocation.
- Alog does not support fancy features supported by Zerolog.
- Alog's performance drops significantly when adding more variables. 
  This is a tradeoff between easy-to-use vs performance.  
- Code
  - Alog:  `al.Info(0).Str("name", "alog").Write("alog")`
  - Zerolog: `zl.Info().Str("name", "zlog").Msg("zlog")`


[^Top](#alog)



EOF
