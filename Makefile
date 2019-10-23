all: build

build:
	@goimports -w .
	@gofmt -s -w .
	@go build -o hubit-mux

run: build
	@./hubit-mux