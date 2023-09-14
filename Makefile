BINARY_NAME=urlshortner
DSN="host=localhost port=5432 user=postgres password=mysecretpassword dbname=urlshortner sslmode=disable timezone=UTC connect_timeout=5"
DB_URL='postgres://postgres:mysecretpassword@localhost/urlshortner?sslmode=disable'
## build: Build binary
build:
	@echo "Building..."
	env CGO_ENABLED=0  go build -ldflags="-s -w" -o ${BINARY_NAME} ./cmd/api
	@echo "Built!"

## run: builds and runs the application
run: build
	@echo "Starting..."
	@env URLSHORTNER_DB_DSN=${DSN} ./${BINARY_NAME} &
	@echo "Started!"

## clean: runs go clean and deletes binaries
clean:
	@echo "Cleaning..."
	@go clean
	@rm ${BINARY_NAME}
	@echo "Cleaned!"

## start: an alias to run
start: run

## stop: stops the running application
stop:
	@echo "Stopping..."
	@-pkill -SIGTERM -f "./${BINARY_NAME}"
	@echo "Stopped!"

## restart: stops and starts the application
restart: stop start


migrateup:
	migrate -path ./migrations -database $(DB_URL) up

## test: runs all tests
test:
	go test -v ./...