// Copyright 2017 The go-vgo Project Developers. See the COPYRIGHT
// file at the top-level directory of this distribution and at
// https://github.com/go-vgo/gt/blob/master/LICENSE
//
// Licensed under the Apache License, Version 2.0 <LICENSE-APACHE or
// http://www.apache.org/licenses/LICENSE-2.0> or the MIT license
// <LICENSE-MIT or http://opensource.org/licenses/MIT>, at your
// option. This file may not be copied, modified, or distributed
// except according to those terms.

package zlog

import (
	// "errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-vgo/gt/conf"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Zlog zlog struct
type Zlog struct {
	// log.Logger
}

type logConfig struct {
	Mode    string
	Path    string
	Name    string
	MaxDays int64 `toml:"max_days"`
	// Srv  Server     `toml:"server"`
}

var (
	logger, errLogger *zap.Logger
	sugar, errSugar   *zap.SugaredLogger
	config            logConfig

	zapErr error
	// ZlogTime zlog time, zapcore.Field
	ZlogTime = zap.String("time", time.Now().Format("2006-01-02 15:04:05"))
)

// Init zap log and config
func Init(tpath string) {
	// if _, err := toml.DecodeFile(tpath, &config); err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	conf.Init(tpath, &config)
	go conf.Watch(tpath, &config)

	go deleteOldLog()

	if config.Mode == "dev" {
		InitDev()
		ZlogTime = zap.Error(nil)
	} else {
		InitLog()
		// InitErrLog()
		go InitErrLog()
	}
}

func deleteOldLog() {
	fileDir, _ := confPath()
	var maxDays int64 = 28

	if config.MaxDays != 0 {
		maxDays = config.MaxDays
	}

	filepath.Walk(fileDir, func(path string, info os.FileInfo, err error) (
		returnErr error) {
		defer func() {
			if r := recover(); r != nil {
				returnErr = fmt.Errorf("Unable to delete old log '%s', error: %+v",
					path, r)
			}
		}()

		if info.IsDir() && info.ModTime().Unix() < (time.Now().Unix()-60*60*24*maxDays) {

			if strings.HasPrefix(filepath.Base(path), filepath.Base(fileDir)) {
				// if err := os.Remove(path); err != nil {
				if err := os.RemoveAll(path); err != nil {
					returnErr = fmt.Errorf("Failed to remove %s: %v", path, err)
				}
			}
		}
		return returnErr
	})
}

// InitDev init dev mode
func InitDev() {
	// logger, _ = zap.NewProduction()
	logCfg := zap.NewDevelopmentConfig()
	logCfg.Sampling = nil
	logger, zapErr = logCfg.Build()
	if zapErr != nil {
		log.Fatal("zap.NewDevelopmentConfig error: ", zapErr)
	}

	errLogger = logger

	defer logger.Sync() // flushes buffer, if any
	sugar = logger.Sugar()
	errSugar = sugar
}

func confPath() (string, string) {
	// var lpath, name string
	var lpath, name string = "./log", "log"

	if config.Path != "" {
		lpath = config.Path
	}

	if config.Name != "" {
		name = config.Name
	}

	return lpath, name
}

// InitLog init log lumberjack
func InitLog() {
	lpath, name := confPath()
	maxDays := 28
	if config.MaxDays != 0 {
		maxDays = int(config.MaxDays)
	}

	logTime := time.Now().Format("2006-01-02")
	logPath := lpath + "/" + logTime + "/" + name + ".json"
	ws := zapcore.AddSync(&lumberjack.Logger{
		Filename:   logPath,
		MaxSize:    500, // megabytes
		MaxBackups: 3,
		MaxAge:     maxDays, // days
	})
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		ws,
		zap.InfoLevel,
	)
	// logger = zap.New(core).WithOptions(zap.AddCaller())
	logger = zap.New(core).WithOptions(zap.AddStacktrace(zap.InfoLevel))

	defer logger.Sync() // flushes buffer, if any
	sugar = logger.Sugar()
}

// InitErrLog init error log and lumberjack
func InitErrLog() {
	// lumberjack.Logger is already safe for concurrent use, so we don't need to
	// lock it.
	lpath, name := confPath()
	maxDays := 28
	if config.MaxDays != 0 {
		maxDays = int(config.MaxDays)
	}

	logTime := time.Now().Format("2006-01-02")
	logPath := lpath + "/" + logTime + "/" + name + "_err.json"
	ws := zapcore.AddSync(&lumberjack.Logger{
		Filename:   logPath,
		MaxSize:    500, // megabytes
		MaxBackups: 3,
		MaxAge:     maxDays, // days
	})

	highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		ws,
		// zap.ErrorLevel,
		highPriority,
	)

	errLogger = zap.New(core).WithOptions(zap.AddStacktrace(zap.ErrorLevel))
	defer logger.Sync() // flushes buffer, if any
	errSugar = errLogger.Sugar()
}

// Err zap.Error
func Err(err error) zapcore.Field {
	return zap.Error(err)
}

// Str zap.String
func Str(key, val string) zapcore.Field {
	return zap.String(key, val)
}

// Int zap.Int
func Int(key string, val int) zapcore.Field {
	return zap.Int(key, val)
}

// Any zap.Any
func Any(key string, val interface{}) zapcore.Field {
	return zap.Any(key, val)
}

// Bool zap.Bool
func Bool(key string, val bool) zapcore.Field {
	return zap.Bool(key, val)
}

// Print fmt.Sprintf
func Print(args ...interface{}) string {
	if len(args) == 0 {
		return ""
	}
	return fmt.Sprintf("%v ", args[0])
}

// Printf fmt.Sprintf
func Printf(args ...interface{}) string {
	if len(args) < 5 {
		return ""
	}

	return fmt.Sprintf(
		"method: %v, statusCode: %v, req: %s, ip: %s, time: %fs",
		args[0],
		args[1],
		args[2],
		args[3],
		args[4])
}

func (z *Zlog) Error(msg string, err error) {
	errLogger.Error(msg,
		ZlogTime,
		zap.Error(err),
	)
}

// LogInfo info log
func LogInfo(msg string, info ...string) {
	var logInfo string
	if len(info) > 0 {
		logInfo = info[0]
	}

	errLogger.Info(msg,
		ZlogTime,
		zap.String("info", logInfo),
	)
}

// Error error log
func Error(msg string, err ...error) {
	var logErr error
	if len(err) > 0 {
		logErr = err[0]
	}
	errLogger.Error(msg,
		ZlogTime,
		zap.Error(logErr),
	)
}

// Errorm more
func Errorm(msg string, fields ...zapcore.Field) {
	errLogger.Error(msg,
		fields...,
	)
}

// SugarErrorm more
func SugarErrorm(msg string, fields ...zapcore.Field) {
	errSugar.Error(msg,
		fields,
	)
}

// Fatal fatal log
func Fatal(msg string, err ...error) {
	var logErr error
	if len(err) > 0 {
		logErr = err[0]
	}
	errLogger.Fatal(msg,
		ZlogTime,
		zap.Error(logErr),
	)
}

// Panic panic log
func Panic(msg string, err ...error) {
	var logErr error
	if len(err) > 0 {
		logErr = err[0]
	}
	errLogger.Panic(msg,
		ZlogTime,
		zap.Error(logErr),
	)
}

// LogsError sugar error log
func LogsError(msg string, err error) {
	errSugar.Error(msg,
		ZlogTime,
		zap.Error(err),
	)
}

// SugarError sugar error log
func SugarError(msg string, err error) {
	errSugar.Error(msg,
		ZlogTime,
		zap.Error(err),
	)
}

// SugarFatal sugar fatal log
func SugarFatal(msg string, err error) {
	errSugar.Fatal(msg,
		ZlogTime,
		zap.Error(err),
	)
}

// SugarPanic sugar panic log
func SugarPanic(msg string, err error) {
	errSugar.Panic(msg,
		ZlogTime,
		zap.Error(err),
	)
}

// Info info log
func Info(msg string, info ...string) {
	var logInfo string
	if len(info) > 0 {
		logInfo = info[0]
	}
	logger.Info(msg,
		ZlogTime,
		zap.String("info", logInfo),
	// fields,
	)
}

// Infom more
func Infom(msg string, fields ...zapcore.Field) {
	logger.Info(msg, fields...)
}

// SugarInfom more
func SugarInfom(msg string, fields ...zapcore.Field) {
	sugar.Info(msg, fields)
}

// Warn warn log
func Warn(msg string, warn ...string) {
	var logWarn string
	if len(warn) > 0 {
		logWarn = warn[0]
	}
	logger.Warn(msg,
		ZlogTime,
		zap.String("warn", logWarn),
	)
}

// Debug debug log
func Debug(msg string, debug ...string) {
	var logDebug string
	if len(debug) > 0 {
		logDebug = debug[0]
	}
	logger.Debug(msg,
		ZlogTime,
		zap.String("debug", logDebug),
	)
}

// Infoff info log
func Infoff(msg string, fields ...zapcore.Field) {
	logger.Info(msg,
		ZlogTime,
		fields[0],
	)
}

// LogError error log
func LogError(msg string, err error) {
	logger.Error(msg,
		ZlogTime,
		zap.Error(err),
	)
}

// LogPanic panic log
func LogPanic(msg string, err error) {
	logger.Panic(msg,
		ZlogTime,
		zap.Error(err),
	)
}

// LogFatal fatal log
func LogFatal(msg string, err error) {
	logger.Fatal(msg,
		ZlogTime,
		zap.Error(err),
	)
}

// Infof infof log
func Infof(msg, info string) {
	sugar.Infof(msg,
		ZlogTime,
		zap.String("info", info),
	)
}

// InfoW infow log
func InfoW(msg, info string) {
	sugar.Infow(msg,
		ZlogTime,
		"info", info,
	)
}

// Errorf errorf log
func Errorf(msg string, err error) {
	sugar.Errorf(msg,
		ZlogTime,
		zap.Error(err),
	)
}

// Warnf warnf log
func Warnf(msg, warn string) {
	sugar.Warnf(msg,
		ZlogTime,
		zap.String("warn", warn),
	)
}
