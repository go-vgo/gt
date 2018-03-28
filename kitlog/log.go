// Copyright 2017 The go-vgo Project Developers. See the COPYRIGHT
// file at the top-level directory of this distribution and at
// https://github.com/go-vgo/gt/blob/master/LICENSE
//
// Licensed under the Apache License, Version 2.0 <LICENSE-APACHE or
// http://www.apache.org/licenses/LICENSE-2.0> or the MIT license
// <LICENSE-MIT or http://opensource.org/licenses/MIT>, at your
// option. This file may not be copied, modified, or distributed
// except according to those terms.

package kitlog

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

const (
	logNameFormat = "2006-01-02_15:04"

	// CallerNum = 1
	CallerNum = 5
)

// LogFormat log format
type LogFormat string

const (
	JsonFormat LogFormat = "json"
	FmtFormat            = "fmt"
	NoFormat             = "no"
)

func init() {
	tmpLog := log.NewJSONLogger(os.Stdout)
	tmpLog = log.With(tmpLog, "caller", log.DefaultCaller)
	// tmpLog = log.With(tmpLog, "caller", log.Caller)
	tmpLog = log.With(tmpLog, "ts", log.DefaultTimestampUTC)
	tmpLog = level.NewFilter(tmpLog, level.AllowAll())
	tmpLog = log.NewSyncLogger(tmpLog)

	lg = &KitLogger{
		Logger:   tmpLog,
		ioWriter: nil,
		sync:     true,
	}
}

// NewKitLogger new kit logger
func NewKitLogger(opt LogOption, format ...LogFormat) (*KitLogger, error) {
	var (
		ioWriter *LogWriter
		tmpLog   log.Logger
		err      error
	)

	if len(format) == 0 {
		if ioWriter, err = NewLogWriter(opt); err != nil {
			return nil, err
		}
		tmpLog = log.NewLogfmtLogger(ioWriter)
	} else {
		switch format[0] {
		case JsonFormat:
			if ioWriter, err = NewLogWriter(opt); err != nil {
				return nil, err
			}
			tmpLog = log.NewJSONLogger(ioWriter)
		case FmtFormat:
			if ioWriter, err = NewLogWriter(opt); err != nil {
				return nil, err
			}
			tmpLog = log.NewLogfmtLogger(ioWriter)
		case NoFormat:
			// ioWriter = &LogWriter{File: nil}
			tmpLog = log.NewNopLogger()
		default:
			panic("log format type is invalid")
		}
	}
	tmpLog = log.With(tmpLog, "caller", log.DefaultCaller)
	tmpLog = log.With(tmpLog, "ts", log.DefaultTimestampUTC)

	var kitOpt level.Option
	switch strings.ToLower(opt.LogLevel) {
	case "info":
		kitOpt = level.AllowInfo()
	case "debug":
		kitOpt = level.AllowDebug()
	case "warn":
		kitOpt = level.AllowWarn()
	case "error", "crit":
		kitOpt = level.AllowError()
	default:
		panic(fmt.Sprintf("logLevel(%s) no in [info|debug|warn|error] ", opt.LogLevel))
	}
	tmpLog = level.NewFilter(tmpLog, kitOpt)

	if opt.Sync {
		tmpLog = log.NewSyncLogger(tmpLog)
	}
	return &KitLogger{Logger: tmpLog}, nil
}

// GlobalLog global log
func GlobalLog() *KitLogger {
	return lg
}

// SetGlobalLog is not thread saftly
func SetGlobalLog(opt LogOption, format ...LogFormat) {
	Close()
	var (
		tmpLog *KitLogger
		err    error
	)
	if len(format) > 0 {
		tmpLog, err = NewKitLogger(opt, format[0])
	} else {
		tmpLog, err = NewKitLogger(opt)
	}
	if err != nil {
		panic(err)
	}
	lg = tmpLog
}

// SetGlobalLogWithLog set global logger with args logger
func SetGlobalLogWithLog(logger log.Logger, levelConf ...level.Option) {
	ioWriter := lg.ioWriter
	defer func() {
		if ioWriter != nil {
			ioWriter.Close()
		}
	}()
	// new logger has a new ioWriter
	lg.Logger = log.With(logger, "caller", log.Caller(CallerNum))
	if len(levelConf) > 0 {
		lg.Logger = level.NewFilter(lg.Logger, levelConf[0])
	} else {
		lg.Logger = level.NewFilter(lg.Logger, levelConf[0])
	}
}

var lg *KitLogger

// KitLogger kit logger
type KitLogger struct {
	log.Logger
	// *levels.Levels
	ioWriter *LogWriter
	sync     bool
}

// Close close the kit logger
func (gklog *KitLogger) Close() error {
	if gklog.ioWriter != nil {
		return gklog.ioWriter.Close()
	}
	return nil
}

// Debug debug
func Debug(args ...interface{}) {
	tmpLog := log.With(lg.Logger, "caller", log.Caller(CallerNum), "level", level.DebugValue())
	logPrint(tmpLog, args)
}

// Debugf debugf log
func Debugf(args ...interface{}) {
	tmpLog := log.With(lg.Logger, "caller", log.Caller(CallerNum), "level", level.DebugValue())
	logPrintf(tmpLog, args)
}

// Info info log
func Info(args ...interface{}) {
	tmpLog := log.With(lg.Logger, "caller", log.Caller(CallerNum), "level", level.InfoValue())
	logPrint(tmpLog, args)
}

func Infof(args ...interface{}) {
	tmpLog := log.With(lg.Logger, "caller", log.Caller(CallerNum), "level", level.InfoValue())
	logPrintf(tmpLog, args)
}

func Warn(args ...interface{}) {
	tmpLog := log.With(lg.Logger, "caller", log.Caller(CallerNum), "level", level.WarnValue())
	logPrint(tmpLog, args)
}

func Warnf(args ...interface{}) {
	tmpLog := log.With(lg.Logger, "caller", log.Caller(CallerNum), "level", level.WarnValue())
	logPrintf(tmpLog, args)
}

func Error(args ...interface{}) {
	tmpLog := log.With(lg.Logger, "caller", log.Caller(CallerNum), "level", level.ErrorValue())
	logPrint(tmpLog, args)
}

func Errorf(args ...interface{}) {
	tmpLog := log.With(lg.Logger, "caller", log.Caller(CallerNum), "level", level.ErrorValue())
	logPrintf(tmpLog, args)
}

func Crit(args ...interface{}) {
	tmpLog := log.With(lg.Logger, "caller", log.Caller(CallerNum), "level", level.ErrorValue())
	logPrint(tmpLog, args)
	os.Exit(1)
}

func Critf(args ...interface{}) {
	tmpLog := log.With(lg.Logger, "caller", log.Caller(CallerNum), "level", level.ErrorValue())
	logPrint(tmpLog, args)
	os.Exit(1)
}

func Log(args ...interface{}) {
	tmpLog := log.With(lg.Logger, "caller", log.Caller(CallerNum))
	logPrint(tmpLog, args)
}

func Close() error {
	if lg.ioWriter != nil {
		return lg.ioWriter.Close()
	}
	return nil
}

func WrapLogLevel(levelsSets []string) []level.Option {
	var kitOpts []level.Option
	for _, logLevel := range levelsSets {
		switch logLevel {
		case "info":
			kitOpts = append(kitOpts, level.AllowInfo())
		case "debug":
			kitOpts = append(kitOpts, level.AllowDebug())
		case "warn":
			kitOpts = append(kitOpts, level.AllowWarn())
		case "error":
			kitOpts = append(kitOpts, level.AllowError())
		}
	}
	return kitOpts
}

// LogOption log option
type LogOption struct {
	// unit in minutes
	SegmentationThreshold int    `toml:"threshold"`
	LogDir                string `toml:"log_dir"`
	LogName               string `toml:"log_name"`
	LogLevel              string `toml:"log_level"`
	Sync                  bool   `toml:"sync"`
}

type LogWriter struct {
	oldTime               time.Time
	segmentationThreshold float64
	logDir                string
	logName               string
	*os.File
}

func NewLogWriter(opt LogOption) (*LogWriter, error) {
	logWriter := &LogWriter{
		oldTime:               time.Now(),
		segmentationThreshold: float64(opt.SegmentationThreshold),
		logDir:                opt.LogDir,
		logName:               opt.LogName,
	}

	fp, err := os.OpenFile(fmt.Sprintf("%s/%s.log", opt.LogDir, opt.LogName),
		os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
	if err != nil {
		return nil, err
	}
	logWriter.File = fp
	return logWriter, nil
}

// TODO
// use bufio buffer
func (lw *LogWriter) Write(p []byte) (n int, err error) {
	if lw.File == nil {
		return
	}
	if time.Since(lw.oldTime).Minutes() > lw.segmentationThreshold {
		if err = lw.renameLogFile(); err != nil {
			return -1, err
		}

		lw.File, err = os.OpenFile(fmt.Sprintf("%s/%s.log", lw.logDir, lw.logName),
			os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm)
		if err != nil {
			return -1, err
		}

	}
	return lw.File.Write(p)
}

func (lw *LogWriter) Close() error {
	if lw.File != nil {
		return lw.File.Close()
	}
	return nil
}

func (lw *LogWriter) renameLogFile() (err error) {
	var (
		stat                     os.FileInfo
		srcFileName, dstFileName string
	)

	if lw.File == nil {
		return
	}

	if stat, err = lw.File.Stat(); err != nil {
		return err
	}
	srcFileName = fmt.Sprintf("%s/%s", lw.logDir, stat.Name())
	if err = lw.File.Close(); err != nil {
		return err
	}
	dstFileName = fmt.Sprintf("%s/%s_%s.log", lw.logDir, lw.logName,
		lw.oldTime.Format(logNameFormat))
	fmt.Println(dstFileName, srcFileName)
	os.Rename(srcFileName, dstFileName)
	lw.oldTime = time.Now()
	return nil
}

func logPrint(logger log.Logger, args []interface{}) {
	if args == nil || len(args) == 0 {
		logger.Log()
		return
	}
	if len(args) == 1 {
		logger.Log("msg", fmt.Sprintf("%v", args[0]))
		return
	}
	for i := 0; i < len(args); i++ {
		args[i] = fmt.Sprintf("%v", args[i])
	}
	logger.Log(args...)
}

func logPrintf(logger log.Logger, args []interface{}) {
	var (
		logFormat, msgContent string
	)
	if args == nil || len(args) == 0 {
		logger.Log("msg")
		return
	}
	if len(args) == 1 {
		logger.Log("msg", fmt.Sprintf("%s", args[0]))
		return
	}

	logFormat = fmt.Sprintf("%v", args[0])
	msgContent = fmt.Sprintf(logFormat, args[1:]...)
	logger.Log("msg", msgContent)
}
