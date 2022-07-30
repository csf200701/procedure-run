build:
	go mod verify
	go mod tidy
	go build -o bin/procedure-run main.go

test:
	go test -race -v -failfast -test.timeout 5m ./...