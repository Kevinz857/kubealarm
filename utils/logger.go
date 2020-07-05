package utils

import (
	"kubealarm/conf"
	"path"
	"time"

	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

/*
func desp:new a logrus object, add a hook.
input: logPath eagle日志目录; maxAge 日志留存时间; rotationTime 日志切割轮转时间.
output: logrus Logger对象
*/
func InitLogger(logPath string, maxAge time.Duration, rotationTime time.Duration) *logrus.Logger {
	if Log != nil {
		return Log
	}

	infoBaseLogPaht := path.Join(logPath, conf.InfoLogFileName)
	infoWriter, err := rotatelogs.New(
		infoBaseLogPaht+".%Y%m%d%H%M",
		rotatelogs.WithLinkName(infoBaseLogPaht),
		rotatelogs.WithMaxAge(maxAge),
		rotatelogs.WithRotationTime(rotationTime),
	)
	if err != nil {
		logrus.Errorf("[InitLogger] info logger file init failed, err:%+s", err.Error())
		return nil
	}
	errBaseLogPaht := path.Join(logPath, conf.ErrLogFileName)
	errWriter, err := rotatelogs.New(
		errBaseLogPaht+".%Y%m%d%H%M",
		rotatelogs.WithLinkName(errBaseLogPaht),
		rotatelogs.WithMaxAge(maxAge),
		rotatelogs.WithRotationTime(rotationTime),
	)
	if err != nil {
		logrus.Errorf("[InitLogger] err logger file init failed, err:%+s", err.Error())
		return nil
	}

	lfHook := lfshook.NewHook(lfshook.WriterMap{
		logrus.DebugLevel: infoWriter,
		logrus.InfoLevel:  infoWriter,
		logrus.WarnLevel:  infoWriter,
		logrus.ErrorLevel: errWriter,
		logrus.FatalLevel: errWriter,
		logrus.PanicLevel: errWriter,
	}, &logrus.TextFormatter{})

	Log = logrus.New()
	Log.Hooks.Add(lfHook)
	return Log
}
