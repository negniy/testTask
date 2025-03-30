package config

import (
	"os"

	"github.com/sirupsen/logrus"
)

var Logger = logrus.New()

func LoadLoger() {

	Logger.SetLevel(logrus.InfoLevel)

	Logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	file, err := os.OpenFile("task.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		Logger.SetOutput(file)
	} else {
		Logger.Warn("Не удалось открыть файл логов, используется стандартный stderr")
	}
}
