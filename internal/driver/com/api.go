package com

// GetRXCount - определяем кол-во принятых байт
func (o *Driver) GetRXCount() int {
	return o.rxCount
}

// GetRXData - забираем принятые данные
func (o *Driver) GetRXData() (int, []byte) {
	tmp := make([]byte, o.rxCount)
	copy(tmp, o.rxBuffer[:o.rxCount])
	return o.rxCount, tmp
}

// SetRXData - записываем отправляемые данные
func (o *Driver) SetTXData(data []byte) {
	if tx := len(data); tx > 4 && tx < len(o.txBuffer) {
		copy(o.txBuffer[:], data)
		o.txCount = tx
	}
}

// Reset - очищаем данные
func (o *Driver) Reset() {
	o.rxCount = 0
	o.txCount = 0
}
