package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

var Logger = logrus.New()

func init() {
    Logger.SetLevel(logrus.DebugLevel)
	Logger.SetFormatter(&logrus.JSONFormatter{})

    file, err := os.OpenFile("/var/log/myapp/app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
    if err != nil {
        Logger.SetOutput(os.Stdout) 
        Logger.Error("Не удалось открыть файл логов, вывод логов на стандартный вывод")
    } else {
        Logger.SetOutput(file)
        defer file.Close()
    }

    Logger.Debug("Функция init() выполнена")
}
