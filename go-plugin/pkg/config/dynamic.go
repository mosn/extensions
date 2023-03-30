package config

import (
	"context"
	"sync"
)

func Register(ctx context.Context, key string, handler ConfigHandler) {
	globalExtendConfig.Register(ctx, key, handler)
}

func GlobalExtendMapByContext(ctx context.Context) (*sync.Map, bool) {
	if cfg, ok := globalExtendConfig.SyncMapByConfig(ctx); ok {
		return cfg, ok
	}
	cfg, ok := ctx.Value(ExtendConfigKey).(*sync.Map)
	return cfg, ok
}

func GlobalExtendConfigByContext(ctx context.Context, key string) (string, bool) {
	cfg, ok := GlobalExtendMapByContext(ctx)
	if !ok {
		return "", false
	}
	info, ok := cfg.Load(key)
	if !ok {
		return "", false
	}
	sinfo, ok := info.(string)
	return sinfo, ok
}
