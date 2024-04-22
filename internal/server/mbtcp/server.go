package mbtcp

import (
	"context"
	"errors"
	"iRo/internal/config"
	"iRo/internal/driver/tcp"
)

// Server - данные MODBUS TCP сервера
type Server struct {
	context context.Context         // контекст для сервера
	cancel  context.CancelFunc      // функция закрытия контекста сервера
	cfg     *config.ModbusTCPConfig // запись конфигурации
	core    *[]byte                 // данные ядра
	nda     *int                    // признак изменения данных
	tcpDrv  *tcp.Driver             // драйвер
}

// New - создание блока данных сервера
func New(ctx context.Context, core *[]byte, nda *int, tcpDrv *tcp.Driver, cfg *config.ModbusTCPConfig) (*Server, error) {
	// проверяем аргументы
	if cfg == nil {
		return nil, errors.New("нет данных конфигурации MODBUS TCP сервера")
	}
	// создаем блок данных
	server := &Server{
		cfg:    cfg,
		core:   core,
		nda:    nda,
		tcpDrv: tcpDrv,
	}
	// если не задан контекст
	if ctx == nil {
		ctx = context.Background()
	}
	// задаем значения
	server.context, server.cancel = context.WithCancel(ctx)
	// сервер создан
	return server, nil
}

// Run - запуск сервера
func (o *Server) Run() error {
	go func() { o.serve() }()
	return nil
}

// Shutdown - завершение работы
func (o *Server) Shutdown() error {
	// закроем контекст сервера
	defer o.cancel()
	// уходим
	return nil
}
