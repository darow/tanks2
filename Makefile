target:
	make run -j2

local:
	go run ./cmd

run: run3

run2:
	make server2 client1 -j2

run3:
	make server3 client1 client2 -j3

server2:
	go run ./cmd -mode=server -players_count=2 -debug

server3:
	go run ./cmd -mode=server -players_count=3 -debug

client1:
	@sleep 1
	go run ./cmd -mode=client -address="127.0.0.1:8080" -player_id=1 -debug

client2:
	@sleep 1
	go run ./cmd -mode=client -address="127.0.0.1:8080" -player_id=2 -debug
