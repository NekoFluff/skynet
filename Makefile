gen:
	go generate ./...

test:
	go test -cover ./...
	golangci-lint run --print-issued-lines=false --out-format=colored-line-number --issues-exit-code=0 ./...

test-cover:
	go test -coverprofile cover.out ./...
	go tool cover -html=cover.out -o cover.html

build:
	go build cmd/bot/main.go

run:
	go mod vendor
	docker compose --file ./docker-compose.yml --project-directory . up --build \
	docker compose --file ./docker-compose.yml --project-directory . down --volumes \
	rmdir /s /q vendor

up:
	go mod vendor
	docker compose --file ./docker-compose.yml --project-directory . up --build -d
	rmdir /s /q vendor

down:
	docker compose --file ./docker-compose.yml --project-directory . down --volumes

deploy-up:
	docker image pull nekofluff/skynet
	docker compose --file ./docker-compose.deploy.yml --project-directory . up -d

deploy-down:
	docker compose --file ./docker-compose.deploy.yml --project-directory . down