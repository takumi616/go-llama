# Build docker image
build:
	docker compose build --no-cache

# Run image as a container
up:
	docker compose up