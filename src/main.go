package main

import (
	"io"
	"os"
	"path"
	"runtime"
	"signaling/src/framework"
	"strconv"

	log "github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

func init() {
	err := framework.LoadConf("./conf/framework.ini")
	if err != nil {
		log.Errorf("load config error %s", err.Error())
		panic(err)
	}
	log.Infof("load config success %+v", framework.Conf)
	logFormatter := new(log.TextFormatter)
	logFormatter.FullTimestamp = true
	logFormatter.TimestampFormat = "2006-01-02 15:04:05.000000"
	logFormatter.CallerPrettyfier = func(frame *runtime.Frame) (function string, file string) {
		fileLine := path.Base(frame.File) + ":" + strconv.Itoa(frame.Line)
		return "", fileLine
	}
	log.SetLevel(log.InfoLevel)
	log.SetReportCaller(true)
	log.SetFormatter(logFormatter)
	logFile := &lumberjack.Logger{
		Filename:   framework.Conf.LogFilePath,
		MaxSize:    framework.Conf.MaxSize, // megabytes
		MaxBackups: framework.Conf.MaxBackups,
		MaxAge:     framework.Conf.MaxAge, //days
		Compress:   framework.Conf.Compress,
	}
	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)
}

func main() {
	// 静态资源处理 /static
	framework.RegisterStaticUrl()

	// 启动http server
	go startHttp()

	startHttps()
}

func startHttp() {
	err := framework.StartHttp()
	if err != nil {
		panic(err)
	}
}

func startHttps() {
	err := framework.StartHttps()
	if err != nil {
		panic(err)
	}
}
