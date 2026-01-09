build:
	go build -o  bin/rapidvpp-server ./cmd/app

run:
	@CGO_ENABLED=0 go run ./cmd/app/main.go

coverage:
	go test -count=1 -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	rm coverage.out
test:
	go test -race ./...