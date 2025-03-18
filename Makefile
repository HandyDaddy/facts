DOCKER_USERNAME ?= fuu
APPLICATION_NAME ?= facts-service
BEARERTOKEN ?= test

.PHONY: docker_build
docker_build:
	docker build --tag ${DOCKER_USERNAME}/${APPLICATION_NAME} -f Dockerfile .

.PHONY: docker_run
docker_run:
	docker run -e BEARERTOKEN -p 8080:8080 -it --rm ${DOCKER_USERNAME}/${APPLICATION_NAME}