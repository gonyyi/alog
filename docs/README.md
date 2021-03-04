# Alog v0.6

(c) 2021 Gon Y Yi. <https://gonyyi.com>.  
[MIT License](https://raw.githubusercontent.com/gonyyi/alog/master/LICENSE)

Version 0.6.0 (3/3/2021)

[![GoDoc](https://godoc.org/github.com/gonyyi/alog?status.svg)](https://godoc.org/github.com/gonyyi/alog)
[![Go Reference](https://pkg.go.dev/badge/github.com/gonyyi/alog.svg)](https://pkg.go.dev/github.com/gonyyi/alog@v0.6.0)
[![License](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://raw.githubusercontent.com/gonyyi/alog/master/LICENSE)
[![Coverage](http://gocover.io/_badge/github.com/gonyyi/alog)](http://gocover.io/github.com/gonyyi/alog)

![Alog Screen Shot 1](https://github.com/gonyyi/alog/blob/master/docs/alog_screen_1.png)

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
