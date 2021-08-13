.PHONY: image
image:
	@docker build -t srt .

.PHONY: run
run:
	@docker run -d --rm --name srt -p 8090:8090/udp srt

.PHONY: stop
stop:
	@docker stop srt

.PHONY: logs
logs:
	@docker logs -f srt
