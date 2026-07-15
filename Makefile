.PHONY: up test down clean

up:
	docker compose --env-file config.env up --build

down:
	docker compose --env-file down

clean:
	docker compose --env-file config.env down -v

test:
	docker compose --env-file test.env -f compose.test.yaml run --build --rm test

load-test:
	docker compose --env-file config.env run --rm load-test
