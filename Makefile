build: server-build agent-build

server: server-build server-run

server-build:
	go build -o ./cmd/server/server ./cmd/server/

server-run:
	./cmd/server/server

agent: agent-build agent-run

agent-build:
	go build -o ./cmd/agent/agent ./cmd/agent/

agent-run:
	./cmd/agent/agent

test-iter4:
	metricstest -test.v -test.run=^TestIteration4$$ -agent-binary-path=cmd/agent/agent -binary-path=cmd/server/server -source-path=/home/arefev/dev/study/golang/mtrcstore -server-port=8080