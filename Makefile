TARGET_ID?=
IPFS_ARGUMENTS?=

CONTAINER_IPFS?=ipfs
STAGING_DIR?=staging
DATA_DIR?=data

SWARM_KEY?=swam.key

.PHONY: up
up: down
	@docker-compose up -d

.PHONY: down
down:
	@docker-compose down -v

.PHONY: logs
logs:
	@docker-compose logs -f $(CONTAINER_IPFS)

.PHONY: peers
peers:
	@docker-compose exec $(CONTAINER_IPFS) ipfs swarm peers

.PHONY: id
id:
	@docker-compose exec $(CONTAINER_IPFS) ipfs id

.PHONY: bootstrap
bootstrap:
	@docker-compose exec $(CONTAINER_IPFS) ipfs bootstrap add $(TARGET_ID)

.PHONY: prepare
prepare:
	@go install github.com/Kubuxu/go-ipfs-swarm-key-gen/ipfs-swarm-key-gen@master
	@-mkdir $(DATA_DIR) $(STAGING_DIR)
	@ipfs-swarm-key-gen > $(DATA_DIR)/$(SWARM_KEY)

.PHONY: clean
clean:
	@-find $(DATA_DIR) $(STAGING_DIR) -type f,d -not -name '$(SWARM_KEY)' -delete

.PHONY: config
config:
	@docker-compose exec $(CONTAINER_IPFS) ipfs bootstrap rm --all
	@docker-compose exec $(CONTAINER_IPFS) ipfs config --bool Swarm.EnableRelayHop true

.PHONY: public-ip
public-ip:
	@-apt install dnsutils -y
	@dig +short myip.opendns.com @resolver1.opendns.com

.PHONY: ipfs
ipfs:
	@docker-compose exec $(CONTAINER_IPFS) ipfs $(IPFS_ARGUMENTS)

.PHONY: list-all
list-all:
	@docker-compose exec $(CONTAINER_IPFS) ipfs refs -r -u local
