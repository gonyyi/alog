# Alog v0.5

(c) 2020 Gon Y Yi. <https://gonyyi.com>.  
[MIT License](https://raw.githubusercontent.com/gonyyi/alog/master/LICENSE)

Version 0.5.0 (2/18/2021)

[![GoDoc](https://godoc.org/github.com/gonyyi/alog?status.svg)](https://godoc.org/github.com/gonyyi/alog)
[![Go Reference](https://pkg.go.dev/badge/github.com/gonyyi/alog.svg)](https://pkg.go.dev/github.com/gonyyi/alog@v0.5.0)
[![License](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://raw.githubusercontent.com/gonyyi/alog/master/LICENSE)
[![Coverage](http://gocover.io/_badge/github.com/gonyyi/alog)](http://gocover.io/github.com/gonyyi/alog)

![Alog with an ExtFormatterTextColor](http://alog.lib.gonyyi.com/images/screenshot1.png)
_(see: <https://github.com/gonyyi/alog/blob/master/benchmark/ext_formatterTextColor_test.go>)_


## Introduction

Alog was built with a very simple goal in mind:

- Tagging (besides the level)
- Minimum memory allocation

If you find any issues, please [create an issue](https://github.com/gonyyi/alog/issues/new).


## Level and Tag 

Alog supports both leveled logging and tagging. Tagging allows minimize total number of logging,
hence can boost the performance, save disk space needed. Tagging and level can be carefully controlled by
control function as well as simply using level and tags alone.


## Extension

Currently as of _v0.5.0_ Alog supports three(3) interfaces for extension: `Writer`, `Formatter`, and `Buffer`.

Alog's extensions are in the `/ext` folder. Currently, there are 3 formatter extensions written.

- `FormatterTextColor`: ANSI colored terminal output 
- `FormatterText`: Regular text based log output
- `FormatterJSON`: JSON output

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
    - Alog:    `al.Info(0,  "alog test", "name", "alog")`
	- Zerolog: `zl.Info().Str("name", "zlog").Msg("zerolog")`

| Logger         | Test# | Counts  | Time        | Memory  | Alloc       |
|:---------------|:------|:--------|:------------|:--------|:------------|
| Alog           | 1     | 4247889 | 271.2 ns/op | 0 B/op  | 0 allocs/op |
| Alog           | 2     | 4410723 | 268.3 ns/op | 0 B/op  | 0 allocs/op |
| Alog           | 3     | 4404769 | 268.2 ns/op | 0 B/op  | 0 allocs/op |
| Zerolog        | 1     | 3938845 | 301.5 ns/op | 0 B/op  | 0 allocs/op |
| Zerolog        | 2     | 3934148 | 301.6 ns/op | 0 B/op  | 0 allocs/op |
| Zerolog        | 3     | 3922243 | 301.5 ns/op | 0 B/op  | 0 allocs/op |
| Builtin Logger | 1     | 4991648 | 238.9 ns/op | 16 B/op | 1 allocs/op |
| Builtin Logger | 2     | 5002394 | 238.8 ns/op | 16 B/op | 1 allocs/op |
| Builtin Logger | 3     | 4994046 | 239.7 ns/op | 16 B/op | 1 allocs/op |

_Detail can be found in the `/benchmark` folder_ <https://github.com/gonyyi/alog/blob/master/benchmark>


[^Top](#alog)


EOF
