package mbrtu

import (
	"context"
	"errors"
	"iRo/internal/config"
	"iRo/internal/driver/com"
)

// Server - данные MODBUS RTU сервера
type Server struct {
	context context.Context         // контекст для сервера
	cancel  context.CancelFunc      // функция закрытия контекста сервера
	cfg     *config.ModbusRTUConfig // запись конфигурации
	core    *[]byte                 // данные ядра
	nda     *int                    // признак изменения данных
	comDrv  *com.Driver             // драйвер
}

// New - создание блока данных сервера
func New(ctx context.Context, core *[]byte, nda *int, comDrv *com.Driver, cfg *config.ModbusRTUConfig) (*Server, error) {
	// проверяем аргументы
	if cfg == nil {
		return nil, errors.New("нет данных конфигурации MODBUS RTU сервера")
	}
	// создаем блок данных
	server := &Server{
		cfg:    cfg,
		core:   core,
		nda:    nda,
		comDrv: comDrv,
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
