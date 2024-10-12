run:
	INPUT_TARGET_DIRECTORY="example" INPUT_CONFIG_FILE="example/config.yaml" go run .

build:
	go build -ldflags="-s -w" -o out .

test:
	go test -v ./...

fmt:
	gofmt -l -s -w .
