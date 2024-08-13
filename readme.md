- Setup
  - 
    - golang version 1.22.x (https://go.dev/dl/)
    - mysql 8.0

  one installed, pull the directory ( https://github.com/cend-org/duval.git ) and run the following command:

- `go mod tidy` to download all Project dépendances and avoid pkg update. If you got a problem just run `go get all` .
- Create a database with the name of `duval`.
- mysql migrator are stored in `data/mysql/migrator.sql`.  It creates the database duval and its tables. 
- in the root project create a config.toml file [ template in config-tpl.go]
