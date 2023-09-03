package framework

import (
	"github.com/Unknwon/goconfig"
	log "github.com/sirupsen/logrus"
)

type FrameworkConf struct {
	LogFilePath string
	LogLevel    log.Level
	MaxSize     int
	MaxBackups  int
	MaxAge      int
	Compress    bool

	HttpPort         int
	HttpStaticDir    string
	HttpStaticPrefix string

	HttpsPort int
	HttpsCert string
	HttpsKey  string
}

var Conf = &FrameworkConf{}

func LoadConf(confFile string) error {
	var err error
	configFile, err := goconfig.LoadConfigFile(confFile)

	if err != nil {
		return err
	}

	Conf.LogFilePath, err = configFile.GetValue("log", "logFilePath")
	if err != nil {
		return err
	}

	level, err := configFile.GetValue("log", "logLevel")
	if err != nil {
		return err
	}
	Conf.LogLevel = converLogLevel(level)

	Conf.MaxSize, err = configFile.Int("log", "maxSize")
	if err != nil {
		return err
	}

	Conf.MaxBackups, err = configFile.Int("log", "maxBackups")
	if err != nil {
		return err
	}

	Conf.MaxBackups, err = configFile.Int("log", "maxAge")
	if err != nil {
		return err
	}

	Conf.MaxAge, err = configFile.Int("log", "maxBackups")
	if err != nil {
		return err
	}

	Conf.Compress, err = configFile.Bool("log", "compress")
	if err != nil {
		return err
	}

	Conf.HttpPort, err = configFile.Int("http", "port")
	if err != nil {
		return err
	}

	Conf.HttpStaticDir, err = configFile.GetValue("http", "staticDir")
	if err != nil {
		return err
	}

	Conf.HttpStaticPrefix, err = configFile.GetValue("http", "staticPrefix")
	if err != nil {
		return err
	}

	Conf.HttpsPort, err = configFile.Int("https", "port")
	if err != nil {
		return err
	}

	Conf.HttpsCert, err = configFile.GetValue("https", "cert")
	if err != nil {
		return err
	}

	Conf.HttpsKey, err = configFile.GetValue("https", "key")
	if err != nil {
		return err
	}

	return err
}

func converLogLevel(logLevel string) log.Level {
	switch logLevel {
	case "Debug":
		return log.DebugLevel
	case "Info":
		return log.InfoLevel
	case "Warn":
		return log.WarnLevel
	case "Error":
		return log.ErrorLevel
	default:
		return log.InfoLevel
	}
}
