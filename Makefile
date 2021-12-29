build:
	go build ./...

run:
	go run .

test:
	go test ./... -v

check-install-swagger:
	which ../../bin/swagger || GO111MODULE=on go install github.com/go-swagger/go-swagger/cmd/swagger@latest

swagger: check-install-swagger
	GO111MODULE=off ../../bin/swagger generate spec -o ./swagger.yaml --scan-models

check-lint:
	which ../../bin/golint || GO111MODULE=on go install golang.org/x/lint/golint/golint@latest

lint:
	GO111MODULE=off ../../bin/golint ./...

fmt:
	go fmt ./...

vet:
	go vet ./...	

build-image:
	docker build --tag scaleflixapi --no-cache=true .

all-images:
	docker images -a -q

remove-images:
	docker rmi $(shell docker images -a -q) 

remove-containers:
	docker rm $(shell docker ps -a -q) -f

remove-volume:
	docker volume rm $(shell docker volume ls -q)

run-docker:
	docker run -d -p 3000:3000 scaleflixapi

run-compose:
	docker-compose up --build

remove-compose:
	docker-compose down