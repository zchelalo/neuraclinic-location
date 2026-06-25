ifneq ("$(wildcard .env)", "")
	include .env
	export $(shell sed 's/=.*//' .env)
endif

DOCKER_COMPOSE_FILE = ./.docker/compose.yml
DOCKER_NETWORK = neuraclinic-network
SQLC_IMAGE = sqlc/sqlc:1.31.1
URI_DB = postgresql://$(DB_USER):$(DB_PASS)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)
MIGRATE = docker run --rm -v $(shell pwd)/internal/db/migrations:/migrations --network $(DOCKER_NETWORK) migrate/migrate -path /migrations -database "$(URI_DB)" -verbose
LOCAL_PROTO_CONTRACTS = ../neuraclinic-proto-contracts
LOCATION_DATA_DIR ?= ./data/sources
SEPOMEX_FILE ?= $(LOCATION_DATA_DIR)/sepomex/CPdescarga.txt
INEGI_ENTITIES_FILE ?= $(LOCATION_DATA_DIR)/inegi-ageeml/catun_entidad.zip
INEGI_MUNICIPALITIES_FILE ?= $(LOCATION_DATA_DIR)/inegi-ageeml/catun_municipio.zip
SOURCE_VERSION ?= 2026-06-22
SOURCE_LICENSE ?= Datos abiertos; validar terminos de Correos de Mexico / SEPOMEX antes de redistribuir
SOURCE_ATTRIBUTION ?= Correos de Mexico / SEPOMEX
SOURCE_URL ?= https://www.correosdemexico.gob.mx/SSLServicios/ConsultaCP/CodigoPostal_Exportar.aspx
SOURCE_ENCODING ?= latin1
INEGI_SOURCE_VERSION ?= 2026-06-12
INEGI_SOURCE_LICENSE ?= Terminos de libre uso INEGI; validar atribucion antes de redistribuir
INEGI_SOURCE_ATTRIBUTION ?= INEGI
INEGI_SOURCE_URL ?= https://www.inegi.org.mx/app/ageeml/

setup:
	$(MAKE) create-envs
	$(MAKE) tls-generate-dev
	$(MAKE) create-network
	$(MAKE) compose-build-detached

create-envs:
	test -f .env || cp .env.example .env

tls-generate-dev:
	./scripts/generate-dev-tls-certs.sh

tls-generate-dev-if-missing:
	test -f certs/public_key.pem -a -f certs/private_key.pem || ./scripts/generate-dev-tls-certs.sh

create-network:
	docker network inspect $(DOCKER_NETWORK) >/dev/null 2>&1 || docker network create $(DOCKER_NETWORK)

proto:
ifneq ("$(wildcard $(LOCAL_PROTO_CONTRACTS)/buf.yaml)", "")
	cd $(LOCAL_PROTO_CONTRACTS) && buf generate \
		--template ../neuraclinic-location/buf.gen.yaml \
		--output ../neuraclinic-location \
		--path proto/location/v1/location.proto
else
	buf generate buf.build/zchelalo-labs/neuraclinic-proto-contracts \
		--template buf.gen.yaml \
		--path location/v1/location.proto
endif

migrate-up:
	$(MIGRATE) up

migrate-down:
	$(MIGRATE) down

compose:
	$(MAKE) create-network
	docker compose -f $(DOCKER_COMPOSE_FILE) up

compose-detached:
	$(MAKE) create-network
	docker compose -f $(DOCKER_COMPOSE_FILE) up -d

compose-build:
	$(MAKE) create-network
	docker compose -f $(DOCKER_COMPOSE_FILE) up --build

compose-build-detached:
	$(MAKE) create-network
	docker compose -f $(DOCKER_COMPOSE_FILE) up --build -d

compose-down:
	docker compose -f $(DOCKER_COMPOSE_FILE) down

download-sepomex:
	./scripts/download-sepomex.sh "$(LOCATION_DATA_DIR)/sepomex"

download-inegi-ageeml:
	./scripts/download-inegi-ageeml.sh "$(LOCATION_DATA_DIR)/inegi-ageeml"

download-location-data: download-sepomex download-inegi-ageeml

import-sepomex:
	@test -n "$(SEPOMEX_FILE)" || (echo "Usage: make import-sepomex SEPOMEX_FILE=/path/to/CPdescarga.txt SOURCE_VERSION=YYYY-MM" >&2; exit 1)
	@test -f "$(SEPOMEX_FILE)" || (echo "SEPOMEX_FILE does not exist: $(SEPOMEX_FILE)" >&2; exit 1)
	@test -n "$(SOURCE_VERSION)" || (echo "SOURCE_VERSION is required, for example SOURCE_VERSION=2026-06" >&2; exit 1)
	$(MAKE) create-envs
	$(MAKE) tls-generate-dev-if-missing
	$(MAKE) create-network
	docker compose -f $(DOCKER_COMPOSE_FILE) up -d --build neuraclinic-location
	docker compose -f $(DOCKER_COMPOSE_FILE) run --rm --no-deps \
		-v "$(abspath $(SEPOMEX_FILE)):/tmp/sepomex-snapshot:ro" \
		neuraclinic-location location-import sepomex \
		--file /tmp/sepomex-snapshot \
		--source-version "$(SOURCE_VERSION)" \
		--source-url "$(SOURCE_URL)" \
		--license "$(SOURCE_LICENSE)" \
		--attribution "$(SOURCE_ATTRIBUTION)" \
		--encoding "$(SOURCE_ENCODING)"

import-inegi-ageeml:
	@test -n "$(INEGI_ENTITIES_FILE)" || (echo "Usage: make import-inegi-ageeml INEGI_ENTITIES_FILE=/path/to/catun_entidad.zip INEGI_MUNICIPALITIES_FILE=/path/to/catun_municipio.zip INEGI_SOURCE_VERSION=YYYY-MM-DD" >&2; exit 1)
	@test -n "$(INEGI_MUNICIPALITIES_FILE)" || (echo "Usage: make import-inegi-ageeml INEGI_ENTITIES_FILE=/path/to/catun_entidad.zip INEGI_MUNICIPALITIES_FILE=/path/to/catun_municipio.zip INEGI_SOURCE_VERSION=YYYY-MM-DD" >&2; exit 1)
	@test -f "$(INEGI_ENTITIES_FILE)" || (echo "INEGI_ENTITIES_FILE does not exist: $(INEGI_ENTITIES_FILE)" >&2; exit 1)
	@test -f "$(INEGI_MUNICIPALITIES_FILE)" || (echo "INEGI_MUNICIPALITIES_FILE does not exist: $(INEGI_MUNICIPALITIES_FILE)" >&2; exit 1)
	@test -n "$(INEGI_SOURCE_VERSION)" || (echo "INEGI_SOURCE_VERSION is required, for example INEGI_SOURCE_VERSION=2026-06-12" >&2; exit 1)
	$(MAKE) create-envs
	$(MAKE) tls-generate-dev-if-missing
	$(MAKE) create-network
	docker compose -f $(DOCKER_COMPOSE_FILE) up -d --build neuraclinic-location
	docker compose -f $(DOCKER_COMPOSE_FILE) run --rm --no-deps \
		-v "$(abspath $(INEGI_ENTITIES_FILE)):/tmp/inegi-entidades.zip:ro" \
		-v "$(abspath $(INEGI_MUNICIPALITIES_FILE)):/tmp/inegi-municipios.zip:ro" \
		neuraclinic-location location-import inegi-ageeml \
		--entities-file /tmp/inegi-entidades.zip \
		--municipalities-file /tmp/inegi-municipios.zip \
		--source-version "$(INEGI_SOURCE_VERSION)" \
		--source-url "$(INEGI_SOURCE_URL)" \
		--license "$(INEGI_SOURCE_LICENSE)" \
		--attribution "$(INEGI_SOURCE_ATTRIBUTION)"

import-location-data:
	$(MAKE) import-sepomex
	$(MAKE) import-inegi-ageeml

fmt:
	go fmt ./...

lint:
	go vet ./...

test:
	go test ./...

coverage:
	go test ./... -coverprofile=coverage.out

build:
	mkdir -p dist
	go build -buildvcs=false -trimpath -o dist/neuraclinic-location ./cmd
	go build -buildvcs=false -trimpath -o dist/location-import ./cmd/location-import

sqlc:
	docker run --rm --user $(shell id -u):$(shell id -g) -v $(shell pwd):/src -w /src $(SQLC_IMAGE) generate

.PHONY: setup create-envs tls-generate-dev tls-generate-dev-if-missing create-network proto migrate-up migrate-down compose compose-detached compose-build compose-build-detached compose-down download-sepomex download-inegi-ageeml download-location-data import-sepomex import-inegi-ageeml import-location-data fmt lint test coverage build sqlc
