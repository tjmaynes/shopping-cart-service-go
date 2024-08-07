ENV_FILE := $(or $(ENV_FILE), .env.development)

include $(ENV_FILE)
export $(shell sed 's/=.*//' $(ENV_FILE))

export CGO_ENABLED=1
export REGISTRY_USERNAME=tjmaynes
export IMAGE_NAME=shopping-cart-service-go
export TAG=$(shell git rev-parse --short HEAD)

install:
	chmod +x ./scripts/install.sh
	./scripts/install.sh

generate_mocks:
	moq -out internal/pkg/item/repository_mock.go internal/pkg/item Repository

generate_seed_data:
	go run ./internal/cmd/shopping-cart-service-seeder \
	--seed-data-destination=${PWD}/internal/db/seed.json \
	--item-count=100 \
	--manufacturer-count=5

migrate:
	DATABASE_URL=$(DATABASE_URL) bin/dbmate --migrations-dir=${PWD}/internal/db/migrations wait
	DATABASE_URL=$(DATABASE_URL) bin/dbmate --migrations-dir=${PWD}/internal/db/migrations up
	DATABASE_URL=$(DATABASE_URL) bin/dbmate --migrations-dir=${PWD}/internal/db/migrations migrate

seed:
	go run ./internal/cmd/shopping-cart-service-db-seeder \
	--db-source=$(DATABASE_URL) \
	--seed-data-source=${PWD}/internal/db/seed.json

test: migrate generate_mocks
	mkdir -p coverage || true
	DATABASE_URL=$(DATABASE_URL) PORT=$(PORT) \
	SEED_DATA_SOURCE=${PWD}/internal/db/seed.json \
	go test -v -coverprofile=coverage/coverage.txt ./internal/...

ci_test:
	mkdir -p coverage || true
	make test 2>&1 | go-junit-report -out coverage/report.xml
	gocov convert coverage/coverage.txt > coverage/coverage.json
	gocov-xml < coverage/coverage.json > coverage/coverage.xml
	gocov-html < coverage/coverage.json > coverage/index.html

build:
	go build -o dist/shopping-cart-service ./internal/cmd/shopping-cart-service

start: build migrate
	DATABASE_URL=$(DATABASE_URL) PORT=$(PORT) ./dist/shopping-cart-service

build_image:
	chmod +x ./scripts/build-image.sh
	./scripts/build-image.sh

push_image:
	chmod +x ./scripts/push-image.sh
	./scripts/push-image.sh

debug_image:
	chmod +x ./scripts/debug-image.sh
	./scripts/debug-image.sh

deploy_k8s:
	kubectl apply -f ./k8s/shopping-cart-common/secret.yml
	kubectl apply -f ./k8s/shopping-cart-db/deployment.yml
	kubectl apply -f ./k8s/shopping-cart-db/persistence.gke.yml

connect_localhost_to_remote_db:
	kubectl port-forward svc/shopping-cart-db 5432:5432

format:
	go fmt github.com/tjmaynes/shopping-cart-service-go

deploy: install ci_test

run_local_db:
	docker compose up

debug_local_db:
	docker run -it --rm \
		--network shopping-cart-service-go_shopping-cart-network \
		postgres:16.3-alpine \
		psql \
		-h shopping-cart-db \
		--username postgres

stop_local_db:
	docker compose down
	docker volume rm shopping-cart-service-go_shopping-cart-db

ship_it: build test
	git push

clean:
	rm -rf dist/ vendor/ coverage* report.xml
