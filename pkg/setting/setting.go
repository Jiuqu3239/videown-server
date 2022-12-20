package setting

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Settings struct {
	ServerSetting   *ServerSettingS
	AppSetting      *AppSettingS
	Web3Setting     *Web3SettingS
	DatabaseSetting *DatabaseSettingS
	RedisSetting    *RedisSettingS
	vp              *viper.Viper
}

func (t *Settings) ChangePassword(newPassword, oldPassword string) error {
	if len(newPassword) < 8 {
		return fmt.Errorf("the new password length must greate and equal than 8")
	}
	if t.AppSetting.Password != oldPassword {
		return fmt.Errorf("the old password is invalid")
	}
	t.AppSetting.Password = newPassword
	t.vp.Set("app.password", newPassword)
	if err := t.vp.WriteConfig(); err != nil {
		return err
	}
	return nil
}

func loadConfigFile(configPath string) (*viper.Viper, error) {
	vp := viper.New()
	vp.SetConfigName("conf")
	vp.AddConfigPath(configPath)
	vp.SetConfigType("yaml")
	err := vp.ReadInConfig()
	if err != nil {
		return nil, err
	}
	return vp, nil
}

func NewSettings() (*Settings, error) {
	return NewSettingsWithDirectory("configs/")
}

func NewSettingsWithDirectory(configDir string) (*Settings, error) {
	vp, err := loadConfigFile(configDir)
	if err != nil {
		return nil, err
	}
	pSettings := new(Settings)
	pSettings.vp = vp
	if err := vp.UnmarshalKey("Server", &pSettings.ServerSetting); err != nil {
		return nil, err
	}
	pSettings.ServerSetting.ReadTimeout *= time.Second
	pSettings.ServerSetting.WriteTimeout *= time.Second

	if err := vp.UnmarshalKey("App", &pSettings.AppSetting); err != nil {
		return nil, err
	}
	if err := vp.UnmarshalKey("Database", &pSettings.DatabaseSetting); err != nil {
		return nil, err
	}
	if err := vp.UnmarshalKey("Web3", &pSettings.Web3Setting); err != nil {
		return nil, err
	}
	if err := vp.UnmarshalKey("Redis", &pSettings.RedisSetting); err != nil {
		return nil, err
	}

	return pSettings, nil
}
