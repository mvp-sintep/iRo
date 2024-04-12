package modbus

import (
	"iRo/internal/pkg/crc16"
	"iRo/internal/pkg/get"
)

// IsValidRTURequest - проверка на валидность modbus rtu запроса
func IsValidRTURequest(packet []byte) bool {
	// not enough data
	if len(packet) < 8 {
		return false
	}
	// more then enoght data
	if len(packet) > 255 {
		return false
	}
	// test node id
	if get.Uint8(packet[0:]) < 1 || get.Uint8(packet[0:]) > 247 {
		return false
	}
	// modbus func 1...6
	if get.Uint8(packet[1:]) > 0 && get.Uint8(packet[1:]) < 7 {
		// test packet size
		if len(packet) != 8 {
			return false
		}
	}
	// modbus func 15
	if get.Uint8(packet[1:]) == 15 {
		// test byte quantity
		if get.Uint16(packet[4:])/8 > uint16(get.Uint8(packet[6:])) {
			return false
		}
	}
	// modbus func 16
	if get.Uint8(packet[1:]) == 16 {
		// test byte quantity
		if get.Uint16(packet[4:])*2 != uint16(get.Uint8(packet[6:])) {
			return false
		}
	}
	// modbus func 15 or 16
	if get.Uint8(packet[1:]) == 15 || get.Uint8(packet[1:]) == 16 {
		// test packet size
		if len(packet) != int(get.Uint8(packet[6:]))+9 {
			return false
		}
	}
	// test checksum
	return (&crc16.Value{}).Reset().Push(packet[:len(packet)-2]).Result() == get.Uint16Swapped(packet[len(packet)-2:])
}
