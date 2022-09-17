### Leafeon
A Wrapper around gorm for migrations

Still a work in progress


#### 🚀 Dependencies
- gorm
- go 1.18

#### How to use
-Import import into you project and reference github.com/dwakel/leafeon/migrator


-Can also use with terminal to run migrations

``` Run sample Migrate
    go run leafeon.go -connstr="host=localhost port=5432 user=postgres password=dbpassword dbname=dbname sslmode=disable" -src="./migrations" up
```
``` Run sample Rollback Migrations
    go run leafeon.go -connstr="host=localhost port=5432 user=postgres password=dbpassword dbname=dbname sslmode=disable" -src="./migrations" down
```

Replace -connstr with your connection string and -src with the path to you directory containing migrations files

##### Naming conventions
- Up migrations should be named {filename}.up.sql
- Rollback migrations should be named {filename}.down.sql

