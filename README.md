# tarantool-rest-api

REST API доступен по адресу http://35.246.166.66/kv

- POST /kv body: {"key": 1, "value": {"name": "Name", "secondName": "secondName"}} 
- PUT kv/{id} body: {"value": {"name": "Name", "secondName": "secondName"}} 
- GET kv/{id} 
- DELETE kv/{id} 
- POST возвращает 409 если ключ уже существует, 
- POST, PUT возвращают 400 если боди некорректное 
- PUT, GET, DELETE возвращает 404 если такого ключа нет - все операции логируются

формат спейса : {id unsigned, map{name, secondName}}

логфайл хранится на сервере, ![logfile](./logfile.jpg)

```bash
tarantool
box.cfg{listen = 3311}
```

Запускает веб-сервер на 8080 порту на всех сетевых интерфейсах
```bash
go build main.go utils.go structs.go api.go
./main
```

api.go содержит обработку всех запросов