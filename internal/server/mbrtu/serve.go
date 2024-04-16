package mbrtu

import (
	"iRo/internal/core"
	"iRo/internal/pkg/crc16"
	"iRo/internal/pkg/get"
	"iRo/internal/pkg/modbus"
	"iRo/internal/pkg/set"
)

// serve - обработка запросов MODBUS RTU
func (o *Server) serve() {

	// рассчитываем заявленную длину блока данных
	length := o.cfg.Core.Start + o.cfg.Core.Bytes
	// если хотим слишком много
	if length > len(*o.core) {
		// задаем доступный максимум
		length = len(*o.core)
	}
	// если хотим слишком мало
	if length < 1 {
		// задаем минимум
		length = 1
	}
	// уменьшаем на 1 для удобства использования в циклах
	length -= 1

	// переменные
	buffer := make([]byte, 256)
	crc := crc16.Value{}

	for {
		select {
		// контекст закрывается
		case <-o.context.Done():
			return

		// сигнал от драйвера - пакет принят
		case <-o.comDrv.SigDSR:

			// забираем данные
			if n, request := o.comDrv.GetRXData(); n > 0 {

				// пустые пакеты и пакеты не для это ноды отбрасываем
				if n < 8 || get.Uint8(request[0:]) != uint8(o.cfg.Node) {
					continue
				}

				// проверяем запрос, формируем ответ если нет ошибок
				if !modbus.IsValidRTURequest(request) {
					set.Uint8(buffer, uint8(o.cfg.Node))
					set.Uint8(buffer[1:], get.Uint8(request[1:])|0x80)
					set.Uint8(buffer[2:], modbus.ErrorIllegalDataValue)
					set.Uint16Swapped(buffer[3:], crc.Reset().Push(buffer[:3]).Result())
					o.comDrv.SetTXData(buffer[:5])
					o.comDrv.SigRTS <- struct{}{}
					continue
				}

				// селектор по номеру modbus функции
				switch get.Uint8(request[1:]) {

				// чтение регистров хранения
				case 3:
					//*** DATA CORE => BUFFER

					// кол-во регистров
					quantity := get.Uint16(request[4:])
					// если больше 125
					if quantity > 125 {
						// отсекаем по пределу (1+1+1+125*2+2=255)
						quantity = 125
					}
					// смещение в ядре данных
					i := o.cfg.Core.Start + int(get.Uint16(request[2:])*2)
					// смещение в выходном буфере
					tx := 3
					// заполняем буфер
					for ; quantity > 0; quantity -= 1 {
						// заполнитель значений за пределами ядра данных
						tmp := uint16(0)
						if i < length {
							tmp = get.Uint16Swapped(core.Data[i:])
						}
						// пишем в буфер
						set.Uint16(buffer[tx:], tmp)
						// один регистр = 2 байта
						tx += 2
						i += 2
					}

					set.Uint8(buffer[0:], uint8(o.cfg.Node))

					// если есть данные по запросу
					if tx > 3 {
						set.Uint8(buffer[1:], 0x03)
						set.Uint8(buffer[2:], uint8(tx-3))
					} else {
						// если нет данных
						set.Uint8(buffer[1:], 0x03|0x80)
						set.Uint8(buffer[2:], modbus.ErrorIllegalDataAddress)
						tx = 3
					}

					// данные готовы, добавляем контрольную сумму, записываем и даем сигнал на отправку
					set.Uint16Swapped(buffer[tx:], crc.Reset().Push(buffer[:tx]).Result())
					o.comDrv.SetTXData(buffer[:tx+2])
					o.comDrv.SigRTS <- struct{}{}

				// запись
				case 16:
					//*** BUFFER => DATA CORE

					// кол-во регистров
					quantity := get.Uint16(request[4:])

					// если кол-во регистров больше 123
					if quantity > 123 {
						set.Uint8(buffer[0:], uint8(o.cfg.Node))
						set.Uint8(buffer[1:], 0x10|0x80)
						set.Uint8(buffer[2:], modbus.ErrorIllegalDataValue)
						set.Uint16Swapped(buffer[3:], crc.Reset().Push(buffer[:3]).Result())
						o.comDrv.SetTXData(buffer[:5])
						o.comDrv.SigRTS <- struct{}{}
						continue
					}
					// смещение в ядре данных
					i := o.cfg.Core.Start + int(get.Uint16(request[2:])*2)
					// смещение в буфере
					tx := 7
					// читаем буфер
					for ; quantity > 0; quantity -= 1 {
						if i >= length || tx >= len(request) {
							break
						}
						// пишем в ядро
						set.Uint16(core.Data[i:], get.Uint16Swapped(request[tx:]))
						// один регистр это 2 байта
						tx += 2
						i += 2
					}

					// действие выполнено, заполняем ответ, вычисляем контрольную сумму, записываем и даем сигнал на отправку
					set.Uint8(buffer[0:], uint8(o.cfg.Node))
					set.Uint8(buffer[1:], 0x10)
					set.Uint16(buffer[2:], get.Uint16(request[2:]))
					set.Uint16(buffer[4:], uint16((tx-7)/2))
					set.Uint16Swapped(buffer[6:], crc.Reset().Push(buffer[:6]).Result())
					o.comDrv.SetTXData(buffer[:8])
					o.comDrv.SigRTS <- struct{}{}

				// прочее
				default:
					set.Uint8(buffer[0:], get.Uint8(request[0:]))
					set.Uint8(buffer[1:], get.Uint8(request[1:])|0x80)
					set.Uint8(buffer[2:], modbus.ErrorIllegalDataValue)
					set.Uint16Swapped(buffer[3:], crc.Reset().Push(buffer[:3]).Result())
					o.comDrv.SetTXData(buffer[:5])
					o.comDrv.SigRTS <- struct{}{}
				}
			}
		// неизвесный сигнал - предотвращаем блокировку рутины
		default:
			continue
		}
	}
}
