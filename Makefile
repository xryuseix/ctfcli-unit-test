run:
	INPUT_TARGET_DIRECTORY="example" INPUT_CONFIG_FILE="config.yaml" go run .

build:
	GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o out .

test:
	go test -v ./...

fmt:
	gofmt -l -s -w .
