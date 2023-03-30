package config

import (
	"context"
	"encoding/json"
	"errors"
	"sync"

	"mosn.io/pkg/log"
	"mosn.io/pkg/variable"
)

const ExtendConfigKey = "global_extend_config"

var globalExtendConfig = NewGlobalExtendConfig()

type ConfigHandler func(val string)

type ExtendConfig struct {
	effective   bool
	once        sync.Once
	dirtyConfig sync.Map
	config      sync.Map
	handler     chan string
	sendh       map[string]ConfigHandler
}

func NewGlobalExtendConfig() *ExtendConfig {
	return &ExtendConfig{
		handler: make(chan string),
		sendh:   make(map[string]ConfigHandler),
	}
}

func (gex *ExtendConfig) Register(ctx context.Context, key string, handler ConfigHandler) {
	gex.sendh[key] = handler
	gex.initConfig(ctx)
}

func (gex *ExtendConfig) initConfig(ctx context.Context) {
	var gerr error
	gex.once.Do(func() {
		recvl, err := variable.Get(ctx, ExtendConfigKey)
		if err != nil {
			gerr = errors.New("the dynamic config chan is not exist")
			return
		}
		rec, ok := recvl.(chan chan string)
		if !ok {
			gerr = errors.New("the dynamic config chan is not exist")
			return
		}
		// 第一次获取数据，阻塞获取
		rec <- gex.handler
		gex.effective = true
		val := <-gex.handler
		gex.parse(val)
		go gex.handlerConfig()
	})
	if gerr != nil {
		log.DefaultLogger.Errorf("init config failed,err:%s", gerr)
	}
}

func (gex *ExtendConfig) handlerConfig() {
	for val := range gex.handler {
		gex.parse(val)
	}
}

func (gex *ExtendConfig) parse(value string) error {
	cc := make(map[string]string)
	if err := json.Unmarshal([]byte(value), &cc); err != nil {
		return err
	}
	// 更新&添加事件
	for key, value := range cc {
		handler, ok := gex.sendh[key]
		if ok {
			handler(value)
		}
		gex.dirtyConfig.Delete(key)
	}
	// 删除事件
	gex.dirtyConfig.Range(func(key, value interface{}) bool {
		handler, ok := gex.sendh[key.(string)]
		if ok {
			handler("")
		}
		return true
	})

	var dirtyConfig, config sync.Map
	for key, value := range cc {
		dirtyConfig.Store(key, value)
		config.Store(key, value)
	}
	gex.config = config
	gex.dirtyConfig = dirtyConfig
	return nil
}

func (gex *ExtendConfig) SyncMapByConfig(ctx context.Context) (sync.Map, bool) {
	gex.initConfig(ctx)
	return gex.config, gex.effective
}

func (gex *ExtendConfig) GetConfig(ctx context.Context, key string) (string, bool) {
	gex.initConfig(ctx)
	val, ok := gex.config.Load(key)
	return val.(string), ok
}
