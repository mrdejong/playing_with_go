run: build
	@./bin/app

build:
	@go build -o bin/app ./cmd/main.go

css:
	@bunx @tailwindcss/cli -i ./views/css/input.css -o ./public/app.css --watch

dev:
	@air
