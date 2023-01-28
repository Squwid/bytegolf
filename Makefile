build:
	docker pull golang:1.19.4-alpine3.17
	docker pull python:3.11.1-alpine3.17
	docker compose build

start: build
	mkdir -p pg
	docker compose up

clean:
	docker compose down
	docker compose rm --all

prune: clean
	docker system prune -a