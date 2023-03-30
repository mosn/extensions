package config

import (
	"context"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"mosn.io/pkg/variable"
	"testing"
	"time"
)

func init() {
	variable.Register(variable.NewVariable(ExtendConfigKey, nil, nil, variable.DefaultSetter, 0))
}

func MockContext(value string) context.Context {
	receiver := make(chan chan string, 1)
	ctx := variable.NewVariableContext(context.Background())
	variable.Set(ctx, ExtendConfigKey, receiver)
	go func() {
		rec := <-receiver
		rec <- value
	}()
	return ctx
}

func TestExtendConfig_Register(t *testing.T) {
	gex := NewGlobalExtendConfig()
	val := `{"key":"v1"}`
	rval := "v1"
	cc := make(map[string]string)
	if err := json.Unmarshal([]byte(val), &cc); err != nil {
	}
	// 更新&添加事件
	handler := func(value string) {
		assert.Equal(t, value, rval)
	}
	// add
	gex.Register(MockContext(val), "key", handler)
	time.Sleep(time.Millisecond * 5)
	// update
	val = `{"key":"value"}`
	rval = "value"
	gex.handler <- val
	time.Sleep(time.Millisecond * 5)
	// delete
	val = `{}`
	rval = ""
	gex.handler <- val
	time.Sleep(time.Millisecond * 5)
}
