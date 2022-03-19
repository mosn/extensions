package main

import (
	"context"
	"encoding/json"
	"mosn.io/api"
	"mosn.io/extensions/go-plugin/pkg/common/encryption"
	"mosn.io/extensions/go-plugin/pkg/common/encryption/xor"
	"mosn.io/pkg/buffer"
	"mosn.io/pkg/log"
	"mosn.io/pkg/protocol/http"
)

// define a function named: CreateFilterFactory, do not need init to register
func CreateFilterFactory(conf map[string]interface{}) (api.StreamFilterChainFactory, error) {
	b, _ := json.Marshal(conf)
	m := make(map[string]string)
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}
	return &BumsDecoderFilterFactory{
		config: m,
	}, nil
}

// An implementation of api.StreamFilterChainFactory
type BumsDecoderFilterFactory struct {
	config map[string]string
}

func (f *BumsDecoderFilterFactory) CreateFilterChain(ctx context.Context, callbacks api.StreamFilterChainFactoryCallbacks) {
	filter := NewBumsDecoderFilter(ctx, f.config)
	// ReceiverFilter, run the filter when receive a request from downstream
	// The FilterPhase can be BeforeRoute or AfterRoute, we use BeforeRoute in this demo
	callbacks.AddStreamReceiverFilter(filter, api.BeforeRoute)
	// SenderFilter, run the filter when receive a response from upstream
	// In the demo, we are not implement this filter type
	callbacks.AddStreamSenderFilter(filter, api.BeforeSend)
}

type BumsDecoderFilter struct {
	config       map[string]string
	secretConfig *encryption.SecretConfig
	ctrlBits     string
	origSender   string
	handler      api.StreamReceiverFilterHandler
	sendHandler  api.StreamSenderFilterHandler
}

// NewBumsDecodersFilter returns a BumsDecoderFilter, the BumsDecoderFilter is an implementation of api.StreamReceiverFilter
// A Filter can implement both api.StreamReceiverFilter and api.StreamSenderFilter.
func NewBumsDecoderFilter(ctx context.Context, config map[string]string) *BumsDecoderFilter {
	secretConfig, err := encryption.ParseSecret(ctx)
	if err != nil {
		log.DefaultLogger.Errorf("[stream_filter][BumsDecoder_decoder] ParseSecret ERR: %s", err)
	}
	return &BumsDecoderFilter{
		config:       config,
		secretConfig: secretConfig,
	}
}

func (f *BumsDecoderFilter) OnReceive(ctx context.Context, headers api.HeaderMap, buf buffer.IoBuffer, trailers api.HeaderMap) api.StreamFilterStatus {
	var serviceId string
	var err error
	if _, ok := headers.(http.RequestHeader); ok {
		if buf == nil {
			return api.StreamFilterContinue
		} else {
			bodyBytes := buf.Bytes()
			//如果是密文，需要先解密
			ctrlBits := f.getCtrlBits(headers)
			origSender := f.getOrigSender(headers)
			if len(ctrlBits) != 8 {
				return api.StreamFilterContinue
			}
			bit := ctrlBits[:1]

			switch bit {
			case "1": //xor 行内加密算法异或
				bodyBytes, _ = xor.XorDecrypt(origSender, buf.Bytes(), f.secretConfig)
			case "2": //3DES
			//todo unimplemented encrypt 3DS algorithm
			case "4": //sm4
				//todo unimplemented encrypt sm4 algorithm
			}
			if bodyBytes == nil {
				return api.StreamFilterStop
			}
			//取DataId
			var body map[string]interface{}
			err = json.Unmarshal(bodyBytes, &body)
			if err != nil {
				log.DefaultLogger.Errorf("Unmarshal ERR %s", err)
				return api.StreamFilterStop
			} else {
				if body["head"] != nil {
					_v, _ := json.Marshal(body["head"])
					var bodyHead map[string]string
					json.Unmarshal(_v, &bodyHead)

					serviceId = bodyHead["tranCode"]
					if serviceId == "" {
						log.DefaultLogger.Errorf("[stream_filter][BumsDecoder_decoder] Not Found ServiceId")
						return api.StreamFilterStop
					}
				} else {
					return api.StreamFilterStop
				}
			}
			f.handler.SetRequestData(buffer.NewIoBufferBytes(bodyBytes))
		}
	}

	// inject http headerX-Target-App
	headers.Set("X-Target-App", serviceId)
	headers.Set("X-Service-Type", "bums")

	f.handler.SetRequestHeaders(headers)

	return api.StreamFilterContinue
}

func (f *BumsDecoderFilter) SetReceiveFilterHandler(handler api.StreamReceiverFilterHandler) {
	f.handler = handler
}

// SetSenderFilterHandler sets the StreamSenderFilterHandler
func (f *BumsDecoderFilter) SetSenderFilterHandler(handler api.StreamSenderFilterHandler) {
	f.sendHandler = handler
}

func (f *BumsDecoderFilter) Append(ctx context.Context, headers api.HeaderMap, buf api.IoBuffer, trailers api.HeaderMap) api.StreamFilterStatus {
	if _, ok := headers.(http.ResponseHeader); ok {
		if buf == nil {
			return api.StreamFilterContinue
		} else {
			bodyBytes := buf.Bytes()
			//如果是密文，需要先解密
			ctrlBits := f.getCtrlBits(headers)
			origSender := f.getOrigSender(headers)
			if len(ctrlBits) != 8 {
				return api.StreamFilterContinue
			}
			bit := ctrlBits[:1]
			switch bit {
			case "1": //xor 行内加密算法异或
				bodyBytes, _ = xor.XorDecrypt(origSender, buf.Bytes(), f.secretConfig)
			case "2": //3DES
			//todo unimplemented encrypt 3DS algorithm
			case "4": //sm4
				//todo unimplemented encrypt sm4 algorithm
			}
			f.handler.SetRequestData(buffer.NewIoBufferBytes(bodyBytes))
		}
	}

	return api.StreamFilterContinue
}

func (f *BumsDecoderFilter) getOrigSender(headers api.HeaderMap) string {
	origSender, ok := headers.Get("OrigSender")
	if ok {
		f.origSender = origSender
	} else {
		origSender = f.origSender
	}
	return origSender
}

func (f *BumsDecoderFilter) getCtrlBits(headers api.HeaderMap) string {
	ctrlBits, ok := headers.Get("CtrlBits")
	if ok {
		f.ctrlBits = ctrlBits
	} else {
		ctrlBits = f.ctrlBits
	}
	return ctrlBits
}

func (f *BumsDecoderFilter) OnDestroy() {}
