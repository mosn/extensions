package common

import "strconv"

func BuildLenOfPacket(mesLen int, size int) string {
	lenOfPacket := strconv.Itoa(mesLen)
	if len(lenOfPacket) < size {
		remain := size - len(lenOfPacket)
		prefix := ""
		for i := 0; i < remain; i++ {
			prefix += "0"
		}
		lenOfPacket = prefix + lenOfPacket
	}
	return lenOfPacket
}
