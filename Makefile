build:
	go build -i -o bin/jackdwyer *.go

release: build

run: build
	./bin/jackdwyer

test:
	go test

refresh_db:
	mkdir -p ./db
	aws --profile personal s3 cp $(S3_DB_LOCATION) db/app.db

get-deps:
	go get github.com/aws/aws-sdk-go
	go get github.com/rwcarlsen/goexif/exif
	go get github.com/codingsince1985/geo-golang
	go get github.com/codingsince1985/geo-golang/openstreetmap
	go get github.com/mattn/go-sqlite3
	go get github.com/disintegration/imaging
	go get github.com/gorilla/mux


dev-test:
	while :; do inotifywait jackdwyer_test.go; go test; done

SHELL=/bin/bash
get-test-data:
	mkdir -p ./test_data
	if [[ ! -f ./test_data/960x540.png ]]; then curl -o ./test_data/960x540.png https://dummyimage.com/960x540/8c8c8c/fff.jpg; fi
	if [[ ! -f ./test_data/961x541.png ]]; then curl -o ./test_data/961x541.png https://dummyimage.com/961x541/8c8c8c/fff.jpg; fi
	if [[ ! -f ./test_data/2160x1440.png ]]; then curl -o ./test_data/2160x1440.png https://dummyimage.com/2160x1440/8c8c8c/fff.jpg; fi

validate-dev-deploy:
	curl --fail -o/dev/null -s http://dev.jackdwyer.org/

sync-production-development:
	aws s3 sync s3://images.jackdwyer.org/960/ s3://dev-images.jackdwyer.org/960 --profile personal

TODOS:
	ag --ignore Makefile TODO
