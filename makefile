OWNER=novacloud
IMAGE_NAME=devops
QNAME=$(OWNER)/$(IMAGE_NAME)

GIT_TAG=$(QNAME):$(TRAVIS_COMMIT)
BUILD_TAG=$(QNAME):0.1.$(TRAVIS_BUILD_NUMBER)
LATEST_TAG=$(QNAME):latest

lint:
	docker run -it --rm -v "$(PWD)/Dockerfile:/Dockerfile:ro" redcoolbeans/dockerlint

build:
	# go get ./...
	# gox -osarch=z"linux/amd64" -output="bin/devops-alpine"
	GO111MODULE=on CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/devops-alpine .
	docker build -t $(GIT_TAG) .
	docker build -t $(GIT_TAG)-golang -f Dockerfile.golang .
	docker build -t $(GIT_TAG)-aws -f Dockerfile.aws .
	docker build -t $(GIT_TAG)-mysql -f Dockerfile.mysql .
	docker build -t $(GIT_TAG)-docker-compose -f Dockerfile.docker-compose .

tag:
	docker tag $(GIT_TAG) $(BUILD_TAG)
	docker tag $(GIT_TAG)-golang $(BUILD_TAG)-golang
	docker tag $(GIT_TAG)-aws $(BUILD_TAG)-aws
	docker tag $(GIT_TAG)-mysql $(BUILD_TAG)-mysql
	docker tag $(GIT_TAG)-docker-compose $(BUILD_TAG)-docker-compose
	docker tag $(GIT_TAG) $(LATEST_TAG)
	docker tag $(GIT_TAG)-golang $(LATEST_TAG)-golang
	docker tag $(GIT_TAG)-aws $(LATEST_TAG)-aws	
	docker tag $(GIT_TAG)-mysql $(LATEST_TAG)-mysql	
	docker tag $(GIT_TAG)-docker-compose $(LATEST_TAG)-docker-compose	

login:
	@docker login -u "$(DOCKER_USER)" -p "$(DOCKER_PASS)"
push: login
	# docker push $(GIT_TAG)
	# docker push $(BUILD_TAG)
	docker push $(LATEST_TAG)
	docker push $(LATEST_TAG)-golang
	docker push $(LATEST_TAG)-aws
	docker push $(LATEST_TAG)-mysql
	docker push $(LATEST_TAG)-docker-compose


build-local:
	# go get ./...
	go build -o devops

deploy-local:
	make build-local
	mv devops /usr/local/bin/
