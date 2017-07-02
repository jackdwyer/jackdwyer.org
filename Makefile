build:
	go build -o bin/jackdwyer jackdwyer.go

release: build

run: build
	./bin/jackdwyer

test:
	go test

refresh_db:
	mkdir -p ./db
	aws --profile personal s3 cp $(S3_DB_LOCATION) db/app.db

get-deps:
	go get -u github.com/aws/aws-sdk-go
	go get github.com/mattn/go-sqlite3
	go get github.com/nfnt/resize

dev-test:
	while :; do inotifywait jackdwyer_test.go; go test; done

SHELL=/bin/bash
get-test-data:
	mkdir -p ./test_data
	if [[ ! -f ./test_data/960x540.png ]]; then curl -o ./test_data/960x540.png https://dummyimage.com/960x540/8c8c8c/fff.jpg; fi
	if [[ ! -f ./test_data/961x541.png ]]; then curl -o ./test_data/961x541.png https://dummyimage.com/961x541/8c8c8c/fff.jpg; fi
	if [[ ! -f ./test_data/2160x1440.png ]]; then curl -o ./test_data/2160x1440.png https://dummyimage.com/2160x1440/8c8c8c/fff.jpg; fi
