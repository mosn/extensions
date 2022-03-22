package xor

import (
	"errors"
	"mosn.io/extensions/go-plugin/pkg/common/encryption"
	"mosn.io/pkg/log"
)

func XorDecrypt(consumerId string, buf []byte, secretConfig *encryption.SecretConfig) ([]byte, error) {
	if "xor" == secretConfig.Type {
		secret := secretConfig.Secret[consumerId]
		if secret != "" {
			body := encryption.XorDecrypt(encryption.Base64Decoder(buf), []byte(secret))
			if body != nil {
				return body, nil
			}
		}
	}

	log.DefaultLogger.Errorf("[encryption][utils] xorDecrypt ERR:consumerId:%s, secretConfig: %v+", consumerId, secretConfig)
	return nil, errors.New("decrypt failed")
}

func XorEncrypt(consumerId string, buf []byte, secretConfig *encryption.SecretConfig) ([]byte, error) {
	if "xor" == secretConfig.Type {
		secret := secretConfig.Secret[consumerId]
		if secret != "" {
			body := encryption.Base64Encoder(encryption.XorEncrypt(buf, []byte(secret)))
			if body != nil {
				return body, nil
			}
		}
	}

	log.DefaultLogger.Errorf("[encryption][utils] xorEncrypt ERR:, consumerId:%s, secretConfig: %v+", consumerId, secretConfig)
	return nil, errors.New("encrypt failed")
}
