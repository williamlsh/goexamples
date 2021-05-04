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

data/swarm.key:
	@go install github.com/Kubuxu/go-ipfs-swarm-key-gen/ipfs-swarm-key-gen@master
	@-mkdir staging data > /dev/null 2>&1
	@ipfs-swarm-key-gen > data/swarm.key
