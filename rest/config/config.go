package config

//Cfg
var Cfg *Config

func init() {
	Cfg = NewDefaultConfig()
}

//Config
type Config struct {
	DBEngine string
	DBName   string
	DBUser   string
	DBPWD    string
	DBHost   string
	DBPort   string
	DBZone   string
}

func NewDefaultConfig() *Config {
	return &Config{
		DBEngine: "mysql",
		DBName:   "mydb",
		DBUser:   "root",
		DBPWD:    "root",
		DBHost:   "192.168.33.10",
		DBPort:   "3306",
		DBZone:   "Asia/Shanghai",
	}
}
