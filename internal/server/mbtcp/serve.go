package mbtcp

import (
	"iRo/internal/pkg/get"
	"iRo/internal/pkg/modbus"
	"iRo/internal/pkg/set"
)

// serve - обработка запросов MODBUS TCP
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
	buffer := make([]byte, 259)

	for {
		select {
		// контекст закрывается
		case <-o.context.Done():
			return
			// сигнал от драйвера - соединение установлено и пакет принят
		case conn := <-o.tcpDrv.SigDSR:
			// забираем данные
			if n, request := conn.GetRXData(); n > 0 {
				//log.Printf("[%d]->[%x]", n, request[:n])
				// drop packet if no data
				if n < 12 || get.Uint16(request[2:]) != 0 || get.Uint16(request[4:]) != uint16(n-6) {
					continue
				}

				// селектор по номеру modbus функции
				switch get.Uint8(request[7:]) {

				// чтение регистров хранения
				case 3:

					//*** DATA CORE => BUFFER

					// кол-во регистров
					quantity := get.Uint16(request[10:])
					// если больше 125
					if quantity > 125 {
						// отсекаем по пределу (2+2+2+1+1+1+125*2=259)
						quantity = 125
					}
					// смещение в ядре данных
					i := o.cfg.Core.Start + int(get.Uint16(request[9:])*2)
					// смещение в выходном буфере
					tx := 9
					// заполняем буфер
					for ; quantity > 0; quantity -= 1 {
						// заполнитель значений за пределами ядра данных
						tmp := uint16(0)
						if i < length {
							tmp = get.Uint16Swapped((*o.core)[i:])
						}
						// пишем в буфер
						set.Uint16(buffer[tx:], tmp)
						// один регистр = 2 байта
						tx += 2
						i += 2
					}

					set.Uint16(buffer[0:], get.Uint16(request[0:]))
					set.Uint16(buffer[2:], 0)
					set.Uint16(buffer[4:], uint16(tx-6))
					set.Uint8(buffer[6:], get.Uint8(request[6:]))

					// если есть данные по запросу
					if tx > 9 {
						set.Uint8(buffer[7:], 0x03)
						set.Uint8(buffer[8:], uint8(tx-9))
					} else {
						// если нет данных
						set.Uint8(buffer[7:], 0x03|0x80)
						set.Uint8(buffer[8:], modbus.ErrorIllegalDataAddress)
					}

					// данные готовы, записываем и даем сигнал на отправку
					conn.SetTXData(buffer[:tx])
					o.tcpDrv.SigRTS <- conn

				// запись
				case 16:

					//*** BUFFER => DATA CORE

					// кол-во регистров
					quantity := get.Uint16(request[10:])

					// если кол-во регистров больше 122 или неправильно указано кол-во байт
					if quantity > 123 || get.Uint8(request[12:]) != byte(quantity*2) {
						set.Uint16(buffer[0:], get.Uint16(request[0:]))
						set.Uint16(buffer[2:], 0)
						set.Uint16(buffer[4:], 3)
						set.Uint8(buffer[6:], get.Uint8(request[6:]))
						set.Uint8(buffer[7:], 16|0x80)
						set.Uint8(buffer[8:], modbus.ErrorIllegalDataAddress)
						conn.SetTXData(buffer[:9])
						o.tcpDrv.SigRTS <- conn
						continue
					}

					// смещение в ядре данных
					i := o.cfg.Core.Start + int(get.Uint16(request[8:])*2)
					// смещение в буфере
					tx := 13
					// читаем буфер
					for ; quantity > 0; quantity -= 1 {
						if i >= length || tx >= len(request) {
							break
						}
						// пишем в ядро
						set.Uint16((*o.core)[i:], get.Uint16Swapped(request[tx:]))
						// one register 2 bytes
						tx += 2
						i += 2
					}

					set.Uint16(buffer[0:], get.Uint16(request[0:]))
					set.Uint16(buffer[2:], 0)
					set.Uint16(buffer[4:], 6)
					set.Uint8(buffer[6:], get.Uint8(request[6:]))
					set.Uint8(buffer[7:], 16)
					set.Uint16(buffer[8:], get.Uint16(request[8:]))
					set.Uint16(buffer[10:], uint16((tx-13)/2))
					conn.SetTXData(buffer[:12])
					o.tcpDrv.SigRTS <- conn

					// данные обновлены
					o.nda <- struct{}{}

				default:
					set.Uint16(buffer[0:], get.Uint16(request[0:]))
					set.Uint16(buffer[2:], 0)
					set.Uint16(buffer[4:], 3)
					set.Uint8(buffer[6:], get.Uint8(request[6:]))
					set.Uint8(buffer[7:], get.Uint8(request[7:])|0x80)
					set.Uint8(buffer[8:], modbus.ErrorIllegalDataValue)
					conn.SetTXData(buffer[:9])
					o.tcpDrv.SigRTS <- conn
				}
			}
		}
	}

}
