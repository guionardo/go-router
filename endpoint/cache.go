package endpoint

import (
	"fmt"
	"log/slog"
	"reflect"
	"sync"

	"github.com/guionardo/go-router/pkg/logging"
	reflections "github.com/guionardo/go-router/pkg/reflect"
)

type structCache struct {
	lock    sync.RWMutex
	structs map[string]any
}

var pool = &structCache{
	structs: map[string]any{},
}

func getEndpoint[T, R any]() (*Endpoint[T, R], reflect.Type, error) {
	t := reflect.TypeFor[T]()
	pool.lock.RLock()
	defer pool.lock.RUnlock()
	if is, ok := pool.structs[t.Name()]; ok {
		return is.(*Endpoint[T, R]), t, nil
	}
	if !reflections.IsStruct[T]() {
		logging.Get().Warn("inspectStructget: expected a struct", slog.String("type", t.Name()))
		return nil, t, fmt.Errorf("expected a struct to make an InspectStruct. Got %s", t.Name())
	}
	var is = &Endpoint[T, R]{
		reqType: t,
	}

	if err := is.buildResponser(); err != nil {
		return nil, t, err
	}

	return is, t, nil
}

func setEndpoint[T, R any](t reflect.Type, is *Endpoint[T, R]) {
	pool.lock.Lock()
	defer pool.lock.Unlock()
	if is == nil {
		delete(pool.structs, t.Name())
	} else {
		pool.structs[t.Name()] = is
	}
}
