# Alog

(c) 2020-2021 Gon Y Yi. <https://gonyyi.com>.  
[MIT License](https://raw.githubusercontent.com/gonyyi/alog/master/LICENSE)

Version 0.7


[![codecov](https://codecov.io/gh/gonyyi/alog/branch/master/graph/badge.svg?token=Y9RT0VRUQZ)](https://codecov.io/gh/gonyyi/alog)
[![Go Reference](https://pkg.go.dev/badge/github.com/gonyyi/alog.svg)](https://pkg.go.dev/github.com/gonyyi/alog@v0.7.4)
[![License](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://raw.githubusercontent.com/gonyyi/alog/master/LICENSE)

![Alog Screen Shot 1](https://github.com/gonyyi/alog/blob/master/docs/alog_screen_1.png)


## Intro

Alog was built with a very simple goal in mind:

- Support Tagging (and also level)
- No memory allocation (or minimum allocation)
- Customizable (see alog/ext for example)

If you find any issues,
please [create an issue](https://github.com/gonyyi/alog/issues/new).

[^Top](#alog)


## Example Usage

### Create logger

  ~~~go
  al := alog.New(os.Stderr)
  ~~~


### Hello World

  ~~~go
  // always has to end with "Write()", otherwise, it won't log, and will cause memory allocations.
  // output example: 
  //   {"date":20210308,"time":203337,"level":"info","tag":[],"message":"Hello World"}
  al.Info().Write("Hello World")
  ~~~


### New Tag

  ~~~go
  // Create a tag "Disk"
  // Create a tag "DB"
  tagDisk := al.NewTag("Disk") 
  tagDB := al.NewTag("DB")     

  al.Info(tagDisk).Str("action", "reading disk").Write()
  al.Info(tagDB).Str("id", "myID").Str("pwd", "myPasswd").Write("Login") // Anything in `Write(string)` will be printed as `message`.
  al.Info(tagDisk,tagDB).Int("status", 200).Write("Login")
  al.Info(tagDisk|tagDB).Int("status", 200).Write("Logout") // tags can be used as `tagDisk|tagDB` or `tagDisk,tagDB` format  
  
  // Output:
  // {"date":20210308,"time":203835,"level":"info","tag":["Disk"],"action":"reading disk"}
  // {"date":20210308,"time":203835,"level":"info","tag":["DB"],"message":"Login","id":"myID","pwd":"myPasswd"}
  // {"date":20210308,"time":203835,"level":"info","tag":["Disk","DB"],"message":"Login","status":200}
  // {"date":20210308,"time":203835,"level":"info","tag":["Disk","DB"],"message":"Logout","status":200}
  ~~~


### Change Format

![Alog Screen Shot 2](https://github.com/gonyyi/alog/blob/master/docs/alog_screen_text_color.png)

  ~~~go
  package main
   
  import (
	  "github.com/gonyyi/alog"
	  "github.com/gonyyi/alog/ext"
	  "os"
  )

  func main() {
    // To color text format. 
    al := alog.New(os.Stderr).Ext(ext.LogFmt.TextColor())
    tagDisk := al.NewTag("Disk")
    tagDB := al.NewTag("DB")
    
    al.Info(tagDisk).Str("action", "reading disk").Write()
    al.Warn(tagDB).Str("id", "myID").Str("pwd", "myPasswd").Write("Login")
    al.Error(tagDisk, tagDB).Int("status", 200).Write("Login")
    al.Fatal(tagDisk|tagDB).Int("status", 200).Write("Logout")
  }
  ~~~



### More example

  ~~~go
  package main

  import (
    "github.com/gonyyi/alog"
    "github.com/gonyyi/alog/ext"
    "os"
  )

  func main() {
    // When creating, if nil is given for io.Writer, 
    // io.Discard will be used. 
    al := alog.New(os.Stderr) 
    
    // Level + Optional Tag + Write(Optional Message)  
    al.Info().Write("info level log message")
    // Output:
    // {"date":20210304,"time":175516,"level":"info","tag":[],"message":"info level log message"}

    // Creating tags
    // Note: if same name of tags are created, it will return same tag value.
    //   eg. DB1 := al.NewTag("DB")
    //       DB2 := al.NewTag("DB") // DB1 == DB2 (true)
    IO := al.NewTag("IO")
    DB := al.NewTag("DB")
    SYS := al.NewTag("SYS")
    NET := al.NewTag("NET")
    TEST := al.NewTag("TEST")

    // Write warning log with id = 1.
    // Note: Args for Write(...string) can omitted. 
    //   eg. al.Warn().Int("id", 1).Write()
    // Note: If more than one string given, Alog will take the first one.
    // Output: {"date":20210304,"time":180012,"level":"warn","tag":[],"id":1}
    al.Warn().Int("id", 1).Write()
    

    // Create a warning log with IO tag, and value of id = 2, message = "this will print IO in the tag"
    // Output: {"date":20210304,"time":180012,"level":"warn","tag":["IO"],"message":"this will print IO in the tag","id":2}
    al.Warn(IO).Int("id", 2).Write("this will print IO in the tag") 
    

    // Create a warning log with IO/DB/SYS/NET tag, id = 3, message = "this will print IO/DB/SYS/NET to tag"
    // Tags can be listed with with pipe as well instead of comma: 
    //   `al.Warn(IO|DB|SYS|NET).Int("id", 3).Write("this will print IO/DB/SYS/NET to tag")`
    // Output: {"date":20210304,"time":180012,"level":"warn","tag":["IO","DB","SYS","NET"],"message":"this will print IO/DB/SYS/NET to tag","id":3}
    al.Warn(IO,DB,SYS,NET).Int("id", 3).Write("this will print IO/DB/SYS/NET to tag") 
    
    
    // Change logging level to Fatal. (highest logging level)
    al.Control.Level = alog.FatalLevel 


    // No output as below Fatal level
    al.Info(IO|SYS).Int("id", 4).Write("this will not print") 
    al.Info(IO).Int("id", 5).Write("this will not print")     
    al.Error(DB).Int("id", 6).Write("this will not print")


    // This will print because it's a Fatal level log entry
    // Output: {"date":20210304,"time":180150,"level":"fatal","tag":["NET"],"message":"this will print","id":7}
    al.Fatal(NET).Int("id", 7).Write("this will print") 
    

    // By adding Tags to control, any log entries with Fatal level or above (as set above),
    // OR any log entries that contains IO tag will show.
    al.Control.Tags = IO 
    

    // Tag contains tag IO; will be printed
    // Output: {"date":20210304,"time":180942,"level":"info","tag":["IO","SYS"],"id":4,"attempt":2}
    al.Info(IO|SYS). 
      Int("id", 4).
      Int("attempt", 2).
      Write()

    // Tag is IO; will be printed  
    // Output: {"date":20210304,"time":180942,"level":"info","tag":["IO"],"id":5,"attempt":2}
    al.Info(IO). 
      Int("id", 5).
      Int("attempt", 2).
      Write()

    // Tag IO isn't used. And level is below Fatal. Will NOT be printed.
    al.Error(DB). 
      Int("id", 6).
      Int("attempt", 2).
      Write()

    // Tag is not matching, but Level is, so will be printed.
    // Output: {"date":20210304,"time":180942,"level":"fatal","tag":["NET"],"id":7,"attempt":2}
    al.Fatal(NET). 
      Int("id", 7).
      Int("attempt", 2).
      Write()


    // EXTENSIONS: Formatter Extension <github.com/gonyyi/alog/ext>
    // revert logging level to INFO level (from fatal)
    al.Control.Level = alog.InfoLevel 


    // Use color formatter extension
    // This will output the log with ANSI colored text format.
    // Output (in Color): 
    //   2021-0304 18:13:38  INF  [TEST] testType="colorText"
    al = al.Ext(ext.LogFmt.TextColor()) 
    al.Info(TEST).Str("testType", "colorText").Write() 
    

    // Use text formatter extension
    // This will output the log with ANSI colored text format.
    // Output: 
    //   2021-0304 18:14:24 INF [TEST] testType="normalText"
    al = al.Ext(ext.LogFmt.Text())
    al.Info(TEST).Str("testType", "normalText").Write() 
    

    // Use default formatter (JSON)
    // This will output the log with default JSON format.
    // Output: 
    //   {"date":20210304,"time":181615,"level":"info","tag":["TEST"],"testType":"backToJSON"}
    al = al.Ext(ext.LogFmt.None())
    al.Info(TEST).Str("testType", "backToJSON").Write() 
    

    // EXTENSION: Custom Log Entry
    // A user can create a custom log entry function.
    // Example below takes a name and set name; also add two additional values
    // of string "testStr" with "myStr", and integer "testInt" with 123.
    myEntry := func(s string) alog.EntryFn {
        return func(entry *alog.Entry) *alog.Entry {
            return entry.Str("name", s).Str("testStr", "myStr").Int("testInt", 123)
        }
    }

    // When using custom log entry function, use Ext() method from the *Entry.
    // Output: {"date":20210305,"time":81210,"level":"info","tag":[],"message":"ok","name":"GON","testStr":"myStr","testInt":123}
    al.Info().Ext(myEntry("GON")).Write("ok")
    
    
    // EXTENSION: Macro Type <github.com/gonyyi/alog/ext>
    // Currently there are 3 log modes: LogMode.Prod(), LogMode.Dev(), and LogMode.Test().
    // - Prod Mode: to a file, buffered writer, JSON
    // - Dev Mode:  to a file, writer, JSON
    // - Test Mode: stderr, TextColor
    // All three modes take same arguments, but for Test mode, the argument (filename)
    // will be ignored. The reason it is still required for Test() is for users to quickly
    // switch between different modes.
    al = al.Ext(ext.LogMode.Prod("output.log")) 
    al.Info(IO).Str("testType", "PROD").Write() 

    // `*Logger.Close() error` will close io.Writer if Close() method is available. 
    // It is not required for testing, but when saved to a file, especially buffered
    // file, it is important to call Close() method.
    al.Close() 
  }
  ~~~

[^Top](#alog)



## Example - Quick Start

`github.com/gonyyi/alog/log` _(see /log suffix)_, Alog can be used right away.

Default format is `*Logger.LEVEL(...TAG).TYPE(KEY, VALUE)...Write(...MSG)`

- __LEVEL:__ `Trace`, `Debug`, `Info`, `Warn`, `Error`, `Fatal`
- __TAG:__   defined by user. Can be used with comma or pipe.
- __TYPE:__  `Int`, `Int64`, `Str`, `Bool`, `Err`, `Float`, `Ext`
- __KEY:__   string of key
- __VALUE:__ values depend on type.
- __MSG:__   optional message. Only first will be used.



### Example 1.

  ~~~go
  // Create an info level log with "MyTag",
  // city = "Gonway", zip = "12345"
  // name = "Gon", age = 50,
  // isMarried = false, height = 5.8
  log.Info(MyTag).
  Str("city", "Gonway").Str("zip", "12345").
  Str("name", "Gon").Int("age", 50).
  Bool("isMarried", false).
  Float("height", 5.8).Write() // Make sure all log entries must end with Write()
  ~~~


### Example 2.

  ~~~go
  package main

  import (
    "github.com/gonyyi/alog"
    "github.com/gonyyi/alog/log" // quick start
  )

  func main() {
    // Alog format is
    //    *Logger.`LOGLEVEL(TAG)`.`Str/Int/Int64/Float/Bool/Err(key, value)`.Write(`OPTIONAL MSG`)
    log.Info().Str("name", "Alog").Int("buildNo", 6).Int("testID", 1).Write("Starting")

    // Tag can be created by `NewTag(name string)`
    tagDB := log.NewTag("DB")
    tagHTTP := log.NewTag("HTTP")
    tagREQ := log.NewTag("REQ")
    tagRES := log.NewTag("RES")

    // Below Debug will not be printed, because default log level is INFO or higher.
    // Final message in Write(...string) is optional.
    log.Debug(tagDB).Str("status", "started").Int("buildNo", 6).Int("testID", 2).Write()

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

`ControlFn` is `func(Level, Tag) bool`. This function directs true (log) or false (not log) with
given Level and Tag. Please, note that, when control function is set, it will supersedes
`*Logger.Control.Level` and `*Logger.Control.Tag`.

[^Top](#alog)



## Extension

Extensions are for users to customize alog. Few examples are written
in ext folder (`github.com/gonyyi/alog/ext`). Current examples are
as below:

- Custom Log Entry
  - Usage: `alog.Info().Ext(ext.EntryHTTP.ReqRx(h)).Write()`
- Custom Formatter
  - Usage
    - Color Text: `al = alog.New(nil).Ext(ext.LogFmt.TextColor())`
    - Text: `al = alog.New(nil).Ext(ext.LogFmt.Text())`
    - None: `al = alog.New(nil).Ext(ext.LogFmt.None())`
- Custom Mode
  - Usage
    - PROD: `al = alog.New(nil).Ext(ext.LogMode.Prod("mylog.log"))`
    - DEV: `al = alog.New(nil).Ext(ext.LogMode.Dev("mylog.log"))`
    - TEST: `al = alog.New(nil).Ext(ext.LogMode.Test("mylog.log"))`

[^Top](#alog)



## Limitation and Benchmark

### TL;DR:

- Alog does not support fancy features supported by Zerolog.
- Alog's performance drops significantly when string escape is needed.

Alog has been optimized for simple leveled and tag-based logging with zero memory allocation.
To get the performance, Alog does not check duplicate keys.

Alog was designed expecting keys and string values aren't needed for string escapes.
Alog does check for this, however, and if found escapes are needed, it will run `strconv.Quote()`.
This part has not been optimized and will slower the performance when happens.

Running a benchmark is very tricky. Depend on which system the benchmark is performed,
the output can be vary and potentially misleading. Please note that this benchmark can be very
differ than your own.

__Please note that, this benchmark test is done based on very limited cases and can be misleading.__


### Benchmark 1

| Type     | Zerolog     | Alog        | Diffs        |
|:---------|:------------|:------------|:-------------|
| Single   | 117.4 ns/op | 95.54 ns/op | 18.6% faster |
| Parallel | 24.21 ns/op | 19.21 ns/op | 20.6% faster |
| Check    | 1.581 ns/op | 1.357 ns/op | 14.1% faster |

- Tested on Mac Mini (M1, 8GB, 2020)
- Both Zerolog and Alog reported zero memory allocation (0 B/op, 0 allocs/op)
- Zerolog version: v1.20.0
- Alog version: v0.7.3


### Benchmark 2

| Type     | Zerolog     | Alog        | Diffs        |
|:---------|------------:|------------:|-------------:|
| Single   | 145.5 ns/op | 142.2 ns/op |  2.2% faster |
| Parallel |  31.5 ns/op |  25.5 ns/op | 19.0% faster |
| Check    |  2.10 ns/op |  1.67 ns/op | 20.4% faster |

- Tested on Intel Macbook Pro 15" (i9-8950HK, 32GB, 2018)
- Both Zerolog and Alog reported zero memory allocation (0 B/op, 0 allocs/op)
- Zerolog version: v1.20.0
- Alog version: v0.7.3


### Benchmark Code 

See `github.com/gonyyi/alog/_tmp` for benchmark used.

  ~~~go
  al.Info(0).
    Str("name", "gonal").
    Int("count", i).
    Str("block", dataComp.StrSlice[i%5]).
    Write(dataComp.Msg)

  zl.Info().
    Str("name", "gonzl").
    Int("count", i).
    Str("block", dataComp.StrSlice[i%5]).
    Msg(dataComp.Msg)
  ~~~

[^Top](#alog)


EOF
