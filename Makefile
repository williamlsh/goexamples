.PHONY: all build down

all: build
	@sudo docker-compose up --build -d

build:
	@go build -o hotrod-linux-amd64 .

down:
	@sudo docker-compose down -v
	@rm hotrod-linux-amd64