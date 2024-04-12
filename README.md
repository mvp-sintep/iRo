# Приложение iRo

### Задача

Реализация тестового стенда

...
MODBUS RTU ── ядро данных ── MODBUS RTU
                   ├──────── WEB SOCKET ── HTML
MODBUS TCP ────────┼──────── MODBUS TCP
    OPC UA ────────┴──────── OPC UA
...

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
|    |     └── com
│    |           ├── api.go
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
│          └── mbrtu
│                ├── serve.go
│                └── server.go
├── vendor
└── web
     ├── files
     |     ├── engine.js
     |     └── *.ico, *.png
     └── templates
           ├── error.html
           └── index.html
```
