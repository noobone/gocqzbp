[![GoDoc](http://godoc.org/github.com/fumiama/cron?status.png)](http://godoc.org/github.com/fumiama/cron)
[![Build Status](https://travis-ci.org/fumiama/cron.svg?branch=master)](https://travis-ci.org/fumiama/cron)

# cron

> 从 tovenja/cron fork 来的，增加了 robfig/cron 的一些未合并的 PR

## tovenja/cron 简介
从 robfig/cron Fork出来的，不过这个cron框架经过使用后发现，大量任务会耗费很多CPU在Sort全部的任务上，每次任务执行完之后就会修改next执行时间，然后数组使用快排进行排序，时间复杂度O(nlogn)，我修改为min-heap的方式，每次添加任务的时候通过堆的属性进行up和down调整，每次添加任务时间复杂度O(logn)，经过验证线上CPU使用降低4~5倍。

From robfig/cron Fork, but after using it, I found that a large number of tasks will consume a lot of CPU on all tasks of Sort. After each task executed, the next execution time will be modified, and then the array will be sorted using quick-sort, which is cost O(nlogn), I modified it to the min-heap method. When adding a task, the up and down adjustments are made through the properties of the heap. Each time the task added, the time cost is O(logn). After product env verification, the online CPU usage is reduced by 4~5 times.

## robfig/cron 简介
Cron V3 has been released!

To download the specific tagged release, run:
```bash
go get github.com/fumiama/cron
```
Import it in your program as:
```go
import "github.com/fumiama/cron"
```
It requires Go 1.11 or later due to usage of Go Modules.

Refer to the documentation here:
http://godoc.org/github.com/fumiama/cron

The rest of this document describes the the advances in v3 and a list of
breaking changes for users that wish to upgrade from an earlier version.

### Upgrading to v3 (June 2019)

cron v3 is a major upgrade to the library that addresses all outstanding bugs,
feature requests, and rough edges. It is based on a merge of master which
contains various fixes to issues found over the years and the v2 branch which
contains some backwards-incompatible features like the ability to remove cron
jobs. In addition, v3 adds support for Go Modules, cleans up rough edges like
the timezone support, and fixes a number of bugs.

New features:

- Support for Go modules. Callers must now import this library as
  `github.com/fumiama/cron`, instead of `gopkg.in/...`

- Fixed bugs:
  - 0f01e6b parser: fix combining of Dow and Dom (#70)
  - dbf3220 adjust times when rolling the clock forward to handle non-existent midnight (#157)
  - eeecf15 spec_test.go: ensure an error is returned on 0 increment (#144)
  - 70971dc cron.Entries(): update request for snapshot to include a reply channel (#97)
  - 1cba5e6 cron: fix: removing a job causes the next scheduled job to run too late (#206)

- Standard cron spec parsing by default (first field is "minute"), with an easy
  way to opt into the seconds field (quartz-compatible). Although, note that the
  year field (optional in Quartz) is not supported.

- Extensible, key/value logging via an interface that complies with
  the https://github.com/go-logr/logr project.

- The new Chain & JobWrapper types allow you to install "interceptors" to add
  cross-cutting behavior like the following:
  - Recover any panics from jobs
  - Delay a job's execution if the previous run hasn't completed yet
  - Skip a job's execution if the previous run hasn't completed yet
  - Log each job's invocations
  - Notification when jobs are completed

It is backwards incompatible with both v1 and v2. These updates are required:

- The v1 branch accepted an optional seconds field at the beginning of the cron
  spec. This is non-standard and has led to a lot of confusion. The new default
  parser conforms to the standard as described by [the Cron wikipedia page].

  UPDATING: To retain the old behavior, construct your Cron with a custom
  parser:
```go
// Seconds field, required
cron.New(cron.WithSeconds())

// Seconds field, optional
cron.New(cron.WithParser(cron.NewParser(
	cron.SecondOptional | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor,
)))
```
- The Cron type now accepts functional options on construction rather than the
  previous ad-hoc behavior modification mechanisms (setting a field, calling a setter).

  UPDATING: Code that sets Cron.ErrorLogger or calls Cron.SetLocation must be
  updated to provide those values on construction.

- CRON_TZ is now the recommended way to specify the timezone of a single
  schedule, which is sanctioned by the specification. The legacy "TZ=" prefix
  will continue to be supported since it is unambiguous and easy to do so.

  UPDATING: No update is required.

- By default, cron will no longer recover panics in jobs that it runs.
  Recovering can be surprising (see issue #192) and seems to be at odds with
  typical behavior of libraries. Relatedly, the `cron.WithPanicLogger` option
  has been removed to accommodate the more general JobWrapper type.

  UPDATING: To opt into panic recovery and configure the panic logger:
```go
cron.New(cron.WithChain(
  cron.Recover(logger),  // or use cron.DefaultLogger
))
```
- In adding support for https://github.com/go-logr/logr, `cron.WithVerboseLogger` was
  removed, since it is duplicative with the leveled logging.

  UPDATING: Callers should use `WithLogger` and specify a logger that does not
  discard `Info` logs. For convenience, one is provided that wraps `*log.Logger`:
```go
cron.New(
  cron.WithLogger(cron.VerbosePrintfLogger(logger)))
```

#### Background - Cron spec format

There are two cron spec formats in common usage:

- The "standard" cron format, described on [the Cron wikipedia page] and used by
  the cron Linux system utility.

- The cron format used by [the Quartz Scheduler], commonly used for scheduled
  jobs in Java software

[the Cron wikipedia page]: https://en.wikipedia.org/wiki/Cron
[the Quartz Scheduler]: http://www.quartz-scheduler.org/documentation/quartz-2.3.0/tutorials/tutorial-lesson-06.html

The original version of this package included an optional "seconds" field, which
made it incompatible with both of these formats. Now, the "standard" format is
the default format accepted, and the Quartz format is opt-in.
