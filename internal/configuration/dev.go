package configuration

import (
	"fmt"
	"github.com/BurntSushi/toml"
	go_console "github.com/DrSmithFr/go-console"
)

func InitiateConfiguration() (err error) {
	cmd := go_console.NewScript().Build()

	_, err = toml.DecodeFile("config.toml", &App)
	if err != nil {
		App = Config{
			Version:                 "v0.0.1",
			Port:                    "8087",
			Host:                    "",
			RunningMode:             RunningModeProd,
			TokenSecret:             "437b059d-bd8b-40d5-920a-341bb8a3f15f",
			DatabaseUserName:        "root",
			DatabaseUserPassword:    "UnderAll4",
			DatabaseName:            "duval",
			DatabaseHost:            "db-duval.ctw4aeiceahd.eu-north-1.rds.amazonaws.com",
			DatabasePort:            "3306",
			DatabaseConnexionString: "",
		}
		fmt.Println("message: \"Your are using production configuration. If its not what you want, please configure your config.toml file.\"")
		cmd.PrintWarnings([]string{
			"Your are using production configuration. If its not what you want, please configure your config.toml file.",
		})

		err = nil
	}

	cmd.PrintNotes([]string{
		fmt.Sprintf("version : %s", App.Version),
		fmt.Sprintf("port : %s", App.Port),
		fmt.Sprintf("host : %s", App.Host),
		fmt.Sprintf("token secret : %s", App.TokenSecret),
		fmt.Sprintf("running mode : %d", App.RunningMode),
		fmt.Sprintf("database host : %s", App.DatabaseHost),
		fmt.Sprintf("database DatabasePort : %s", App.DatabasePort),
	})

	App.DatabaseConnexionString = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", App.DatabaseUserName,
		App.DatabaseUserPassword, App.DatabaseHost, App.DatabasePort,
		App.DatabaseName)

	return err
}
