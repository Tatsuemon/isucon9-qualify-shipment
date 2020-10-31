# isucon9-qualify-shipment
isucon9-qualifyのshipment serviceのようなものを作ってみた

ここでは配達サービスを表現

[役割]
1. 配達の生成
2. 配達の問合せ(配達IDからのみ) -> 本体に配達IDを持たせる
3. 配達statusを変化(initial -> wait_pickup -> shipping -> done, cancel)

[参考.go](https://github.com/isucon/isucon9-qualify/blob/1409a5ca6883f343e024a72fb1fa6227fa57b293/bench/server/shipment.go)
[参考.md](https://github.com/isucon/isucon9-qualify/blob/master/webapp/docs/EXTERNAL_SERVICE_SPEC.md)



### DB migrate

#### migrationファイル作成
```
$ docker-compose run app goose create create_{table名}_table sql
```

#### migration
```
$ docker-compose run app goose up
```

#### rollback
```
$ docker-compose run app goose down
```

### TODO
- testコード
- github actionsでテストを走らせる
