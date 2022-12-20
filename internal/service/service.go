package service

import (
	"videown-server/internal/service/auth"
	"videown-server/pkg/logger"
	"videown-server/pkg/setting"

	"gorm.io/gorm"
)

func SetupService(settings *setting.Settings, gorm *gorm.DB, log *logger.Logger) {
	auth.SetupAuth(settings.AppSetting.JwtSecret,
		int64(settings.AppSetting.JwtDuration),
		auth.DefaultAdminVerifyProvider{Settings: settings},
	)
}
