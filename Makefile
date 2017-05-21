build:
	go build -o bin/jackdwyer jackdwyer.go

run:
	./bin/jackdwyer

refresh_db:
	mkdir ./db &>/dev/null
	aws --profile personal s3 cp $(S3_DB_LOCATION) db/app.db

get-deps:
	go get github.com/mattn/go-sqlite3
