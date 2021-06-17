.PHONY: run
run:
	@go run -race .

.PHONY: bufgen
bufgen: bufdep
	@buf generate
	@go mod tidy

.PHONY: bufdep
bufdep:
	@buf beta mod update
