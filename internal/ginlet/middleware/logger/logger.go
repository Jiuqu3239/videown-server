package logger

import (
	"fmt"
	"os"
	"time"
	"videown-server/global"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func LoggerToFile() gin.HandlerFunc {
	fileName := global.Settings.AppSetting.FullLogFilePath()

	// Read/write mode of log files
	f, err := os.OpenFile(fileName, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("err:", err)
	}

	logger := logrus.New()
	// set out
	logger.Out = f
	// Setting a Log Level
	logger.SetLevel(logrus.TraceLevel)
	// Setting the Log Format
	logger.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: global.Time_FMT,
	})

	return func(c *gin.Context) {
		startTime := time.Now()               // start time
		c.Next()                              // Handling the Request
		endTime := time.Now()                 // end time
		latencyTime := endTime.Sub(startTime) // execution time
		reqMethod := c.Request.Method         // request method
		reqUri := c.Request.RequestURI        // required parameter
		statusCode := c.Writer.Status()       // status code
		clientIP := c.ClientIP()              // IP
		logger.Infof(" %13v | %15s | %7s | %3d | %s ",
			latencyTime,
			clientIP,
			reqMethod,
			statusCode,
			reqUri,
		)
	}
}
