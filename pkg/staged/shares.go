package staged

import (
	"context"
	"sync"
)

func NewShared(ctx context.Context) context.Context {
	sm := ctx.Value(SharedMemoryScope)
	if smInstance, ok := sm.(*SharedMemory); ok {
		return context.WithValue(ctx, SharedMemoryScope, smInstance.Children())
	}

	return context.WithValue(ctx, SharedMemoryScope, NewSharedMemory())
}

func NewSharedMemory() *SharedMemory {
	return &SharedMemory{
		values: &sync.Map{},
		father: nil,
	}
}

type SharedMemory struct {
	father *SharedMemory
	values *sync.Map
}

func (sm *SharedMemory) Set(key, value any) {
	sm.values.Store(key, value)
}

func (sm *SharedMemory) Children() *SharedMemory {
	return &SharedMemory{
		values: &sync.Map{},
		father: sm,
	}
}

func (sm *SharedMemory) Get(key any, defaultValue any) (any, bool) {
	res, has := sm.values.Load(key)
	if !has {
		if sm.father == nil {
			return defaultValue, false
		} else {
			return sm.father.Get(key, defaultValue)
		}
	}

	return res, true
}

func Get[T any](sm *SharedMemory, key any, defaultValue T) (T, bool) {
	if sm == nil {
		return defaultValue, false
	}
	result, has := sm.Get(key, defaultValue)
	if has {
		if res, ok := result.(T); ok {
			return res, true
		}
	}

	return defaultValue, false
}
func SetFromContext(ctx context.Context, key, value any) context.Context {
	sm := ctx.Value(SharedMemoryScope)
	if _, ok := sm.(*SharedMemory); !ok {
		ctx = NewShared(ctx)
	}
	sm = ctx.Value(SharedMemoryScope)
	sm.(*SharedMemory).Set(key, value)
	return ctx
}

func GetFromContext[T any](ctx context.Context, key any, defaultValue T) (T, bool) {
	sm := ctx.Value(SharedMemoryScope)
	if smInstance, ok := sm.(*SharedMemory); ok {
		return Get(smInstance, key, defaultValue)
	}
	return defaultValue, false
}
