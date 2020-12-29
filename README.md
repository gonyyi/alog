# Alog

(c) 2020 Gon Y Yi. <https://gonyyi.com/copyright.txt>  
Version 0.1.1  
Last update: 12/29/2020


## Introduction

Alog is a simple dependency-free logger with a goal of zero memory allocation.
Alog supports leveled logging with optional category support.

If you find any bug/conceren about performance, 
please [create an issue](https://github.com/gonyyi/alog/issues/new).


## Changes


### v0.1.1

- Added new method `*Logger.NewPrint(level, Category) func(string)`
    - This will be used to create custom logging with defining category each time.
    - `Printf` has not been implemented due to memory allocation.


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


## Examples

More examples are available in `alog_example_test.go` or from the Go Doc.

```sh
$ godoc -http=:8080
```

Then, visit <http://localhost:8080/pkg/github.com/gonyyi/alog/>


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


### Using Buffered Writer (bufio)

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


## Benchmark 

Alog has been focused on memory allocation rather than the speed or fancy function.
Test was done on 2018 MacBook Pro (15-inch):

- MacOS 10.15.7 Catalina
- 2.9 GHz 6-Core Intel Core i9
- 32 GB 2400 MHz DDR4
- Radeon Pro 560X 4 GB / Intel UHD Graphics 630 1536 MB


### Baseline Go-Builtin Logger

| Name                           | Total   | TimeTook  | MemoryUsed | Allocation  |
|:-------------------------------|:--------|:----------|:-----------|:------------|
| BenchmarkBuiltinLoggerBasic-12 | 2883834 | 410 ns/op | 80 B/op    | 2 allocs/op |
| BenchmarkBuiltinLoggerFmt-12   | 2396258 | 508 ns/op | 88 B/op    | 3 allocs/op |


### Alog Logger

| Name                           | Total      | TimeTook   | MemoryUsed | Allocation  |
|:-------------------------------|:-----------|:-----------|:-----------|:------------|
| BenchmarkAlogPrintf-12         | 3314240    | 361 ns/op  | 0 B/op     | 0 allocs/op |
| BenchmarkAlogPrint-12          | 4725747    | 250 ns/op  | 0 B/op     | 0 allocs/op |
| BenchmarkAlogInfof-12          | 3195420    | 365 ns/op  | 0 B/op     | 0 allocs/op |
| BenchmarkAlogInfo-12           | 4710241    | 258 ns/op  | 0 B/op     | 0 allocs/op |
| Benchmark_ALog_Level_3-12      | 1482876    | 808 ns/op  | 0 B/op     | 0 allocs/op |
| Benchmark_ALog_Levelf_3-12     | 974910     | 1191 ns/op | 0 B/op     | 0 allocs/op |
| Benchmark_ALog_Levelf5_0-12    | 211717482  | 5.68 ns/op | 0 B/op     | 0 allocs/op |
| Benchmark_ALog_Levelf1_0-12    | 1000000000 | 1.19 ns/op | 0 B/op     | 0 allocs/op |
| Benchmark_ALogPrint_Cat5_1-12  | 4502719    | 263 ns/op  | 0 B/op     | 0 allocs/op |
| Benchmark_ALogPrint_Cat5_2-12  | 2311810    | 518 ns/op  | 0 B/op     | 0 allocs/op |
| Benchmark_ALogPrintf_Cat5_1-12 | 2869698    | 403 ns/op  | 0 B/op     | 0 allocs/op |
| Benchmark_ALogPrintf_Cat5_2-12 | 1564646    | 774 ns/op  | 0 B/op     | 0 allocs/op |
