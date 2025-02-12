include .env

GOLANGCI_LINT_CACHE?=/tmp/praktikum-golangci-lint-cache
T_AGENT_BINARY_PATH=cmd/agent/agent
T_BINARY_PATH=cmd/server/server
T_SOURCE_PATH=/home/arefev/dev/study/golang/mtrcstore
T_SERVER_PORT=8080
FILE_STORAGE_PATH=./storage.json
USER=CURRENT_UID=$$(id -u):0
DOCKER_PROJECT_NAME=mtrcstore
DATABASE_DSN="host=${DB_HOST} user=${DB_USER} password=${DB_PASSWORD} dbname=${DB_NAME} sslmode=disable"


.PHONY: build server-build server server-run server-build agent agent-run agent-build gofmt test-iter

build: server-build agent-build

server: server-run

server-run: server-build
	./cmd/server/server -d=${DATABASE_DSN} -k="${SECRET_KEY}" -a="localhost:8081"

server-build:
	go build -o ./cmd/server/server ./cmd/server/

server-build-cover:
	go build -cover -o ./cmd/server/server ./cmd/server/
.PHONY: server-build-cover

agent: agent-run

agent-run: agent-build
	./cmd/agent/agent -r 2 -p 1 -k="${SECRET_KEY}"

agent-build:
	go build -o ./cmd/agent/agent ./cmd/agent/

gofmt:
	gofmt -s -w ./

containers:
	$(USER) docker-compose --project-name $(DOCKER_PROJECT_NAME) up -d

test: server-build-cover
	go test ./... -cover -coverprofile=coverage.out && \
	go tool cover -html coverage.out -o test.html && \
	go tool cover -func=coverage.out
.PHONY: test

test-clear: 
	rm -f coverage.out && rm -f test.html
.PHONY: test-clear

test-iter: test-iter1 test-iter2a test-iter2b test-iter3a test-iter3b test-iter4 test-iter5 test-iter6 test-iter7 test-iter8 test-iter9 test-iter10 test-iter11 test-iter12 test-iter13 test-iter14

test-iter1:
	metricstest -test.v -test.run=^TestIteration1$$ -agent-binary-path=${T_AGENT_BINARY_PATH} -binary-path=${T_BINARY_PATH} -source-path=${T_SOURCE_PATH} -server-port=${T_SERVER_PORT} -file-storage-path=${FILE_STORAGE_PATH}

test-iter2a:
	metricstest -test.v -test.run=^TestIteration2A$$ -agent-binary-path=${T_AGENT_BINARY_PATH} -binary-path=${T_BINARY_PATH} -source-path=${T_SOURCE_PATH} -server-port=${T_SERVER_PORT} -file-storage-path=${FILE_STORAGE_PATH}

test-iter2b:
	metricstest -test.v -test.run=^TestIteration2B$$ -agent-binary-path=${T_AGENT_BINARY_PATH} -binary-path=${T_BINARY_PATH} -source-path=${T_SOURCE_PATH} -server-port=${T_SERVER_PORT} -file-storage-path=${FILE_STORAGE_PATH}

test-iter3a:
	metricstest -test.v -test.run=^TestIteration3A$$ -agent-binary-path=${T_AGENT_BINARY_PATH} -binary-path=${T_BINARY_PATH} -source-path=${T_SOURCE_PATH} -server-port=${T_SERVER_PORT} -file-storage-path=${FILE_STORAGE_PATH}

test-iter3b:
	metricstest -test.v -test.run=^TestIteration3B$$ -agent-binary-path=${T_AGENT_BINARY_PATH} -binary-path=${T_BINARY_PATH} -source-path=${T_SOURCE_PATH} -server-port=${T_SERVER_PORT} -file-storage-path=${FILE_STORAGE_PATH}

test-iter4:
	metricstest -test.v -test.run=^TestIteration4$$ -agent-binary-path=${T_AGENT_BINARY_PATH} -binary-path=${T_BINARY_PATH} -source-path=${T_SOURCE_PATH} -server-port=${T_SERVER_PORT} -file-storage-path=${FILE_STORAGE_PATH}

test-iter5:
	metricstest -test.v -test.run=^TestIteration5$$ -agent-binary-path=${T_AGENT_BINARY_PATH} -binary-path=${T_BINARY_PATH} -source-path=${T_SOURCE_PATH} -server-port=${T_SERVER_PORT} -file-storage-path=${FILE_STORAGE_PATH}

test-iter6:
	metricstest -test.v -test.run=^TestIteration6$$ -agent-binary-path=${T_AGENT_BINARY_PATH} -binary-path=${T_BINARY_PATH} -source-path=${T_SOURCE_PATH} -server-port=${T_SERVER_PORT} -file-storage-path=${FILE_STORAGE_PATH}

test-iter7:
	metricstest -test.v -test.run=^TestIteration7$$ -agent-binary-path=${T_AGENT_BINARY_PATH} -binary-path=${T_BINARY_PATH} -source-path=${T_SOURCE_PATH} -server-port=${T_SERVER_PORT} -file-storage-path=${FILE_STORAGE_PATH}

test-iter8:
	metricstest -test.v -test.run=^TestIteration8$$ -agent-binary-path=${T_AGENT_BINARY_PATH} -binary-path=${T_BINARY_PATH} -source-path=${T_SOURCE_PATH} -server-port=${T_SERVER_PORT} -file-storage-path=${FILE_STORAGE_PATH}

test-iter9:
	metricstest -test.v -test.run=^TestIteration9$$ -agent-binary-path=${T_AGENT_BINARY_PATH} -binary-path=${T_BINARY_PATH} -source-path=${T_SOURCE_PATH} -server-port=${T_SERVER_PORT} -file-storage-path=${FILE_STORAGE_PATH}

test-iter10:
	metricstest -test.v -test.run=^TestIteration10$$ -agent-binary-path=${T_AGENT_BINARY_PATH} -binary-path=${T_BINARY_PATH} -source-path=${T_SOURCE_PATH} -server-port=${T_SERVER_PORT} -file-storage-path=${FILE_STORAGE_PATH}

test-iter11:
	metricstest -test.v -test.run=^TestIteration11$$ -agent-binary-path=${T_AGENT_BINARY_PATH} -binary-path=${T_BINARY_PATH} -source-path=${T_SOURCE_PATH} -server-port=${T_SERVER_PORT} -file-storage-path=${FILE_STORAGE_PATH} -database-dsn=${DATABASE_DSN}

test-iter12:
	metricstest -test.v -test.run=^TestIteration12$$ -agent-binary-path=${T_AGENT_BINARY_PATH} -binary-path=${T_BINARY_PATH} -source-path=${T_SOURCE_PATH} -server-port=${T_SERVER_PORT} -file-storage-path=${FILE_STORAGE_PATH} -database-dsn=${DATABASE_DSN}

test-iter13:
	metricstest -test.v -test.run=^TestIteration13$$ -agent-binary-path=${T_AGENT_BINARY_PATH} -binary-path=${T_BINARY_PATH} -source-path=${T_SOURCE_PATH} -server-port=${T_SERVER_PORT} -file-storage-path=${FILE_STORAGE_PATH} -database-dsn=${DATABASE_DSN}

test-iter14:
	metricstest -test.v -test.run=^TestIteration14$$ -agent-binary-path=${T_AGENT_BINARY_PATH} -binary-path=${T_BINARY_PATH} -source-path=${T_SOURCE_PATH} -server-port=${T_SERVER_PORT} -file-storage-path=${FILE_STORAGE_PATH} -database-dsn=${DATABASE_DSN} -key="${SECRET_KEY}"

.PHONY: golangci-lint-run
golangci-lint-run: _golangci-lint-rm-unformatted-report

.PHONY: _golangci-lint-reports-mkdir
_golangci-lint-reports-mkdir:
	mkdir -p ./golangci-lint

.PHONY: _golangci-lint-run
_golangci-lint-run: _golangci-lint-reports-mkdir
	-docker run --rm \
    -v $(shell pwd):/app \
    -v $(GOLANGCI_LINT_CACHE):/root/.cache \
    -w /app \
    golangci/golangci-lint:v1.57.2 \
        golangci-lint run \
            -c .golangci.yml \
	> ./golangci-lint/report-unformatted.json

.PHONY: _golangci-lint-format-report
_golangci-lint-format-report: _golangci-lint-run
	cat ./golangci-lint/report-unformatted.json | jq > ./golangci-lint/report.json

.PHONY: _golangci-lint-rm-unformatted-report
_golangci-lint-rm-unformatted-report: _golangci-lint-format-report
	rm ./golangci-lint/report-unformatted.json

.PHONY: golangci-lint-clean
golangci-lint-clean:
	sudo rm -rf ./golangci-lint 