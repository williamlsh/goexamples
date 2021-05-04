.PHONY: up
up: down
	@docker-compose up -d

.PHONY: down
down:
	@docker-compose down -v

.PHONY: logs
logs:
	@docker-compose logs -f ipfs

.PHONY: peers
peers:
	@docker-compose exec ipfs ipfs swarm peers

.PHONY: prepare
prepare:
	@-mkdir data staging
	@ipfs-swarm-key-gen > data/swarm.key

.PHONY: clean
clean:
	@-rm -rf data staging

.PHONY: config
config:
	@docker-compose exec ipfs ipfs bootstrap rm --all
	@docker-compose exec ipfs ipfs config --bool Swarm.EnableRelayHop true
