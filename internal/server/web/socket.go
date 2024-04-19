package web

import (
	"context"
	"errors"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

// object - экземпляр gorilla web object updater
var object = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// hWEBSocket - обработка запросов на подключение
func (o *Server) hWEBSocket(w http.ResponseWriter, r *http.Request) {
	// сообщаем о начале работы
	// log.Print("open websocket connection")
	// запускаем экземпляр
	socket, err := object.Upgrade(w, r, nil)
	// в случае ошибки
	if err != nil {
		// забираем ошибку
		var handshakeError websocket.HandshakeError
		// если ошибка не позволит работь
		if !errors.As(err, &handshakeError) {
			// сообщаем
			log.Printf("websocket fault [%s]", err)
			// уходим
			return
		}
	}
	// создаем контекст
	ctx, cancel := context.WithCancel(o.context)
	// обработчик закрытия соединения
	socket.SetCloseHandler(
		func(code int, text string) error {
			// закрываем контекст
			cancel()
			// выводим сообщение
			// log.Printf("close websocket connection")
			// уходим
			return nil
		})
	// если соединение установлено запускаем рутины и выходим
	if socket != nil {
		go func() { wsWrite(ctx, socket, o.core) }()
		go func() { wsRead(ctx, socket) }()
		return
	}
	// при попытке считать ресурс выводим сообщение
	io.WriteString(w, "=================================\nYou can't get websocket.json directly\n=================================")
}

// wsRead - доступ на чтение
func wsRead(ctx context.Context, ws *websocket.Conn) {
	// закроем при выходе
	defer ws.Close()
	// установим ограничения по кол-ву и времени
	ws.SetReadLimit(1024)
	ws.SetReadDeadline(time.Now().Add(time.Second * 3))
	// установим ответчика ping
	ws.SetPongHandler(
		func(string) error {
			ws.SetReadDeadline(time.Now().Add(time.Second * 3))
			return nil
		})
	// бесконечный цикл
	for {
		// ждем сигналы
		select {
		// контекст закрывается, выходим
		case <-ctx.Done():
			return
		// читаем
		default:
			_, p, err := ws.ReadMessage()
			if err != nil {
				break
			}
			log.Printf("=> [%v]", p)
		}
	}
}

// wsWrite - пишем данные
func wsWrite(ctx context.Context, ws *websocket.Conn, data *[]byte) {
	// создаем ограничители
	pingTicker := time.NewTicker(time.Second * 1)
	dataTicker := time.NewTicker(time.Millisecond * 250)
	// закроем всех при выходе
	defer func() {
		dataTicker.Stop()
		pingTicker.Stop()
		ws.Close()
	}()
	// бесконечный цикл
	for {
		// ждем сигналы
		select {
		// контекст закрывается, выходим
		case <-ctx.Done():
			return
		// пора отправлять ping
		case <-pingTicker.C:
			ws.SetWriteDeadline(time.Now().Add(time.Second * 5))
			if err := ws.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		// пора обновлять данные
		case <-dataTicker.C:
			ws.SetWriteDeadline(time.Now().Add(time.Second * 5))
			// отправим одним блоком все данные ядра
			if err := ws.WriteMessage(websocket.BinaryMessage, *data); err != nil {
				return
			}
		}
	}
}
