target:
	make run -j2

local:
	go run ./cmd

run: server client

run2:
	make server2 client2 -j2

run3:
	make server3 client3-player1 client3-player2 -j3

server:
	go run ./cmd -mode=server -debug

client:
	@sleep 1
	go run ./cmd -mode=client -debug

server2:
	go run ./cmd -mode=server -players_count=2 -debug

client2:
	@sleep 1
	go run ./cmd -mode=client -address="127.0.0.1:8080" -player_id=1 -debug

server3:
	go run ./cmd -mode=server -players_count=3 -debug

client3-player1:
	@sleep 1
	go run ./cmd -mode=client -address="127.0.0.1:8080" -player_id=1 -debug

client3-player2:
	@sleep 1
	go run ./cmd -mode=client -address="127.0.0.1:8080" -player_id=2 -debug
