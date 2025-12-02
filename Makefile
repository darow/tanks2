target:
	make run -j2

local:
	go run ./cmd

run: server client

server:
	go run ./cmd -mode=server -debug

client:
	@sleep 1
	go run ./cmd -mode=client -debug
