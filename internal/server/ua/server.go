package ua

// #include <stdlib.h>
// #include "export.h"
import "C"

import (
	"context"
	"errors"
	"fmt"
	"iRo/internal/config"
	"iRo/internal/core"
	"iRo/internal/pkg/get"
	"log"
	"unsafe"
)

var count int = 0

//export go_callback_info
func go_callback_info(x []byte) {
	count += 1
	if count == 100 {
		log.Printf("[%s]->", x)
		count = 0
	}
}

//export go_callback_int32
func go_callback_int32(address C.int16_t, value *C.int32_t) C.int16_t {
	if value == nil {
		return 1
	}
	ptr := (*int32)(unsafe.Pointer(value))
	if address < 0 || int(address) >= len(core.Data) {
		*ptr = 0
		return 1
	}
	*ptr = int32(get.Uint32(core.Data[address:]))
	return 1
}

//export go_callback_int16
func go_callback_int16(address C.int16_t, value *C.int16_t) C.int16_t {
	if value == nil {
		return 1
	}
	ptr := (*int16)(unsafe.Pointer(value))
	if address < 0 || int(address) >= len(core.Data) {
		*ptr = 0
		return 1
	}
	*ptr = int16(get.Uint16(core.Data[address:]))
	return 1
}

//export go_callback_float
func go_callback_float(address C.int16_t, value *C.int32_t) C.int16_t {
	if value == nil {
		return 1
	}
	ptr := (*float32)(unsafe.Pointer(value))
	if address < 0 || int(address) >= len(core.Data) {
		*ptr = 0
		return 1
	}
	*ptr = get.Float32(core.Data[address:])
	return 1
}

// Server - блок данных сервера
type Server struct {
	context context.Context    // контекст для сервера
	cancel  context.CancelFunc // функция закрытия контекста сервера
	cfg     *config.UAConfig   // запись конфигурации
	Server  *C.UA_Server
}

// New - создание блока данных сервера
func New(ctx context.Context, cfg *config.UAConfig) (*Server, error) {

	// проверяем аргументы
	if cfg == nil {
		return nil, errors.New("нет данных конфигурации UA сервера")
	}

	// создаем блок данных
	server := &Server{
		cfg:    cfg,
		Server: C.New(C.uint32_t(cfg.Port), C.uint32_t(0), C.int(cfg.Namespace)),
	}

	// если не задан контекст
	if ctx == nil {
		ctx = context.Background()
	}
	// задаем значения
	server.context, server.cancel = context.WithCancel(ctx)

	return server, nil
}

// CreateObjectNode - создание ноды объекта
func (o *Server) CreateObjectNode(nodeId string, nodeName string) {
	C.CreateObjectNode(o.Server, C.CString(nodeId), C.CString(nodeName))
}

// CreateI32DataSource - создание тега int32
func (o *Server) CreateI32DataSource(nodeID string, nodeName string, parentNodeID string) {
	C.CreateI32DataSource(o.Server, C.CString(nodeID), C.CString(nodeName), C.CString(parentNodeID), 0)
}

// CreateI16DataSource - создание тега int16
func (o *Server) CreateI16DataSource(nodeID string, nodeName string, parentNodeID string) {
	C.CreateI16DataSource(o.Server, C.CString(nodeID), C.CString(nodeName), C.CString(parentNodeID), 0)
}

// CreateFloatDataSource - создание тега float
func (o *Server) CreateFloatDataSource(nodeID string, nodeName string, parentNodeID string) {
	C.CreateFloatDataSource(o.Server, C.CString(nodeID), C.CString(nodeName), C.CString(parentNodeID), 0)
}

// Run - запуск сервера
func (o *Server) Run() error {
	// создаем ноду
	o.CreateObjectNode("tags", "tags")
	// создаем теги для доступа к ядру как к массиву int16
	for i := 0; i < int(core.Size/2); i++ {
		tmp := fmt.Sprintf("i16%03d", i)
		o.CreateI16DataSource(tmp, tmp, "tags")
	}
	// создаем теги для доступа к ядру как к массиву int32
	for i := 0; i < int(core.Size/4); i++ {
		tmp := fmt.Sprintf("i32%03d", i)
		o.CreateI32DataSource(tmp, tmp, "tags")
	}
	// создаем теги для доступа к ядру как к массиву float
	for i := 0; i < int(core.Size/4); i++ {
		tmp := fmt.Sprintf("f%03d", i)
		o.CreateFloatDataSource(tmp, tmp, "tags")
	}
	// запускаем
	C.Run(o.Server)
	// выходим
	return nil
}

// Shutdown - завершение работы
func (o *Server) Shutdown() error {
	// закроем контекст сервера
	defer o.cancel()
	// уходим
	return nil
}
