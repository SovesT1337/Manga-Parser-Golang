module x.localhost/scripts

go 1.23.11

require (
	github.com/joho/godotenv v1.5.1
	golang.org/x/net v0.42.0
	gorm.io/gorm v1.30.0
)

require (
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/mattn/go-sqlite3 v1.14.22 // indirect
	golang.org/x/text v0.27.0 // indirect
	gorm.io/driver/sqlite v1.6.0
)

replace x.localhost/scripts/parsers => ./parsers

replace x.localhost/scripts/telegraph => ./telegraph

replace x.localhost/scripts/telegram => ./telegram

replace x.localhost/scripts/database => ./database

replace x.localhost/scripts/schemas => ./schemas
