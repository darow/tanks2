target:
	make run -j2

local:
	go run ./cmd

run: run_server run_client

run_server:
	go run ./cmd -mode=server -debug

run_client:
	@sleep 1
	go run ./cmd -mode=client -debug
