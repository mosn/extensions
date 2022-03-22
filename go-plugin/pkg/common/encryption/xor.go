package encryption

func XorEncrypt(msg []byte, key []byte) []int8 {
	keyLen := len(key)
	msgLen := len(msg)

	msgResult := make([]int8, len(msg))

	for i := 0; i < msgLen; i++ {
		msgResult[i] = int8(msg[i]) ^ int8(key[i%keyLen])
	}

	return msgResult
}

func XorDecrypt(msg []int8, key []byte) []byte {
	keyLen := len(key)
	msgLen := len(msg)

	msgResult := make([]byte, len(msg))

	for i := 0; i < msgLen; i++ {
		msgResult[i] = byte(msg[i]) ^ key[i%keyLen]
	}

	return msgResult
}
