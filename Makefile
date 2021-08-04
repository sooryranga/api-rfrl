#@IgnoreInspection BashAddShebang
export ROOT=$(realpath $(dir $(lastword $(MAKEFILE_LIST))))
export DEBUG=true
export APP=api-rfrl
export LDFLAGS="-w -s"

all: build test

build:
	go build -race  .

build-static:
	CGO_ENABLED=0 go build -race -v -o $(APP) -a -installsuffix cgo -ldflags $(LDFLAGS) .

start-rebuild:
	docker-compose up -d --force-recreate

start:
	docker-compose up -d

stop:
	docker-compose down

logs:
	docker-compose logs -f

docker-build:
	docker build .

############################################################
# Test
############################################################

test:
	go test -v -race ./...

container:
	docker build -t api-rfrl .

run-container:
	docker run --rm -it -p 8080:8080 -v ${ROOT}/id_rsa:/app/id_rsa api-rfrl 

.PHONY: build run build-static test container