package configuration

const (
	RunningModeTest = 0
	RunningModeDev  = 1
	RunningModeProd = 2
)

type Config struct {
	Version                 string `toml:"version"`
	RunningMode             int    `toml:"running_mode"`
	Port                    string `toml:"port"`
	Host                    string `toml:"host"`
	TokenSecret             string `toml:"token_secret"`
	DatabaseUserName        string `toml:"database_user_name"`
	DatabaseUserPassword    string `toml:"database_user_password"`
	DatabaseName            string `toml:"database_name"`
	DatabaseHost            string `toml:"database_host"`
	DatabasePort            string `toml:"database_port"`
	DatabaseConnexionString string
}

var App Config

func init() {
	err := InitiateConfiguration()
	if err != nil {
		panic(err)
	}
}

func (c *Config) IsDev() bool {
	return c.RunningMode == RunningModeDev
}

func (c *Config) IsProd() bool {
	return c.RunningMode == RunningModeProd
}

func (c *Config) IsTest() bool {
	return c.RunningMode == RunningModeTest
}
