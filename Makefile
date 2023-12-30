run:
	INPUT_TARGET_DIRECTORY="example" go run .

build:
	GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o out .

test:
	go test -v ./...

fmt:
	gofmt -l -s -w .
