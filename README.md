# Alog

(c) 2020 Gon Y Yi. <https://gonyyi.com>  
[MIT License](https://raw.githubusercontent.com/gonyyi/alog/master/LICENSE)

Version 0.3.0 (1/3/2020)

[![GoDoc](https://godoc.org/github.com/gonyyi/alog?status.svg)](https://godoc.org/github.com/gonyyi/alog)
[![Go Reference](https://pkg.go.dev/badge/github.com/gonyyi/alog.svg)](https://pkg.go.dev/github.com/gonyyi/alog@v0.2.0)
[![License](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://raw.githubusercontent.com/gonyyi/alog/master/LICENSE)
[![Coverage](http://gocover.io/_badge/github.com/gonyyi/alog)](http://gocover.io/github.com/gonyyi/alog)


## Table of Contents

1. [Introduction](#introduction)
2. [Examples](#examples)
    - [Leveled Log](#leveled-log)
    - [Tag Support](#tag-support)
    - [With a Buffered Writer](#with-a-buffered-writer)
    - [NewPrint](#newprint)
    - [NewWriter](#newwriter)
    - [Using a Custom/Predefined Configuration](#using-a-custompredefined-configuration)
3. [Changes](#changes)
4. [Note: Formatted Output](#note-formatted-output)
5. [Benchmark](#benchmark)


## Introduction

Alog is a simple dependency-free logger with a goal of zero memory allocation.
Alog supports leveled logging and tagging.

If you find any bug/concern about performance, 
please [create an issue](https://github.com/gonyyi/alog/issues/new).

[^Top](#alog)


## Examples

More examples are available in `alog_example_test.go` or from the Go Doc.

```sh
$ godoc -http=:8080
```

Then, visit <http://localhost:8080/pkg/github.com/gonyyi/alog/>

[^Top](#alog)


### Leveled Log

```go
package main

import (
	"github.com/gonyyi/alog"
	"os"
)

func main() {
	// Without FilterLevel defined, the default value for level is Info.
	l := alog.New(os.Stdout).SetFlag(alog.Fprefix|alog.Flevel).FilterLevel(alog.Ldebug)

	// Trace will NOT be printed
	l.Trace("hello trace")

	// Debug/Info/Warn/Error/Fatal will be printed
	l.Debug("hello debug")
	l.Info("hello info")
	l.Warn("hello warn")
	l.Error("hello error")
}
```

[^Top](#alog)


### Tag Support

```go
package main

import (
	"github.com/gonyyi/alog"
	"os"
)

func main() {
	var BACK, FRONT, USER alog.Tag

	// Assume I want to see BACK and FRONT with a level DEBUG or above.
	l := alog.New(os.Stdout).SetPrefix("test ").SetFlag(alog.Fprefix | alog.Flevel).
		FilterLevel(alog.Ldebug).SetTags(&BACK, &FRONT, &USER).FilterTag(BACK | FRONT)

	f := func(c alog.Tag, s string) {
		l.Printf(alog.Ltrace, c, "%s.trace", s)
		l.Printf(alog.Ldebug, c, "%s.debug", s)
		l.Printf(alog.Linfo, c, "%s.info", s)
		l.Printf(alog.Lwarn, c, "%s.warn", s)
		l.Printf(alog.Lerror, c, "%s.error", s)
	}

	f(BACK, "BACK")   // prints debug-error
	f(FRONT, "FRONT") // prints debug-error
	f(USER, "USER")   // not print anything
}
```

[^Top](#alog)


### With a Buffered Writer

```go
package main

import (
	"bufio"
	"github.com/gonyyi/alog"
	"os"
)

func main() {
	// Create a file and bufio writer
	fLog, _ := os.Create("./test.log")
	bLog := bufio.NewWriter(fLog)

	// Create an Alog with default option (MMDD, Time, Level) + UTC time.
	l := alog.New(bLog).SetFlag(alog.Fdefault|alog.FtimeUTC)

	for i := 0; i < 500; i++ {
		l.Infof("Test %s %d", "info name", i)
	}

	// Flush bufio writer
	bLog.Flush()
}
```

[^Top](#alog)


### NewPrint

```go
package main

import (
	"github.com/gonyyi/alog"
	"os"
)

func main() {
	// Create an Alog with default option (MMDD, Time, Level)
	l := alog.New(os.Stderr)

	// Another way of adding tags instead of SetTags() are assign each tag with NewTag()
	USER := l.NewTag()
	DB := l.NewTag()

	l.FilterTag(USER)

	UserInfo := l.NewPrint(alog.Linfo, USER, "USER: ")
	DBInfo := l.NewPrint(alog.Linfo, DB, "DB: ")

	UserInfo("test cat: user, lvl: info") // Printed
	DBInfo("test cat: DB, lvl: info")     // Not printed as tag is set to USER
}
```

[^Top](#alog)


### NewWriter

`*Logger.NewWriter` takes a level and a tag then creates an alogw object which is io.Writer compatible.
This can be used as a writer hook. Assume there is an API that takes io.Writer, you can preset the level and
tag and just plug it in.

```go
package main

import (
	"fmt"
	"github.com/gonyyi/alog"
	"os"
)

func main() {
	l := alog.New(os.Stdout).SetPrefix("nptest ").SetFlag(alog.Fprefix|alog.Flevel).SetLevel(alog.Ldebug)

	TEST1 := l.NewTag()

	wT1D := l.NewWriter(alog.Ldebug, TEST1, "T1D: ")
	wT1I := l.NewWriter(alog.Linfo, TEST1, "T1I: ")

	// Assume API takes an io.Writer interface,
	fmt.Fprintf(wT1D, "test: %s fprintf", "T1D")
	fmt.Fprintf(wT1I, "test: %s fprintf", "T1I")
}
```


### Using a Custom/Predefined Configuration

Preconfigured configuration function can be created and used using `*Logger.Do`.

```go
package main

import (
	"github.com/gonyyi/alog"
	"os"
)

func main() {
	myFunc := func(al *alog.Logger) {
		al.SetPrefix("alogtest ").
			SetFlag(alog.Fprefix | alog.Flevel).
			FilterLevel(alog.Ltrace)
	}

	// use myFunc defined above,
	// and also use color level by using predefined alog.DoColor function.
	// `Do` takes any number of `Do` functions.
	l := alog.New(os.Stderr).Do(myFunc, alog.DoColor()) // now with color

	// Output below will print colored level as output is set to os.Stderr in this example.
	l.Trace("test trace test")
	l.Debug("test debug test")
	l.Info("test info test")
	l.Warn("test warn test")
	l.Error("test error test")
	l.Fatal("test fatal test")
}
```


[^Top](#alog)



## Changes

### v0.3.0

1. New `*alog.Do(*Logger)` method.
2. New `*alog.NewTag() Tag`, this replaced `alog.NewCategory()`.
3. New `*alog.SetTags(...*alog.Tag)` method. This will set tags for given pointers.
4. New `*alog.Output(lvl level, tag Tag, b []byte)` method. This takes byte slice and
   should be used where bytes are used.


### v0.2.0

Due to backward compability issue, the version went up from v0.1.6 to v0.2.0. This is mainly
because of the changes of constructor: `alog.New(io.Writer, string, flag) *Logger`
became `alog.New(io.Writer) *Logger`.

1. `alog.New(io.Writer) *Logger`: Constructor Most of the time, people don't set logger
   prefix, also uses basic default setting. Therefore it's bit cumbersome to require
   two (prefix, flag), often, unused parameters.

2. `SetOutput`, `SetPrefix`, `SetFlag`, `FilterLevel`, `SetLevelPrefix`, `SetCategory` are now
   returning `*Logger` pointer which means, when a logger is created, you can add a
   configuration only when it's necessary.

    - Initially set discard for output but overridden to os.Stderr

        ```go
        l := alog.New(nil).SetOutput(os.Stderr)
        ```

   - Set prefix and level

        ```go
        l = alog.New(os.Stderr).SetPrefix("TestLog: ").SetLevel(alog.Linfo)
        ```

   - Set prefix, flag together, and level separately

        ```go
        l := alog.New(os.Stderr).SetPrefix("TestLog: ").SetFlag(alog.Fdefault|alog.FtimeUTC)
        l.FilterLevel(alog.Ltrace)

        ```


### v0.1.x

- v0.1.6
    - Added a badge for a coverage
    - Any level-predefined and formatted methods such as `Tracef`, `Debugf`, ... `Fatalf` will evaluate if any additional arguments are present besides format string. If there is no additional argument, it will run without formatting to save processing time.
    - `*.Logger.IfError(error)`, `*.Logger.IfFatal(error)` has been added. These methods are taking error (or `nil`) for an argument. If it's `nil`, it will ignore, but if actual error is given, it will log the error message.
    - `*Logger.Close()` method has been added back. If an `io.Writer` that logger uses have `Close()` method, it will call the `Close()` method of the writer.
    - `Fatal`, `Fatalf`, `IfFatal` will call `*Logger.Close()` right before the `os.Exit()`
- v0.1.5
    - Added a license (MIT) <https://raw.githubusercontent.com/gonyyi/alog/master/LICENSE>
- v0.1.4
    - New flag option `Fnewline` has been added. The default behavior is not allowing newline
      within the log message. However, using this option will allow newlines in the log message.
- v0.1.3
    - NewWriter and NewPrint now takes additional string argument for prefix.
        - `*Logger.NewWriter(level, Category, string) *alogw`
        - `*Logger.NewPrint(level, Category, string) func(string)`
- v0.1.2
    - Added new method `*Logger.NewWriter(level, Category) *alogw`
        - This is compatible with io.Writer interface.
        - This can be used as a log hook for libraries.
- v0.1.1
    - Added new method `*Logger.NewPrint(level, Category) func(string)`
        - This will be used to create custom logging with defining category each time.
        - `Printf` has not been implemented due to memory allocation.
- v0.1.0
    - **Removed buffer flags**: this can be replaced with `bufio.Writer`.
    - **Removed logging level flags from constructor**: constructor `New()` was taking
      configuration bitflag for formatting as well as logging level. As those two
      aren't the same kind, logging level flag has been removed. Default logging
      level is `INFO`, and a user can reset it by using `*logger.SetLevel()` method.
    - [Issue #8](https://github.com/gonyyi/alog/issues/8) type Format should be renamed
      to configuration or something
    - **Removed `SetExitOnFatal()`**: `*logger.SetExitOnFatal()` was used to set if
      a user want to exit when fatal level log is received. Now, `Print()` and `Printf()`
      will not exit when received fatal level log, but `Fatal()` and `Fatalf()` will
      exit with exit code 1.
    - **Renamed `SetLevels()` into `SetLevelPrefix()`**: A method `SetLevels()` was
      not an intuitive name, created confusion. So, renamed it to `SetLevelPrefix()`
    - **Unexport few types**: previously, alog exported type for `Level` and `Format`.
      those two types are not unexported.
        - [Issue #8](https://github.com/gonyyi/alog/issues/8) type Format should be renamed to
          configuration or something
    - **Added Writer()**: `*logger.Writer()` will return `io.Writer` used in the logger.
    - **Documentation update**: added comments and examples for many code base and
      compatible with GoDoc.
        - [Issue #7](https://github.com/gonyyi/alog/issues/7) Add `example_` files
        - [Issue #4](https://github.com/gonyyi/alog/issues/4) Make compatible with Godoc
        - [Issue #3](https://github.com/gonyyi/alog/issues/3) Add a comparison for internal functions
    - **Use more of switch instead of if-else**: although there isn't any performance
      gain, many codes where it has multiple if-else blocks are now converted to
      switch for better code readability.

[^Top](#alog)


## Note: Formatted Output

(_applies to `Tracef()`, `Debugf()`, `Infof()`, `Warnf()`, `Errorf()`, `Fatalf()` and `Printf()`_)

`*logger.Printf`: Alog's `Printf` has been re-written from the scratch for better
memory usages, and achieved zero allocation. However, it does **NOT** support same 
as `fmt.Printf`:

- `*logger.Printf` only supports following formats:
	- `%s`: string
	- `%d`: integer
	- `%f`: float (only up to 2 decimal points)
- For float, current version does round down, and this can cause a difference
	in 2nd decimal place.

[^Top](#alog)


## Benchmark 

Alog has been focused on memory allocation rather than the speed or fancy function.
Test was done on 2018 MacBook Pro (15-inch):

- MacOS 10.15.7 Catalina
- 2.9 GHz 6-Core Intel Core i9
- 32 GB 2400 MHz DDR4
- Radeon Pro 560X 4 GB / Intel UHD Graphics 630 1536 MB

Below benchmark is as of v0.1.6. No significant changes in v0.2.0.

| Type    | Name                           | Test               | Count      | Speed       | Mem     | Alloc       |
|:--------|:-------------------------------|:-------------------|:-----------|:------------|:--------|:------------|
| Builtin | BenchmarkBuiltinLoggerBasic-12 |                    | 2969100    | 408 ns/op   | 80 B/op | 2 allocs/op |
| Builtin | BenchmarkBuiltinLoggerFmt-12   |                    | 2534346    | 477 ns/op   | 88 B/op | 3 allocs/op |
| Alog    | BenchmarkLogger_Info           | 1_eval_0_print-12  | 1000000000 | 0.420 ns/op | 0 B/op  | 0 allocs/op |
| Alog    | BenchmarkLogger_Info           | 5_eval_0_prints-12 | 651124946  | 1.89 ns/op  | 0 B/op  | 0 allocs/op |
| Alog    | BenchmarkLogger_Info           | 5_eval_3_prints-12 | 1517264    | 777 ns/op   | 0 B/op  | 0 allocs/op |
| Alog    | BenchmarkLogger_Infof          | 1_eval_0_print-12  | 1000000000 | 1.10 ns/op  | 0 B/op  | 0 allocs/op |
| Alog    | BenchmarkLogger_Infof          | 5_eval_0_prints-12 | 227060984  | 5.48 ns/op  | 0 B/op  | 0 allocs/op |
| Alog    | BenchmarkLogger_Infof          | 5_eval_3_prints-12 | 1078039    | 1112 ns/op  | 0 B/op  | 0 allocs/op |
| Alog    | BenchmarkLogger_Print          | 5_eval_1_prints-12 | 4073192    | 304 ns/op   | 0 B/op  | 0 allocs/op |
| Alog    | BenchmarkLogger_Print          | 5_eval_2_prints-12 | 2320356    | 478 ns/op   | 0 B/op  | 0 allocs/op |
| Alog    | BenchmarkLogger_Printf         | 5_eval_1_prints-12 | 2938634    | 419 ns/op   | 0 B/op  | 0 allocs/op |
| Alog    | BenchmarkLogger_Printf         | 5_eval_2_prints-12 | 1502662    | 805 ns/op   | 0 B/op  | 0 allocs/op |
| Alog    | BenchmarkLogger_NewPrint-12    |                    | 3516364    | 337 ns/op   | 0 B/op  | 0 allocs/op |
| Alog    | BenchmarkLogger_NewWriter-12   |                    | 2853715    | 409 ns/op   | 0 B/op  | 0 allocs/op |

[^Top](#alog)


EOF
