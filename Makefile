build:
	go mod verify
	go mod tidy
	go build -o bin/admin-cli main.go

test:
	go test -race -v -failfast -test.timeout 5m ./...