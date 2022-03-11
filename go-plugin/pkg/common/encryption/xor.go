package encryption

func xorBytes(key []byte, msg []byte) []byte {
	keyLen := len(key)
	msgLen := len(msg)

	for i := 0; i < msgLen; i++ {
		msg[i] ^= key[i%keyLen]
	}

	return msg
}

func XorEncoder(msg []byte, key []byte) []byte {
	return xorBytes(key, msg)
}

func XorDecoder(msg []byte, key []byte) []byte {
	return xorBytes(key, msg)
}
