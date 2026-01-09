.PHONY: test test-coverage lint build clean

test:
	go test -v -cover ./...

test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

lint:
	golangci-lint run

build:
	go build ./...

clean:
	rm -f coverage.out coverage.html

