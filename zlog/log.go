package zlog

import (
	// "errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Zlog struct {
	// log.Logger
}

var (
	logger, errLogger *zap.Logger
	sugar, errSugar   *zap.SugaredLogger
	zErr              error
	zlogTime          zapcore.Field = zap.String("time", time.Now().Format("2006-01-02 15:04:05"))
)

type logConfig struct {
	Mode    string
	Path    string
	Name    string
	MaxDays int64
	// Srv  Server     `toml:"server"`
}

var config logConfig

func Init(tpath string) {
	if _, err := toml.DecodeFile(tpath, &config); err != nil {
		fmt.Println(err)
		return
	}

	go deleteOldLog()

	if config.Mode == "dev" {
		InitDev()
	} else {
		InitLog()
		InitErrLog()
	}
}

func deleteOldLog() {
	fileDir, _ := conf()
	var maxDays int64 = 28

	if config.MaxDays != 0 {
		maxDays = config.MaxDays
	}

	filepath.Walk(fileDir, func(path string, info os.FileInfo, err error) (returnErr error) {
		defer func() {
			if r := recover(); r != nil {
				returnErr = fmt.Errorf("Unable to delete old log '%s', error: %+v", path, r)
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

func InitDev() {
	logger, _ = zap.NewProduction()
	errLogger = logger

	defer logger.Sync() // flushes buffer, if any
	sugar = logger.Sugar()
	errSugar = sugar
}

func conf() (string, string) {
	// var lpath, name string
	var lpath, name string = "./log", "foo"

	if config.Path != "" {
		lpath = config.Path
	}

	if config.Name != "" {
		name = config.Name
	}

	return lpath, name
}

func InitLog() {
	lpath, name := conf()

	logTime := time.Now().Format("2006-01-02")
	logPath := lpath + "/" + logTime + "/" + name + ".json"
	ws := zapcore.AddSync(&lumberjack.Logger{
		Filename:   logPath,
		MaxSize:    500, // megabytes
		MaxBackups: 3,
		MaxAge:     28, // days
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

func InitErrLog() {
	// lumberjack.Logger is already safe for concurrent use, so we don't need to
	// lock it.
	lpath, name := conf()

	logTime := time.Now().Format("2006-01-02")
	logPath := lpath + "/" + logTime + "/" + name + "_err.json"
	ws := zapcore.AddSync(&lumberjack.Logger{
		Filename:   logPath,
		MaxSize:    500, // megabytes
		MaxBackups: 3,
		MaxAge:     28, // days
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

func Info(msg, info string) {
	logger.Info(msg,
		zlogTime,
		zap.String("info", info),
	)
}

func Error(msg string, err error) {
	logger.Error(msg,
		zlogTime,
		zap.Error(err),
	)
}

func (z *Zlog) Error(msg string, err error) {
	errLogger.Error(msg,
		zlogTime,
		zap.Error(err),
	)
}

func LogInfo(msg, info string) {
	errLogger.Info(msg,
		zlogTime,
		zap.String("info", info),
	)
}

func LogError(msg string, err error) {
	errLogger.Error(msg,
		zlogTime,
		zap.Error(err),
	)
}

func LogsError(msg string, err error) {
	errSugar.Error(msg,
		zlogTime,
		zap.Error(err),
	)
}

func Warn(msg, warn string) {
	logger.Warn(msg,
		zlogTime,
		zap.String("warn", warn),
	)
}

func Debug(msg, debug string) {
	logger.Debug(msg,
		zlogTime,
		zap.String("debug", debug),
	)
}

func Panic(msg, warn string, err error) {
	logger.Panic(msg,
		zlogTime,
		zap.String("warn", warn),
		zap.Error(err),
	)
}

func Fatal(msg, warn string, err error) {
	logger.Fatal(msg,
		zlogTime,
		zap.String("warn", warn),
		zap.Error(err),
	)
}

func InfoW(msg, info string) {
	sugar.Infow(msg,
		zlogTime,
		"info", info,
	)
}

func Infof(msg, info string) {
	sugar.Infof(msg,
		zlogTime,
		zap.String("info", info),
	)
}

func Errorf(msg string, err error) {
	sugar.Errorf(msg,
		zlogTime,
		zap.Error(err),
	)
}

func Warnf(msg, warn string) {
	sugar.Warnf(msg,
		zlogTime,
		zap.String("warn", warn),
	)
}
