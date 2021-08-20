.PHONY: cli
cli:
	@go run -race ./cli

.PHONY: server
server:
	@go run -race ./server

.PHONY: chore
chore:
	@go mod tidy
	@go mod download all
