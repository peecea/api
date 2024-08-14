package main

import (
	"github.com/jmoiron/sqlx"
	"peec/database"
	"peec/internal/route"
)

func main() {
	defer func(Client *sqlx.DB) {
		err := Client.Close()
		if err != nil {
			return
		}
	}(database.Client)

	err := route.Serve()
	if err != nil {
		panic(err)
	}
}
