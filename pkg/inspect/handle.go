package inspect

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/guionardo/go-router/pkg/config"
	"github.com/guionardo/go-router/pkg/logging"
)

func (s *InspectStruct[T, R]) Handle(w http.ResponseWriter, r *http.Request) {
	payload, err := s.Parse(r)
	if err != nil {
		answerError(w, err)
		return
	}
	s.handlerFunc(w, r, payload)
}

func answerError(w http.ResponseWriter, err error) {
	pe := NewParseError(err)
	w.Header().Add("Content-Type", "application/json")
	status := http.StatusBadRequest
	if responseErr, ok := err.(responseError); ok && status <= 0 {
		status = responseErr.StatusCode()
	}

	if status <= 0 {
		status = http.StatusBadGateway
	}
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(pe)
}

func (s *InspectStruct[T, R]) handleCustom(w http.ResponseWriter, r *http.Request, payload *T) error {
	startTime := time.Now()
	if err := s.customResponser.Handle(w, r, payload); err != nil {
		var statusCode any
		if hErr, ok := err.(responseError); ok {
			statusCode = hErr.StatusCode()
		} else {
			statusCode = "unknown"
		}
		logging.Get().Warn("", slog.String("method", r.Method),
			slog.String("handler", s.reqType.String()),
			slog.Int64("ms", time.Since(startTime).Milliseconds()),
			slog.Any("status", statusCode),
			slog.String("path", r.URL.String()))

		return err
	}
	return nil
}

func (s *InspectStruct[T, R]) handleSimple(w http.ResponseWriter, r *http.Request, payload *T) error {
	if config.DevelopmentMode {
		w.Header().Add("X-Router-Request-Type", s.reqType.String())
		w.Header().Add("X-Router-Response-Type", s.respType.String())
	}
	// Simple handler
	var (
		err       error
		status    int
		response  *R
		startTime = time.Now()
	)

	response, status, err = s.responser.Handle(r, payload)
	if responseErr, ok := err.(responseError); ok && status <= 0 {
		status = responseErr.StatusCode()
	}

	if status <= 0 {
		status = http.StatusBadGateway
	}

	if err != nil || response != nil {
		w.Header().Add("Content-Type", "application/json")

	}

	if err != nil {
		w.Header().Add("X-Handler-Error", err.Error())
	}
	w.WriteHeader(status)

	if response != nil {
		json.NewEncoder(w).Encode(response)
	} else if err != nil {
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
	}

	logging.Get().Debug("REQUEST",
		slog.String("method", r.Method),
		slog.String("handler", s.reqType.String()),
		slog.Int64("ms", time.Since(startTime).Milliseconds()),
		slog.Int("status", status),
		slog.String("path", r.URL.String()))
	return err
}
