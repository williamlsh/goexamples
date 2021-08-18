.PHONY: cli
cli:
	@go run ./cli

.PHONY: server
server:
	@go run ./server

.PHONY: chore
chore:
	@go mod tidy
	@go mod download all
