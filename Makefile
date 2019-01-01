
GO_IMAGE=golang:1.10.4
SWAGGER_IMAGE=quay.io/goswagger/swagger:v0.17.2

REPO_NAME=task2trip
REPO_PATH=github.com/itimofeev/$(REPO_NAME)


deploy-db:
	docker stack deploy --compose-file tools/db.docker-stack.yml db

deploy:
	docker stack deploy --compose-file tools/docker-stack.yml db

gen-server: download
	docker run --rm -v $(GOPATH):/go/ -w /go/src/$(REPO_PATH) -t $(SWAGGER_IMAGE) \
		generate server \
		--target=server \
		-f tools/swagger.yml

gen-client:
	docker run --rm -v $(GOPATH):/go/ -w /go/src/$(REPO_PATH) -t $(SWAGGER_IMAGE) \
		generate client \
		--target=$(REPO_NAME) \
		-f tools/swagger.yml

download:
	wget -O tools/swagger.yml\
    		--header="Accept: application/yaml" \
    		https://api.swaggerhub.com/apis/itimofeev/task2trip/1.0.0

rm:
	docker service rm $(shell docker service ls -q) || true

rm-volume:
	docker volume rm db_pg_data || true

rm-containers:
	docker rm $(shell docker ps -a -f status=exited -q) || true

release: build-docker build-image upload clean run-remote

upload:
	scp -r $(REPO_NAME).img root@159.69.121.222:/root/$(REPO_NAME)
	scp -r tools/docker-stack.yml root@159.69.121.222:/root/$(REPO_NAME)/tools
	scp -r tools/cmd.sql root@159.69.121.222:/root/$(REPO_NAME)/tools
	scp -r Makefile root@159.69.121.222:/root/$(REPO_NAME)
	ssh root@159.69.121.222 "cd $(REPO_NAME); ls"

run-remote:
	ssh root@159.69.121.222 "cd $(REPO_NAME); \
		docker load -i $(REPO_NAME).img; \
		make rm; \
		make deploy \
		"

clean:
	rm $(REPO_NAME) $(REPO_NAME).img

build-docker:
	docker run \
		-v $$GOPATH:/go \
		-w /go/src/$(REPO_PATH) \
		-t $(GO_IMAGE) \
		make build
	docker cp `docker ps -q -n=1`:/go/bin/linux_amd64/$(REPO_NAME) ./

build-image:
	docker build  \
		--force-rm=true \
		-t $(REPO_NAME) \
		-f tools/Dockerfile `pwd`
	docker image save $(REPO_NAME) -o $(REPO_NAME).img

import-products:
	docker cp tools/cmd.sql `docker ps -q -f name=db_db`:/cmd.sql \
		&& docker exec `docker ps -q -f name=db_db` psql postgres postgres -f /cmd.sql

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o $$GOPATH/bin/linux_amd64/task2trip $(REPO_PATH)/server/cmd/api-for-task2trip-server