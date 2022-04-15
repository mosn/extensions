package config

import (
	"context"
	"sync"
)

const ExtendConfigKey = "global_extend_config"

func GlobalExtendMapByContext(ctx context.Context) (*sync.Map, bool) {
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
