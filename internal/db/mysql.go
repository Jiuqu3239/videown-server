package db

import (
	"videown-server/pkg/setting"

	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	_ "gorm.io/driver/mysql"
)

func NewGormDbForMySql(dbSetting *setting.DatabaseSettingS, debug bool) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s&parseTime=%t",
		dbSetting.UserName,
		dbSetting.Password,
		dbSetting.Host,
		dbSetting.DBName,
		dbSetting.Charset,
		dbSetting.ParseTime,
	)
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       dsn,
		Conn:                      nil,
		SkipInitializeWithVersion: false,
		DefaultStringSize:         256,
		DefaultDatetimePrecision:  nil,
		DisableDatetimePrecision:  false,
		DontSupportRenameIndex:    false,
		DontSupportRenameColumn:   false,
		DontSupportForShareClause: false,
	}), &gorm.Config{
		// QueryFields: true,
	})
	//db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	//db, err := gorm.Open(databaseSetting.DBType, )

	if err != nil {
		return nil, err
	}

	if debug {
		db = db.Debug()
	}
	//db.SingularTable(true)
	//db.DB().SetMaxIdleConns(databaseSetting.MaxIdleConns)
	//db.DB().SetMaxOpenConns(databaseSetting.MaxOpenConns)
	return db, nil
}
