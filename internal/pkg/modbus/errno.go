package modbus

const (
	ErrorIllegalFunction     = 0x01
	ErrorIllegalDataAddress  = 0x02
	ErrorIllegalDataValue    = 0x03
	ErrorSlaveDeviceFailure  = 0x04
	ErrorAcknowledge         = 0x05
	ErrorSlaveDeviceBusy     = 0x06
	ErrorNegativeAcknowledge = 0x07
)
