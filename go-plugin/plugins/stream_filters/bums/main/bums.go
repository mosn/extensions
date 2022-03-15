package main

import (
	"context"
	"encoding/json"
	"fmt"
	"mosn.io/pkg/protocol/http"
	"sync/atomic"

	"github.com/mosn/extensions/go-plugin/pkg/common/encryption"
	"mosn.io/api"
	"mosn.io/pkg/buffer"
	"mosn.io/pkg/log"
)

// define a function named: CreateFilterFactory, do not need init to register
func CreateFilterFactory(conf map[string]interface{}) (api.StreamFilterChainFactory, error) {
	b, _ := json.Marshal(conf)
	m := make(map[string]string)
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}
	return &BumsFilterFactory{
		config: m,
	}, nil
}

// An implementation of api.StreamFilterChainFactory
type BumsFilterFactory struct {
	config map[string]string
}

func (f *BumsFilterFactory) CreateFilterChain(ctx context.Context, callbacks api.StreamFilterChainFactoryCallbacks) {
	filter := NewBumsFilter(ctx, f.config)
	// ReceiverFilter, run the filter when receive a request from downstream
	// The FilterPhase can be BeforeRoute or AfterRoute, we use BeforeRoute in this demo
	callbacks.AddStreamReceiverFilter(filter, api.BeforeRoute)
	// SenderFilter, run the filter when receive a response from upstream
	// In the demo, we are not implement this filter type
	// callbacks.AddStreamSenderFilter(filter, api.BeforeSend)
}

type BumsFilter struct {
	config       map[string]string
	secretConfig *encryption.SecretConfig
	handler      api.StreamReceiverFilterHandler
}

// NewBumssFilter returns a BumsFilter, the BumsFilter is an implementation of api.StreamReceiverFilter
// A Filter can implement both api.StreamReceiverFilter and api.StreamSenderFilter.
func NewBumsFilter(ctx context.Context, config map[string]string) *BumsFilter {
	value := ctx.Value("codec_config").(*atomic.Value)
	secretConfig, err := encryption.ParseSecret(value)
	if err != nil {
		log.DefaultLogger.Errorf("[stream_filter][bums] ParseSecret ERR: %s", err)
	}
	return &BumsFilter{
		config:       config,
		secretConfig: secretConfig,
	}
}

func (f *BumsFilter) OnReceive(ctx context.Context, headers api.HeaderMap, buf buffer.IoBuffer, trailers api.HeaderMap) api.StreamFilterStatus {
	passed := true
	var serviceId string
	var err error
	if _, ok := headers.(http.RequestHeader); ok {
		if buf == nil {
			passed = false
		} else {
			//如果是密文，需要现解密
			ctrlBits := getCtrlBits(headers)
			origSender := getOrigSender(headers)
			bit := ctrlBits[0]
			bodyBytes := buf.Bytes()
			switch bit {
			case 1: //xor 行内加密算法异或
				bodyBytes, err = f.decrypt("xor", origSender, buf)
			case 2: //3DES
			//todo
			case 4: //sm4
				//todo
			}
			if bodyBytes == nil {
				passed = false
			}
			var body map[string]interface{}
			err = json.Unmarshal(bodyBytes, &body)
			if err != nil {
				passed = false
				log.DefaultLogger.Errorf("Unmarshal ERR %s", err)
			} else {
				if body["head"] != nil {
					_v, _ := json.Marshal(body["head"])
					var bodyHead map[string]string
					json.Unmarshal(_v, &bodyHead)

					serviceId = bodyHead["tranCode"]
					if serviceId == "" {
						passed = false
						log.DefaultLogger.Errorf("[stream_filter][bums] Not Found ServiceId")
					}
				} else {
					passed = false
				}
			}

		}
	}

	if !passed {
		return api.StreamFilterStop
	}
	// inject http header
	headers.Set("X-Target-App", serviceId)
	headers.Set("X-Service-Type", "springcloud")

	return api.StreamFilterContinue
}

func (f *BumsFilter) SetReceiveFilterHandler(handler api.StreamReceiverFilterHandler) {
	f.handler = handler
}

func (f *BumsFilter) decrypt(secType string, consumerId string, buf buffer.IoBuffer) ([]byte, error) {
	if secType == f.secretConfig.Type {
		secret := f.secretConfig.Secret[consumerId]
		if secret != "" {
			body := encryption.XorDecrypt(buf.Bytes(), []byte(secret))
			if body != nil {
				return body, nil
			}
		}
	}

	log.DefaultLogger.Errorf("[stream_filter][bums] decrypt ERR:secType:%s, consumerId:%s, secretConfig: %v+", secType, consumerId, secretConfig)
	return nil, fmt.Errorf("decrypt failed")
}

func getOrigSender(headers api.HeaderMap) string {
	origSender, _ := headers.Get("OrigSender")
	return origSender
}

func getCtrlBits(headers api.HeaderMap) string {
	ctrlBits, _ := headers.Get("CtrlBits")
	return ctrlBits
}

func (f *BumsFilter) OnDestroy() {}
