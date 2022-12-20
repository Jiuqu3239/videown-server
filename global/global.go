package global

import (
	"videown-server/pkg/logger"
	"videown-server/pkg/setting"

	"gorm.io/gorm"
)

var Time_FMT = "2006/01/02 15:04"
var Settings *setting.Settings
var Logger *logger.Logger
var GormDb *gorm.DB

const COVER_IMAGE_PATH = "./cover_images/"
const DEFAULT_COVER_IMAGE = "./cover_images/default.png"
