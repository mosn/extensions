package encryption

import (
	"encoding/json"
	"fmt"
	"mosn.io/pkg/log"
	"sync/atomic"
)

type SecretConfig struct {
	Enable bool              `json:"enable"`
	Type   string            `json:"type"`
	Secret map[string]string `json:"secrets"`
}

func ParseSecret(value *atomic.Value) (*SecretConfig, error) {

	secretConfigValue := value.Load().(string)

	if log.DefaultLogger.GetLogLevel() >= log.DEBUG {
		log.DefaultLogger.Debugf("[common] [encryption]: SecretConfig parse: %s", secretConfigValue)
	}

	if secretConfigValue == "" {
		return nil, fmt.Errorf("secretConfigValue is empty")
	}

	secretConfig := &SecretConfig{}
	err := json.Unmarshal([]byte(secretConfigValue), secretConfig)

	if err != nil {
		if log.DefaultLogger.GetLogLevel() >= log.DEBUG {
			log.DefaultLogger.Debugf("[common] [encryption]: SecretConfig unmarshal error, json:%s, error:%v", secretConfigValue, err)
		}
		return nil, err
	}

	if !secretConfig.Enable {
		return nil, fmt.Errorf("secretConfig is not enabled")
	}

	return secretConfig, nil
}
