package web

import (
	"context"
	"errors"
	"fmt"
	"iRo/internal/config"
	"net/http"
	"time"
)

// Server - данные http сервера
type Server struct {
	context context.Context    // контекст для сервера
	cancel  context.CancelFunc // функция закрытия контекста сервера
	cfg     *config.HTTPConfig // запись конфигурации
	eCh     chan<- error       // канал ошибок
	core    *[]byte            // данные ядра
	server  *http.Server       // сам сервер
}

// New - создание блока данных сервера
func New(ctx context.Context, eCh chan<- error, core *[]byte, cfg *config.HTTPConfig) (*Server, error) {

	// проверяем аргументы
	if cfg == nil {
		return nil, errors.New("нет данных конфигурации HTTP сервера")
	}
	if cfg.Address == "" {
		return nil, errors.New("нет данных конфигурации сетевого адаптера HTTP сервера")
	}
	// создаем блок данных
	server := &Server{
		cfg:  cfg,
		eCh:  eCh,
		core: core,
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

	// создаем сервер
	o.server = &http.Server{
		Addr:         fmt.Sprintf("%s:%s", o.cfg.Address, o.cfg.Port),
		Handler:      o,
		ReadTimeout:  time.Second * time.Duration(o.cfg.Read),
		WriteTimeout: time.Second * time.Duration(o.cfg.Write),
	}
	// запускаем обработку запросов
	return o.server.ListenAndServe()
}

// Shutdown - завершение работы
func (o *Server) Shutdown() error {
	// закроем контекст сервера
	defer o.cancel()
	// создадим новый контекст с отсрочкой закрытия
	ctx, cancel := context.WithTimeout(o.context, time.Duration(o.cfg.Shutdown*int(time.Second)))
	// закроем при выходе
	defer cancel()
	// дадим команду на остановку
	o.server.Shutdown(ctx)
	// уходим
	return nil
}
