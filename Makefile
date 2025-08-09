run: build
	@./bin/app

build:
	@go build -o bin/app ./cmd/main.go

web:
	@cd ./web && bun run dev

dev:
	@air
