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