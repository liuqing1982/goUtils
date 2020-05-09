// 描述:
//    全局配置日志

// 全局配置:
//    logger := Logger{log_path, level, RotationTime, MaxAge, Console}
//    logger.LogInit()

// 各文件使用日志
// import (
//      log "github.com/sirupsen/logrus"
//     )
// 	log.Info("this is info")
//	log.Debug("this is debug")
//	log.Error("this is error")

package logger

import (
	"bufio"
	"fmt"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	log "github.com/sirupsen/logrus"
	"os"
	"path"
	"runtime"
	"strings"
	"time"
)

type Logger struct {
	Path         string
	Level        string
	RotationTime int
	MaxAge       int
	Console      bool
}

func (logger *Logger) LogInit() {
	setLogLevel(logger.Level)
	setLogPath(logger.Path, logger.RotationTime, logger.MaxAge)
	setLogConsole(logger.Console)
}

type MyFormatter struct {
}

func caller(skip int) (string, bool) {
	// 日志打印文件名、 行号， 方法名称
	pc, file, line, ok := runtime.Caller(skip)
	pcName := runtime.FuncForPC(pc).Name()
	msg := fmt.Sprintf("%s %d %s\n", file, line, pcName)
	return msg, ok
}

func fileTrack() string {
	// 追踪多层文件日志记录
	var msg string
	for i := 0; i < 32; i++ {
		m, ok := caller(i)
		if ok == false {
			break
		}
		// 过滤日志管理器的文件打印
		if strings.Index(m, "rifflock/lfshook") == -1&strings.Index(m, "sirupsen/logrus") {
			msg = msg + m
		}
	}
	return msg
}

func (s MyFormatter) Format(entry *log.Entry) ([]byte, error) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	level := strings.ToUpper(entry.Level.String())
	message := entry.Message
	var msg string
	if level == "ERROR" {
		file := fileTrack()
		msg = fmt.Sprintf("[%s]  [%s]  [%s]\n %s", timestamp, level, message, file)
	} else {
		msg = fmt.Sprintf("[%s]  [%s]  [%s]\n", timestamp, level, message)
	}
	return []byte(msg), nil
}

func setLogLevel(level string) {
	switch level {
	case "info":
		log.SetLevel(log.InfoLevel)
		break
	case "error":
		log.SetLevel(log.ErrorLevel)
		break
	case "debug":
		log.SetLevel(log.DebugLevel)
		break
	default:
		panic("Section \"log\" Option \"level\" error. choose (info, debug, error)")
	}
}

func setLogPath(logPath string, rotationTime, maxAge int) {
	rotation_time_ := time.Duration(rotationTime) * time.Hour
	max_age_ := time.Duration(maxAge) * time.Hour
	info_path := path.Join(logPath, "info.log")
	debug_path := path.Join(logPath, "debug.log")
	error_path := path.Join(logPath, "error.log")
	access, _ := rotatelogs.New(
		info_path+".%Y%m%d%H%M",
		rotatelogs.WithLinkName(info_path),
		rotatelogs.WithMaxAge(max_age_),
		rotatelogs.WithRotationTime(rotation_time_),
	)
	debug, _ := rotatelogs.New(
		debug_path+".%Y%m%d%H%M",
		rotatelogs.WithLinkName(debug_path),
		rotatelogs.WithMaxAge(max_age_),
		rotatelogs.WithRotationTime(rotation_time_),
	)
	err, _ := rotatelogs.New(
		error_path+".%Y%m%d%H%M",
		rotatelogs.WithLinkName(error_path),
		rotatelogs.WithMaxAge(max_age_),
		rotatelogs.WithRotationTime(rotation_time_),
	)

	lfHook := lfshook.NewHook(lfshook.WriterMap{
		log.DebugLevel: debug, // 为不同级别设置不同的输出目的
		log.InfoLevel:  access,
		log.ErrorLevel: err},
		&MyFormatter{})
	//&log.TextFormatter{})

	log.AddHook(lfHook)
}

func setLogConsole(console bool) {
	if console == false {
		src, _ := os.OpenFile(os.DevNull, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
		writer := bufio.NewWriter(src)
		log.SetOutput(writer)
	}
}
