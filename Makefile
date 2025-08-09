templ:
	templ generate --watch --proxy="http://localhost:3000" --open-browser=false

server:
	air \
		--build.cmd "go build -o tmp/bin/main ./cmd/main.go" \
		--build.bin "tmp/bin/main" \
		--build.delay "100" \
		--build.exclude_dir "node_modules" \
		--build.include_ext "go" \
		--build.stop_on_error "false" \
		--misc.clean_on_exit true

tailwind:
	bunx @tailwindcss/cli -i ./views/css/input.css -o ./public/app.css --watch

dev:
	make -j3 tailwind templ server

