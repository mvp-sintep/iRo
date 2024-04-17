package tcp

// GetRXCount - определить кол-во принятых байт
func (o *Connection) GetRXCount() int {
	return o.rxCount
}

// GetRXData - забрать принятые данные
func (o *Connection) GetRXData() (int, []byte) {
	tmp := make([]byte, o.rxCount)
	copy(tmp, o.rxBuffer[:o.rxCount])
	return o.rxCount, tmp
}

// SetRXData - записать отправляемые данные
func (o *Connection) SetTXData(data []byte) {
	if tx := len(data); tx > 4 && tx < len(o.txBuffer) {
		copy(o.txBuffer[:], data)
		o.txCount = tx
	}
}

// Reset - очистить данные
func (o *Connection) Reset() {
	o.rxCount = 0
	o.txCount = 0
}
