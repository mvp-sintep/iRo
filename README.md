# Приложение iRo

### Задача

Реализация тестового стенда

MODBUS RTU + ядро данных + WEB SOCKET + HTML + MODBUS TCP  + OPC UA

### Структура проекта

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
│    ├── core
│    |     └── data.go
│    ├── daemon
│    |     └── daemon.go
│    ├── driver
|    |     ├── com
│    |     |     ├── api.go
│    |     |     └── port.go
|    |     └── tcp
│    |           ├── api.go
│    |           ├── connection.go
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
│    └── server
│          ├── mbrtu
│          |     ├── serve.go
│          |     └── server.go
│          ├── mbtcp
│          |     ├── serve.go
│          |     └── server.go
│          ├── ua
│          |     ├── export.c
│          |     ├── export.h
│          |     ├── open62541.*
│          |     └── server.go
│          └── web
│                ├── serve.go
│                ├── server.go
│                └── socket.go
├── vendor
└── web
     ├── files
     |     ├── engine.js
     |     └── *.ico, *.png
     └── templates
           ├── error.html
           └── index.html
```
