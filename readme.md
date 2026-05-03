**quick start:**

host server:
```shell
go run ./cmd -mode=server
```

client:
```shell
go run ./cmd -mode=client -address="127.0.0.1:8080"
```

host server for 4 players:
```shell
go run ./cmd -mode=server -players_count=4
```

clients for 4 players:
```shell
go run ./cmd -mode=client -address="127.0.0.1:8080" -players_count=4 -player_id=1
go run ./cmd -mode=client -address="127.0.0.1:8080" -players_count=4 -player_id=2
go run ./cmd -mode=client -address="127.0.0.1:8080" -players_count=4 -player_id=3
```

глобальный план:
1. сделать соло танки
2. сделать мультиплеер танки
3. сделать платформер шутер
4. сделать мультиплеер баттл рояль в стиле Noita

бэклог:
1.пофиксить пролетание снаряда через две рядом стоящие стены
2.сделать отталкивание танка от стены при повороте, чтобы не было блокировок
