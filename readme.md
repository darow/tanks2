**quick start:**

host server for 2 players:
```shell
go run ./cmd -mode=server -players_count=2
```

client for 2 players:
```shell
go run ./cmd -mode=client -address="127.0.0.1:8080" -player_id=1
```

host server for 3 players:
```shell
go run ./cmd -mode=server -players_count=3
```

clients for 3 players:
```shell
go run ./cmd -mode=client -address="127.0.0.1:8080" -player_id=1
go run ./cmd -mode=client -address="127.0.0.1:8080" -player_id=2
```

make shortcuts:
```shell
make run2
make run3
```

глобальный план:
1. сделать соло танки
2. сделать мультиплеер танки
3. сделать платформер шутер
4. сделать мультиплеер баттл рояль в стиле Noita

бэклог:
1.пофиксить пролетание снаряда через две рядом стоящие стены
2.сделать отталкивание танка от стены при повороте, чтобы не было блокировок
