# Alog

(c) 2020 Gon Y Yi. <https://gonyyi.com/copyright.txt>  
Version 0.1.4 (12/29/2020)


## Table of Contents

1. [Introduction](#introduction)
2. [Changes](#changes)
    - [v0.1.4](#v014)
    - [v0.1.3](#v013)
    - [v0.1.2](#v012)
    - [v0.1.1](#v011)
    - [v0.1.0](#v010)
3. [Examples](#examples)
    - [Leveled Log](#leveled-log)
    - [Category Support](#category-support)
    - [With a Buffered Writer](#with-a-buffered-writer)
    - [NewPrint](#newprint)
    - [NewWriter](#newwriter)
4. [Note: Formatted Output](#note-formatted-output)
5. [Benchmark](#benchmark)


## Introduction

Alog is a simple dependency-free logger with a goal of zero memory allocation.
Alog supports leveled logging with optional category support.

If you find any bug/concern about performance, 
please [create an issue](https://github.com/gonyyi/alog/issues/new).

[^Top](#alog)


## Changes

### v0.1.4

New flag option `Fnewline` has been added. The default behavior is not allowing newline
within the log message. However, using this option will allow newlines in the log message.


[^Top](#alog)


### v0.1.3

NewWriter and NewPrint now takes additional string argument for prefix.

- `*Logger.NewWriter(level, Category, string) *alogw`
- `*Logger.NewPrint(level, Category, string) func(string)`

[^Top](#alog)


### v0.1.2

- Added new method `*Logger.NewWriter(level, Category) *alogw`
    - This is compatible with io.Writer interface.
    - This can be used as a log hook for libraries.

[^Top](#alog)


### v0.1.1

- Added new method `*Logger.NewPrint(level, Category) func(string)`
    - This will be used to create custom logging with defining category each time.
    - `Printf` has not been implemented due to memory allocation.

[^Top](#alog)


### v0.1.0

- **Removed buffer flags**: this can be replaced with `bufio.Writer`.
- **Removed logging level flags from constructor**: constructor `New()` was taking
    configuration bitflag for formatting as well as logging level. As those two
    aren't the same kind, logging level flag has been removed. Default logging
    level is `INFO`, and a user can reset it by using `*logger.SetLevel()` method.
    - [Issue #8](https://github.com/gonyyi/alog/issues/8) type Format should be renamed to configuration or something
- **Removd `SetExitOnFatal()`**: `*logger.SetExitOnFatal()` was used to set if
    a user want to exit when fatal level log is received. Now, `Print()` and `Printf()`
    will not exit when received fatal level log, but `Fatal()` and `Fatalf()` will
    exit with exit code 1.
- **Renamed `SetLevels()` into `SetLevelPrefix()`**: A method `SetLevels()` was
    not an intuitive name, created confusion. So, renamed it to `SetLevelPrefix()`
- **Unexport few types**: previously, alog exported type for `Level` and `Format`. 
    those two types are not unexported.
    - [Issue #8](https://github.com/gonyyi/alog/issues/8) type Format should be renamed to configuration or something
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
    l := alog.New(os.Stdout, "test ", alog.Fprefix|alog.Flevel)

    // Trace/Debug will NOT be printed
    l.Trace("hello trace")
    l.Debug("hello debug")

    // Info/Warn/Error/Fatal will be printed
    l.Info("hello info")
    l.Warn("hello warn")
    l.Error("hello error")
}
```

[^Top](#alog)


### Category Support

```go
package main

import (
    "github.com/gonyyi/alog"
    "os"
)

func main() {
    cat := alog.NewCategory()
    BACK := cat.Add()
    FRONT := cat.Add()
    USER := cat.Add()
    l := alog.New(os.Stdout, "test ", alog.Flevel|alog.Fprefix)
    l.SetLevel(alog.Ltrace)

    // Assume I want to see BACK and FRONT with a level INFO or above.
    l.SetCategory(BACK | FRONT) // only show BACK and FRONT
    l.SetLevel(alog.Linfo)      // this will override config "alog.CLevelTrace"

    f := func(c alog.Category, s string) {
        l.Printf(alog.Ltrace, c, "%s.trace", s)
        l.Printf(alog.Ldebug, c, "%s.debug", s)
        l.Printf(alog.Linfo, c, "%s.info", s)
        l.Printf(alog.Lwarn, c, "%s.warn", s)
        l.Printf(alog.Lerror, c, "%s.error", s)
    }

    f(BACK, "BACK")
    f(FRONT, "FRONT")
    f(USER, "USER")
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
	fLog, _ := os.Create("./alogtest/test.log")
	bLog := bufio.NewWriter(fLog)

	// Create an Alog with default option (MMDD, Time, Level) + UTC time.
	l := alog.New(bLog, "", alog.Fdefault|alog.FtimeUTC)

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
	l := alog.New(os.Stderr, "", alog.Fdefault)

	cat := alog.NewCategory()
	USER := cat.Add()
	DB := cat.Add()

	l.SetCategory(USER)

	UserInfo := l.NewPrint(alog.Linfo, USER, "USER: ")
	DBInfo := l.NewPrint(alog.Linfo, DB, "DB: ")

	UserInfo("test cat: user, lvl: info") // Printed
	DBInfo("test cat: DB, lvl: info")     // Not printed as category is set to USER
}
```

[^Top](#alog)


### NewWriter

`*Logger.NewWriter` takes a level and a category then creates an alogw object which is io.Writer compatible.
This can be used as a writer hook. Assume there is an API that takes io.Writer, you can preset the level and
category and just plug it in.

```go
package main

import (
	"github.com/gonyyi/alog"
	"os"
)

func main() {
	l := alog.New(os.Stdout, "nptest ", alog.Fprefix|alog.Flevel) // Default level is INFO and higher
	l.SetLevel(alog.Ldebug) // set logging level to DEBUG

	cat := alog.NewCategory()
	TEST1 := cat.Add()

	wT1D := l.NewWriter(alog.Ldebug, TEST1, "T1D: ")
	wT1I := l.NewWriter(alog.Linfo, TEST1, "T1I: ")

    // Assume API takes an io.Writer interface,
	fmt.Fprintf(wT1D, "test: %s fprintf", "T1D")
	fmt.Fprintf(wT1I, "test: %s fprintf", "T1I")
}
```

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
