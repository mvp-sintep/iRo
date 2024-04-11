# Приложение iRo

### Задача

Реализация тестового стенда MODBUS RTU -> MODBUS TCP

### Структура

```
├── cmd
│    └── iRo
│          └──  main.go
├── config
│    └── system.yaml
├── internal
│    ├── application
│    |     └── run.go
│    ├── config
│    |     └── system.go
├── script
├── vendor
└── web
     ├── files
     |     ├── engine.js
     |     └── *.ico, *.png
     └── templates
           ├── error.html
           └── index.html
```



│    ├── core
│    |     └── data.go
│    ├── driver
│    |     ├── db
│    |     |     └── pool.go
│    |     ├── rtu
│    |     |     └── port.go
│    |     └── tcp
│    |           └── port.go
│    ├── pkg
│    |     ├── crc16
│    |     |     └── calculate.go
│    |     ├── get
│    |     |     └── get.go
│    |     ├── modbus
│    |     |     ├── errno.go
│    |     |     └── validate.go
│    |     └── set
│    |           └── set.go
│    ├── server
│    |     ├── http
│    |     |     ├── serve.go
│    |     |     ├── server.go
│    |     |     ├── socket.go
│    |     |     └── sql.go
│    |     ├── modbus
│    |     |     ├── rtu.go
│    |     |     |     ├── serve.go
│    |     |     |     └── server.go
│    |     |     └── tcp.go
│    |     |           ├── serve.go
│    |     |           └── server.go
│    |     └── ua
│    |           ├── export.c
│    |           ├── export.h
│    |           ├── open62541.c
│    |           ├── open62541.h
│    |           └── server.go
│    └── service
│          ├── ae
│          |     └── service.go
│          └── trend
│                └── service.go



