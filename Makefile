include .env
export

run: clean
	go run cmd/main.go

clean:
	gofmt -w .
