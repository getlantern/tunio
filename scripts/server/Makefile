docker:
	docker build -t getlantern/tunio-proxy .

docker-run: docker
	(docker stop tunio-proxy || exit 0) && \
	(docker rm tunio-proxy || exit 0) && \
	docker run \
		-d \
		-p 0.0.0.0:2099:2099 \
		-p 0.0.0.0:5353:5353 \
		--name tunio-proxy \
		-t getlantern/tunio-proxy
