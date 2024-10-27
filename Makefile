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