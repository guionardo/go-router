package inspect

import (
	"fmt"
	"log/slog"
	"reflect"
	"runtime"
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
	var (
		s  any = new(T)
		is     = &InspectStruct[T, R]{
			typeName: t.Name(),
		}
	)

	if sr, srt := s.(Responser[T, R]); srt {
		is.handlerFunc = is.handleSimple
		is.handlerName = runtime.FuncForPC(reflect.ValueOf(sr.Handle).Pointer()).Name()
	} else if cr, crt := s.(CustomResponser[T, R]); crt {
		is.handlerName = runtime.FuncForPC(reflect.ValueOf(cr.Handle).Pointer()).Name()
		is.handlerFunc = is.handleCustom
	} else {
		tcr := reflect.TypeFor[CustomResponser[T, R]]()
		tsr := reflect.TypeFor[Responser[T, R]]()
		return nil, t, fmt.Errorf("type %s should implements interfaces %s or %s", t.Name(), tcr.Name(), tsr.Name())
	}

	if _, ok := s.(Responser[T, R]); !ok {
		it := reflect.TypeFor[Responser[T, R]]()
		return nil, t, fmt.Errorf("struct '%s' must implement the interface %s", t.Name(), it.Name())
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
