package tcp

import (
	"context"
	"errors"
	"fmt"
	"iRo/internal/config"
	"net"
	"time"
)

// Driver - блок данных RTU драйвера
type Driver struct {
	context context.Context         // контекст для драйвера
	cancel  context.CancelFunc      // функция закрытия контекста драйвера
	eCh     chan<- error            // канал ошибок
	cfg     *config.ModbusTCPConfig // запись конфигурации
	SigDSR  chan *Connection        // сигнал, принимаемые данные готовы для обработки (data set ready)
	SigRTS  chan *Connection        // сигнал, отправляемые данные готовы для обработки (ready to send)
}

// New - создание блока данных драйвера
func New(ctx context.Context, eCh chan<- error, cfg *config.ModbusTCPConfig) (*Driver, error) {
	// проверяем аргументы
	if cfg == nil {
		return nil, errors.New("нет данных конфигурации TCP драйвера")
	}
	// создаем блок данных драйвера
	driver := &Driver{
		eCh:    eCh,
		cfg:    cfg,
		SigDSR: make(chan *Connection),
		SigRTS: make(chan *Connection),
	}
	// если не задан контекст
	if ctx == nil {
		ctx = context.Background()
	}
	// задаем значения
	driver.context, driver.cancel = context.WithCancel(ctx)
	// драйвер создан
	return driver, nil
}

// Run - запуск драйвера
func (o *Driver) Run() error {

	// создаем слушателя
	listener, err := net.Listen("tcp4", fmt.Sprintf("%s:%d", o.cfg.Address, o.cfg.Port))
	// проверяем
	if err != nil {
		// вернем ошибку
		return err
	}
	// закроем при выходе
	defer listener.Close()

	// бесконечный цикл
	for {
		// селектор сигналов
		select {
		// контекст закрывается
		case <-o.context.Done():
			// уходим
			return nil
		// периодический контроль
		case <-time.After(time.Duration(o.cfg.Control) * time.Millisecond):
			// проверяем на соединение
			conn, err := listener.Accept()
			// проверяем на наличие ошибок
			if err != nil {
				// сигнал ошибки
				o.eCh <- err
				// продолжаем ждать
				continue
			}
			// создаем обработчика
			connection, err := NewConnection(o.context, o.cfg, conn)
			// проверяем на наличие ошибок
			if err != nil {
				// сигнал ошибки
				o.eCh <- err
				// продолжаем ждать
				continue
			}
			// обработка соединения
			go func() {
				if err := connection.Serve(o.SigDSR, o.SigRTS); err != nil {
					o.eCh <- err
				}
			}()
		}

	}
}

// Shutdown - завершение работы
func (o *Driver) Shutdown() error {
	// закроем контекст драйвера
	defer o.cancel()
	// уходим
	return nil
}
