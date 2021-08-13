.PHONY: image
image:
	@docker build -t srt/echo .

.PHONY: run
run:
	@docker run -d --rm --name srt -p 8090:8090/udp srt/echo

.PHONY: stop
stop:
	@docker stop srt

.PHONY: logs
logs:
	@docker logs -f srt
