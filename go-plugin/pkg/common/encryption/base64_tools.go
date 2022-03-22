package encryption

var (
	alphabet = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/=")
	codes    = make([]int16, 256)
)

func init() {
	for i := 0; i < 256; i++ {
		if i >= 48 && i <= 57 {
			codes[i] = (int16)(52 + i - 48)
		} else if i >= 65 && i <= 90 {
			codes[i] = (int16)(i - 65)
		} else if i >= 97 && i <= 122 {
			codes[i] = (int16)(26 + i - 97)
		} else {
			codes[i] = -1
		}
	}

	codes[43] = 62
	codes[47] = 63
}

func Base64Encoder(srcByte []int8) []byte {
	i := 0

	destByte := make([]byte, (len(srcByte)+2)/3*4)
	for index := 0; i < len(srcByte); index += 4 {
		quad := false
		trip := false

		val := 255 & int(srcByte[i])
		val <<= 8
		if i+1 < len(srcByte) {
			val |= 255 & int(srcByte[i+1])
			trip = true
		}

		val <<= 8
		if i+2 < len(srcByte) {
			val |= 255 & int(srcByte[i+2])
			quad = true
		}
		if quad {
			destByte[index+3] = alphabet[val&63]
		} else {
			destByte[index+3] = alphabet[64]
		}
		val >>= 6
		if trip {
			destByte[index+2] = alphabet[val&63]
		} else {
			destByte[index+2] = alphabet[64]
		}
		val >>= 6
		destByte[index+1] = alphabet[val&63]
		val >>= 6
		destByte[index+0] = alphabet[val&63]
		i += 3
	}

	return destByte
}

func Base64Decoder(base64Bytes []byte) []int8 {
	srcByte := base64Bytes
	length := (len(srcByte) + 3) / 4 * 3
	if len(srcByte) > 0 && srcByte[len(srcByte)-1] == 61 {
		length--
	}

	if len(srcByte) > 1 && srcByte[len(srcByte)-2] == 61 {
		length--
	}

	destByte := make([]int8, length)
	var shift int16 = 0
	var accum int16 = 0
	var index = 0

	for ix := 0; ix < len(srcByte); ix++ {
		value := codes[srcByte[ix]&255]
		if value >= 0 {
			accum <<= 6
			shift += 6
			accum |= value
			if shift >= 8 {
				shift -= 8
				destByte[index] = (int8)(accum >> shift & 255)
				index++
			}
		}
	}

	if index != len(destByte) {
		return nil
	} else {
		return destByte
	}
}
