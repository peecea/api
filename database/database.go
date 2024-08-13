package database

import (
	"database/sql"
	"duval/database/db"
	"duval/internal/configuration"
	_ "github.com/go-sql-driver/mysql"
	"github.com/iancoleman/strcase"
	"github.com/jmoiron/sqlx"
	"time"
)

const (
	defaultDriver        = "mysql"
	maxOpenConnexion     = 120
	maxIdleConnexion     = 8
	maxConnexionLifeTime = time.Minute
)

var Client *sqlx.DB

func init() {
	var err error
	Client, err = sqlx.Connect(defaultDriver, configuration.App.DatabaseConnexionString)
	if err != nil {
		panic(err)
	}

	Client.SetMaxOpenConns(maxOpenConnexion)
	Client.SetMaxIdleConns(maxIdleConnexion)
	Client.SetConnMaxLifetime(maxConnexionLifeTime)
	Client.MapperFunc(strcase.ToSnake)
}

func Insert(T any) (lastId int64, err error) {
	var result sql.Result
	result, err = Client.Exec(db.I(T))
	if err != nil {
		return 0, err
	}

	lastId, err = result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return lastId, err
}

func Select(R any, Q string, A ...any) (err error) {
	err = Client.Select(R, Q, A...)
	if err != nil {
		return err
	}

	return err
}

func Get(R any, Q string, A ...any) (err error) {
	err = Client.Get(R, Q, A...)
	if err != nil {
		return err
	}
	return err
}

func InsertOne(T any) (id uint, err error) {
	lastId, err := Insert(T)
	if err != nil {
		return 0, err
	}

	return uint(lastId), err
}

func Update(T any) (err error) {
	_, err = Client.Exec(db.U(T))
	if err != nil {
		return err
	}
	return err
}

func Delete(T any) (err error) {
	_, err = Client.Exec(db.D(T))
	if err != nil {
		return err
	}

	return err
}

func Exec(Q string, A ...any) (err error) {
	_, err = Client.Exec(Q, A...)
	if err != nil {
		return err
	}
	return err
}

func InsertMany(T []any) (err error) {
	for i := 0; i < len(T); i++ {
		_, err = InsertOne(T[i])
		if err != nil {
			return err
		}
	}
	return err
}

func GetMany(R interface{}, Q string, A ...interface{}) (err error) {
	err = Client.Select(R, Q, A...)
	if err != nil {
		return err
	}
	return err
}
