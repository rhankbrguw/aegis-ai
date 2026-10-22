.PHONY: all setup docker-up docker-down run-backend run-mobile test test-backend test-mobile build-backend build-mobile

setup:
	cp -n deployments/.env.example deployments/.env || true
	cd services/gateway && go mod download
	cd mobile && flutter pub get

docker-up:
	docker compose -f deployments/docker-compose.yml up -d --build

docker-down:
	docker compose -f deployments/docker-compose.yml down

run-backend:
	cd services/gateway && go run cmd/server/main.go

run-mobile:
	cd mobile && flutter run

run-mobile-linux:
	cd mobile && flutter run -d linux

test: test-backend test-mobile

test-backend:
	cd services/gateway && go test -v -race -cover ./...

test-mobile:
	cd mobile && flutter test

check:
	cd services/gateway && go vet ./...
	cd mobile && flutter analyze
	@$(MAKE) test
	@echo "[SUCCESS] All constraints verified successfully!"

build-backend:
	cd services/gateway && CGO_ENABLED=0 go build -o bin/aegis-gateway cmd/server/main.go

build-mobile:
	cd mobile && flutter build appbundle --release
