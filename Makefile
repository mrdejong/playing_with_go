server:
	air

tailwind:
	bunx @tailwindcss/cli -i ./views/css/input.css -o ./public/app.css --watch
dev:
	make -j2 tailwind server

