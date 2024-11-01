GOLANGCI_LINT_CACHE?=/tmp/praktikum-golangci-lint-cache

.PHONY: build server-build server server-run server-build agent agent-run agent-build gofmt

build: server-build agent-build

server: server-run

server-run: server-build
	./cmd/server/server

server-build:
	go build -o ./cmd/server/server ./cmd/server/

agent: agent-run

agent-run: agent-build
	./cmd/agent/agent

agent-build:
	go build -o ./cmd/agent/agent ./cmd/agent/

gofmt:
	gofmt -s -w ./

test-iter1:
	metricstest -test.v -test.run=^TestIteration1$$ -agent-binary-path=cmd/agent/agent -binary-path=cmd/server/server

test-iter2a:
	metricstest -test.v -test.run=^TestIteration2A$$ -agent-binary-path=cmd/agent/agent -binary-path=cmd/server/server

test-iter2b:
	metricstest -test.v -test.run=^TestIteration2B$$ -agent-binary-path=cmd/agent/agent -binary-path=cmd/server/server -source-path=/home/arefev/dev/study/golang/mtrcstore

test-iter3a:
	metricstest -test.v -test.run=^TestIteration3A$$ -agent-binary-path=cmd/agent/agent -binary-path=cmd/server/server -source-path=/home/arefev/dev/study/golang/mtrcstore

test-iter3b:
	metricstest -test.v -test.run=^TestIteration3B$$ -agent-binary-path=cmd/agent/agent -binary-path=cmd/server/server -source-path=/home/arefev/dev/study/golang/mtrcstore

test-iter4:
	metricstest -test.v -test.run=^TestIteration4$$ -agent-binary-path=cmd/agent/agent -binary-path=cmd/server/server -source-path=/home/arefev/dev/study/golang/mtrcstore -server-port=8080

test-iter5:
	metricstest -test.v -test.run=^TestIteration5$$ -agent-binary-path=cmd/agent/agent -binary-path=cmd/server/server -source-path=/home/arefev/dev/study/golang/mtrcstore -server-port=8080

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