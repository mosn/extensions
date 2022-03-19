package encryption

import (
	"context"
	"encoding/json"
	"errors"
	"mosn.io/pkg/log"
	"sync/atomic"
)

type SecretConfig struct {
	Enable bool              `json:"enable"`
	Type   string            `json:"type"`
	Secret map[string]string `json:"secrets"`
}

func ParseSecret(ctx context.Context) (*SecretConfig, error) {
	value := ctx.Value("code_config")
	if value == nil {
		return nil, errors.New("code_config is empty")
	}
	atomicValue := value.(*atomic.Value)
	valueStr := atomicValue.Load()
	if valueStr == nil {
		return nil, errors.New("code_config is empty")
	}
	secretConfigValue := valueStr.(string)

	if log.DefaultLogger.GetLogLevel() >= log.DEBUG {
		log.DefaultLogger.Debugf("[common] [encryption]: SecretConfig parse: %s", secretConfigValue)
	}

	if secretConfigValue == "" {
		return nil, errors.New("secretConfigValue is empty")
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
		return nil, errors.New("secretConfig is not enabled")
	}

	return secretConfig, nil
}
