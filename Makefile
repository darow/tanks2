target:
	make run -j2

run: run_server run_client

run_server:
	go run ./cmd -mode=server

run_client:
	@sleep 1
	go run ./cmd -mode=client
