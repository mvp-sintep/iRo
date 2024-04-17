package tcp

import (
	"context"
	"iRo/internal/config"
	"net"
	"time"
)

// Connection - блок данных соединения
type Connection struct {
	context  context.Context         // контекст для соединения
	cancel   context.CancelFunc      // функция закрытия контекста соединения
	cfg      *config.ModbusTCPConfig // запись конфигурации
	conn     net.Conn                // соединение
	rxBuffer [259]byte               // приемный буфер
	rxCount  int                     // кол-во байт в приемном буфере
	txBuffer [259]byte               // буфер отправки
	txCount  int                     // кол-во байт в буфере отправки
}

// NewConnection - создание блока данных соединения
func NewConnection(ctx context.Context, cfg *config.ModbusTCPConfig, conn net.Conn) (*Connection, error) {
	// создаем блок данных соединения
	connection := &Connection{
		cfg:      cfg,
		conn:     conn,
		rxBuffer: [259]byte{},
		rxCount:  0,
		txBuffer: [259]byte{},
		txCount:  0,
	}
	// если не задан контекст
	if ctx == nil {
		ctx = context.Background()
	}
	// задаем значения
	connection.context, connection.cancel = context.WithCancel(ctx)
	// соединение создано
	return connection, nil
}

// Serve - обработка соединения
func (o *Connection) Serve(dsr chan *Connection, rts chan *Connection) error {
	// создаем обработчик для отправки данных
	go func() {
		// бесконечный цикл
		for {
			select {
			// контекст соединения закрывается
			case <-o.context.Done():
				// уходим
				return
			// есть данные для отправки
			case <-rts:
				// выводим сообщение
				//log.Printf("[%d]<-[%x]", o.txCount, o.txBuffer[:o.txCount])
				// настраиваем таймаут
				if err := o.conn.SetWriteDeadline(time.Now().Add(time.Duration(o.cfg.Write) * time.Second)); err != nil {
					return
				}
				// проверяем данные
				if o.conn != nil && o.rxCount > 7 && o.rxCount < len(o.txBuffer) {
					// отправляем
					o.conn.Write(o.txBuffer[:o.txCount])
				}
			}
		}
	}()
	// закрываем соединение при выходе
	defer func() {
		o.conn.Close()
	}()
	// таймаут контроля
	timeout := time.Duration(o.cfg.Control) * time.Millisecond
	// бесконечный цикл
	for {
		// селектор сигналов
		select {
		// контекст закрывается
		case <-o.context.Done():
			// уходим
			return nil
		// нет событий
		case <-time.After(timeout):
			// устанавливаем время ожидания
			if err := o.conn.SetReadDeadline(time.Now().Add(time.Duration(o.cfg.Read) * time.Second)); err != nil {
				return err
			}
			// читаем входящий поток
			if count, err := o.conn.Read(o.rxBuffer[:]); err == nil && count > 11 && count < len(o.rxBuffer) {
				// кол-во данных
				o.rxCount = count
				// выводим сообщение
				//log.Printf("[%d]<-[%x]", o.rxCount, o.rxBuffer[:o.rxCount])
				// сообщаем о наличии данных
				dsr <- o
				// сохраняем соединение
				continue
			}
			// если соединение было закрыто или нет данных выходим
			return nil
		}
	}
}
