package util

// CRC16ModbusCheckCode generate modbus crc16 check code
func CRC16ModbusCheckCode(data []byte) (hi, lo byte, checkCode uint16) {
	checkCode = 0xffff
	l := len(data)
	for i := 0; i < l; i++ {
		checkCode ^= uint16(data[i])
		for j := 0; j < 8; j++ {
			if checkCode&0x0001 > 0 {
				checkCode = (checkCode >> 1) ^ 0xA001
			} else {
				checkCode >>= 1
			}
		}
	}
	return byte(checkCode & 0xFF), byte(checkCode >> 8), checkCode
}

// CRC16ModbusCheck modbus crc 16 check
func CRC16ModbusCheck(data []byte) bool {
	n := len(data)
	if n < 2 {
		return false
	}
	hi, low, _ := CRC16ModbusCheckCode(data[:(n - 2)])
	if data[(n-2)] != hi && data[(n-1)] != low {
		return false
	}
	return true
}
