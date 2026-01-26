.PHONY: compose-up compose-down

compose-up:
	docker compose up -d

compose-down:
	docker compose down