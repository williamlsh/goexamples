.PHONY: up
up: down
	@docker-compose up -d

.PHONY: down
down:
	@docker-compose down -v --remove-orphan

.PHONY: image
image:
	@docker build -t demo .

.PHONY: run
run:
	@docker run --rm --name demo demo
