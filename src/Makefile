.PHONY:	help start_server send_req stop_server clean launch-cluster \
				create-kubernetes-cluster kubernetes-test kubernetes-clean
.DEFAULT_GOAL: build

VERSION=v1
GOCMD=go
GOTEST=$(GOCMD) test
BINARY_NAME=birthday
PLATFORM ?= darwin
BINARY_UNIX=$(BINARY_NAME)_unix
DOCKER_REPO ?= gcr.io/birthday-server

## test: run tests
test: 
	@echo "  >  Running tests..."
	$(GOTEST) -v ./... 

## build[mac]: build binary
build: test
	@echo "  >  Building binary..."
	env GOOS=${PLATFORM} GOARCH=amd64 go build -o $(BINARY_NAME)

## docker-build[linux]: build inside docker container
docker-build:
	@echo "  >  Building ${BINARY_UNIX}..."
	docker run --rm -it -v "$(PWD)":/go golang:latest go build -o "$(BINARY_UNIX)" -v

start_server: build
	@echo "  >  Starting the server..."
	./birthday -port 8080 \
		2>/dev/null &

## send_req: build start_server and send test request

send_req: start_server
	@sleep 10
	@echo "  >  Sending test requests..."
	curl -s -X PUT -H "Content-Type: application/json" -d '{"birthday": "2018-02-01"}' http://localhost:8080/hello/suresh
	curl -s http://localhost:8080/hello/suresh|jq


stop_server:
	@killall -q birthday || true

clean: stop_server
	rm -f $(BINARY_NAME)

docker-clean:
	docker rmi ${NAME} &>/dev/null ||true

build-images: build
	docker build --pull=true --no-cache -t ${DOCKER_REPO}/birthday_server:$(VERSION) --rm .

## docker-images: push docker images to gcp using gcloud command
docker-images: build-images
	gcloud docker -- push ${DOCKER_REPO}/birthday_server:$(VERSION)

deploy-app:
	kubectl create -f birthday-server.yaml
	@echo "please wait deploy app...."
	@sleep 120


## kubernetes-test: deploy app to kubernetes
kubernetes-test: deploy-app
	curl "http://$(shell kubectl get svc birthday-server-svc|tail -1|awk '{print $4}')/"

kubernetes-clean:
	kubectl delete -f birthday-server.yaml

all: help
help: Makefile
	@echo
	@echo " Choose a command run :"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo
