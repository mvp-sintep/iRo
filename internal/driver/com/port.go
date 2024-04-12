package com

import (
	"context"
	"errors"
	"iRo/internal/config"
	"log"
	"time"

	"github.com/goburrow/serial"
)

// Driver - блок данных RTU драйвера
type Driver struct {
	context  context.Context       // контекст для драйвера
	cancel   context.CancelFunc    // функция закрытия контекста драйвера
	cfg      *config.COMPortConfig // запись конфигурации
	rxBuffer [256]byte             // приемный буфер
	rxCount  int                   // кол-во байт в приемном буфере
	txBuffer [256]byte             // буфер отправки
	txCount  int                   // кол-во байт в буфере отправки
	SigDSR   chan struct{}         // сигнал, принимаемые данные готовы для обработки (data set ready)
	SigRTS   chan struct{}         // сигнал, отправляемые данные готовы для обработки (ready to send)
}

// New - создание блока данных драйвера
func New(ctx context.Context, cfg *config.COMPortConfig) (*Driver, error) {
	// проверяем аргументы
	if cfg == nil {
		return nil, errors.New("нет данных конфигурации RTU драйвера")
	}
	// проверим настройку паузы
	if cfg.Pause < 100 || cfg.Pause > 55000 {
		cfg.Pause = 5000
	}
	// создаем блок данных
	driver := &Driver{
		cfg:      cfg,
		rxBuffer: [256]byte{},
		rxCount:  0,
		txBuffer: [256]byte{},
		txCount:  0,
		SigDSR:   make(chan struct{}),
		SigRTS:   make(chan struct{}),
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
	// откроем порт
	port, err := serial.Open(&serial.Config{
		Address:  o.cfg.File,
		BaudRate: o.cfg.BaudRate,
		DataBits: o.cfg.DataBits,
		StopBits: o.cfg.StopBits,
		Parity:   o.cfg.Parity,
		Timeout:  time.Duration(o.cfg.Timeout) * time.Millisecond,
	})
	if err != nil {
		// вернем ошибку
		return err
	}
	defer func() {
		// если порт был создан
		if port != nil {
			// закроем порт
			port.Close()
		}
	}()
	// переменные
	timeout := time.Duration(o.cfg.Pause) * time.Microsecond
	buffer := make([]byte, 512)
	count := int(0)
	// бесконечный цикл
	for {
		// селектор сигналов
		select {
		// контекст закрывается
		case <-o.context.Done():
			return nil
		// сработал таймаут ожидания минимальной порции данных
		case <-time.After(timeout):
			// если идет отправка данных, а не прием
			if o.txCount > 0 {
				// ждем начала приема
				continue
			}
			// считываем данные в локальный буфер
			x, err := port.Read(buffer)
			// если ошибка
			if err != nil {
				// если ошибка отсутствия данных
				if err.Error() == "serial: timeout" {
					// если ничего не получили в последний раз, но есть данные в буфере
					if count > 0 {
						// указываем кол-во данных
						o.rxCount = count
						o.txCount = 0
						// готовим новый прием
						count = 0
						// выводим сообщение
						// log.Printf("dsr<-[%d]<-[%x]", o.rxCount, o.rxBuffer[:o.rxCount])
						// сигнал о наличии принятых данных
						o.SigDSR <- struct{}{}
					}
				} else {
					// другая ошибка, выводим сообщение
					log.Printf("com<-[error]<-[%v]", err)
				}
				// продолжаем
				continue
			}
			// есть данные
			if x > 0 {
				// данных не слишком много
				if (x + count) <= len(o.rxBuffer) {
					// копируем
					for i := 0; i < x && count < len(o.rxBuffer); i++ {
						o.rxBuffer[count] = buffer[i]
						count += 1
					}
				} else {
					// данных слишком много отбрасываем
					count = 0
				}
			}
		// данные готовы к отправке
		case <-o.SigRTS:
			// log.Printf("rts->[%d]->[%x]", o.txCount, o.txBuffer[:o.txCount])
			// если кол-во данных корректное
			if o.txCount > 0 && o.txCount < len(o.txBuffer) {
				// записываем в порт
				n, err := port.Write(o.txBuffer[:o.txCount])
				// при наличии ошибки
				if err != nil {
					// выводим сообщение
					log.Printf("found serial error [%d]<-[%v]", n, err)
				}
			}
			// отправили
			o.txCount = 0
		}
	}
}

// Shutdown - завершение работы
func (o *Driver) Shutdown() error {
	// очистим
	o.rxCount = 0
	o.txCount = 0
	// закроем контекст драйвера
	defer o.cancel()
	// уходим
	return nil
}
