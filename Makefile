HTTP_PROXY?=""

.PHONY: all build down logs

all: build
	@docker-compose up -d

build:
	@COMPOSE_DOCKER_CLI_BUILD=1 docker-compose build --build-arg HTTP_PROXY=${HTTP_PROXY}

down:
	@docker-compose down -v

logs:
	@docker-compose logs -f hotrod