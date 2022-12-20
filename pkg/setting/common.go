package setting

import (
	"path"
	"time"
)

const RunModeLive = "release"
const RunModeDev = "debug"

type ServerSettingS struct {
	RunMode      string
	HttpPort     string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

func (self *ServerSettingS) IsDebugMode() bool {
	return self.RunMode == RunModeDev
}

type AppSettingS struct {
	DefaultPageSize int
	MaxPageSize     int
	LogSavePath     string
	LogFileName     string
	LogFileExt      string
	CmpHttpUrl      string
	JwtSecret       string
	JwtDuration     int
	Username        string
	Password        string

	StaticPath string

	Limiter struct {
		IsOpen bool
		Count  int64
		Gap    int32
	}
}

func (self *AppSettingS) FullLogFilePath() string {
	return path.Join(self.LogSavePath, self.LogFileName+self.LogFileExt)
}
