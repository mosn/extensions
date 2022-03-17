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
	return &BumsEncoderFilterFactory{
		config: m,
	}, nil
}

// An implementation of api.StreamFilterChainFactory
type BumsEncoderFilterFactory struct {
	config map[string]string
}

func (f *BumsEncoderFilterFactory) CreateFilterChain(ctx context.Context, callbacks api.StreamFilterChainFactoryCallbacks) {
	filter := NewBumsEncoderFilter(ctx, f.config)
	// ReceiverFilter, run the filter when receive a request from downstream
	// The FilterPhase can be BeforeRoute or AfterRoute, we use BeforeRoute in this demo
	callbacks.AddStreamReceiverFilter(filter, api.AfterRoute)
	// SenderFilter, run the filter when receive a response from upstream
	// In the demo, we are not implement this filter type
	// callbacks.AddStreamSenderFilter(filter, api.BeforeSend)
}

type BumsEncoderFilter struct {
	config       map[string]string
	secretConfig *encryption.SecretConfig
	handler      api.StreamReceiverFilterHandler
}

// NewBumsEncodersFilter returns a BumsEncoderFilter, the BumsEncoderFilter is an implementation of api.StreamReceiverFilter
// A Filter can implement both api.StreamReceiverFilter and api.StreamSenderFilter.
func NewBumsEncoderFilter(ctx context.Context, config map[string]string) *BumsEncoderFilter {
	//value := ctx.Value("codec_config").(*atomic.Value)
	value := &atomic.Value{}
	value.Store("{\"enable\":true, \"type\":\"xor\", \"secrets\":{\"ESB002\":\"12345678\",\"ESB001\":\"13213211\"}}")
	secretConfig, err := encryption.ParseSecret(value)
	if err != nil {
		log.DefaultLogger.Errorf("[stream_filter][BumsEncoder_decoder] ParseSecret ERR: %s", err)
	}
	return &BumsEncoderFilter{
		config:       config,
		secretConfig: secretConfig,
	}
}

func (f *BumsEncoderFilter) OnReceive(ctx context.Context, headers api.HeaderMap, buf buffer.IoBuffer, trailers api.HeaderMap) api.StreamFilterStatus {
	passed := true
	bodyBytes := buf.Bytes()
	if _, ok := headers.(http.RequestHeader); ok {
		if buf == nil {
			passed = false
		} else {
			//如果是密文，需要现解密
			ctrlBits := getCtrlBits(headers)
			origSender := getOrigSender(headers)
			bit := ctrlBits[:1]
			switch bit {
			case "1": //xor 行内加密算法异或
				bodyBytes, _ = f.xorEncrypt(origSender, buf)
			case "2": //3DES
			//todo
			case "4": //sm4
				//todo
			}
			if bodyBytes == nil {
				passed = false
			}
		}
	}

	if !passed {
		return api.StreamFilterStop
	}
	// inject http header
	headers.Del("X-Target-App")
	headers.Del("X-Service-Type")

	f.handler.SetRequestData(buffer.NewIoBufferBytes(bodyBytes))
	f.handler.SetRequestHeaders(headers)

	return api.StreamFilterContinue
}

func (f *BumsEncoderFilter) SetReceiveFilterHandler(handler api.StreamReceiverFilterHandler) {
	f.handler = handler
}

func (f *BumsEncoderFilter) Append(ctx context.Context, headers api.HeaderMap, buf api.IoBuffer, trailers api.HeaderMap) api.StreamFilterStatus {

	bodyBytes := buf.Bytes()
	if _, ok := headers.(http.ResponseHeader); ok {
		if buf == nil {
			return api.StreamFilterContinue
		} else {
			//如果是密文，需要先解密
			ctrlBits := getCtrlBits(headers)
			origSender := getOrigSender(headers)
			bit := ctrlBits[:1]

			switch bit {
			case "1": //xor 行内加密算法异或
				bodyBytes, _ = f.xorEncrypt(origSender, buf)
			case "2": //3DES
			//todo
			case "4": //sm4
				//todo
			}
			if bodyBytes == nil {
				return api.StreamFilterContinue
			}
		}
	}

	f.handler.SetRequestData(buffer.NewIoBufferBytes(bodyBytes))

	return api.StreamFilterContinue
}

func (f *BumsEncoderFilter) xorEncrypt(consumerId string, buf buffer.IoBuffer) ([]byte, error) {
	if "xor" == f.secretConfig.Type {
		secret := f.secretConfig.Secret[consumerId]
		if secret != "" {
			body := encryption.Base64Encoder(encryption.XorEncrypt(buf.Bytes(), []byte(secret)))
			if body != nil {
				return body, nil
			}
		}
	}

	log.DefaultLogger.Errorf("[stream_filter][BumsEncoder_decoder] xorEncrypt ERR:, consumerId:%s, secretConfig: %v+", consumerId, f.secretConfig)
	return nil, fmt.Errorf("encrypt failed")
}

func getOrigSender(headers api.HeaderMap) string {
	origSender, _ := headers.Get("OrigSender")
	return origSender
}

func getCtrlBits(headers api.HeaderMap) string {
	ctrlBits, _ := headers.Get("CtrlBits")
	return ctrlBits
}

func (f *BumsEncoderFilter) OnDestroy() {}
