package setting

import "time"

type DatabaseSettingS struct {
	DBType    string
	UserName  string
	Password  string
	Host      string
	DBName    string
	Charset   string
	ParseTime bool

	MaxIdleConns int
	MaxOpenConns int
}

type RedisSettingS struct {
	AddressPort string
	Password    string
	DefaultDB   int
	DialTimeout time.Duration
}

type Web3SettingS struct {
	ChainId       int64
	RpcEndpoints  []string
	Mnemonic      string
	SupperAddress string
}
