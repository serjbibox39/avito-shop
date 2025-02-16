PUBLIC_REGISTRY_HOST=docker.io
PUBLIC_REGISTRY_OWNER=serjbibox
PUBLIC_REGISTRY_APP_NAME=avito-shop

CI_COMMIT_REF_NAME=latest

all: deps build

deps:
	@go mod download
	@echo "Dependencies installed successfully"

build:
	go build ./cmd

image:
	@docker build -t ${PUBLIC_REGISTRY_HOST}/${PUBLIC_REGISTRY_OWNER}/${PUBLIC_REGISTRY_APP_NAME}:${CI_COMMIT_REF_NAME} ./
	@docker push ${PUBLIC_REGISTRY_HOST}/${PUBLIC_REGISTRY_OWNER}/${PUBLIC_REGISTRY_APP_NAME}:${CI_COMMIT_REF_NAME}
	

init: docker-down-clear \
	docker-pull docker-build docker-up \

docker-down-clear:
	docker-compose -f docker-compose.yml down -v --remove-orphans

docker-pull:
	docker-compose -f docker-compose.yml pull
	
docker-build:
	docker-compose -f docker-compose.yml build --pull

docker-up:
	docker-compose -f docker-compose.yml up -d

stop:
	docker-compose -f docker-compose.yml down
run:
	docker-compose -f docker-compose.yml up -d
log: 
	docker-compose logs -f -t		
test:
	docker-compose -f docker-compose_test.yml up --build --abort-on-container-exit
	docker-compose -f docker-compose_test.yml down --volumes

test-db-up:
	docker-compose -f docker-compose_test.yml up --build db

test-db-down:
	docker-compose -f docker-compose_test.yml down --volumes db