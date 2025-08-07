package inspect

import (
	"fmt"
	"log/slog"
	"reflect"
	"sync"

	"github.com/guionardo/go-router/pkg/logging"
)

type structCache struct {
	lock    sync.RWMutex
	structs map[string]any
}

var (
	pool = &structCache{
		structs: map[string]any{},
	}
)

func inspectStructGet[T, R any]() (*InspectStruct[T, R], reflect.Type, error) {
	t := reflect.TypeFor[T]()
	pool.lock.RLock()
	defer pool.lock.RUnlock()
	if is, ok := pool.structs[t.Name()]; ok {
		return is.(*InspectStruct[T, R]), t, nil
	}
	if !IsStruct[T]() {
		logging.Get().Warn("inspectStructget: expected a struct", slog.String("type", t.Name()))
		return nil, t, fmt.Errorf("expected a struct to make an InspectStruct. Got %s", t.Name())
	}
	var is = &InspectStruct[T, R]{
		reqType: t,
	}

	if err := is.buildResponser(); err != nil {
		return nil, t, err
	}

	return is, t, nil
}

func poolSet[T, R any](t reflect.Type, is *InspectStruct[T, R]) {
	pool.lock.Lock()
	defer pool.lock.Unlock()
	if is == nil {
		delete(pool.structs, t.Name())
	} else {
		pool.structs[t.Name()] = is
	}
}
