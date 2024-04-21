package natsworker

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/leonardinius/go-service-template/internal/apiserv"
)

type Worker interface {
	ListenAndServe(ctx context.Context) error
	Shutdown(ctx context.Context) error
}

type worker struct {
	server  *http.Server
	handler http.Handler
	natsCon *nats.Conn
	routes  []apiserv.Route
}

var _ Worker = (*worker)(nil)

func NewWorker(ctx context.Context, serverURL string, routes []apiserv.Route, options ...Option) (Worker, error) {
	config, natsioOptions, err := buildNatsIOOptions(ctx, options...)
	if err != nil {
		return nil, err
	}

	server, err := apiserv.NewDefaultServer(ctx, config.metricsAddress, routes...)
	if err != nil {
		return nil, err
	}

	// save original originalHandler
	originalHandler := server.Handler
	server.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != apiserv.MetricsRoutePath {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte("Use NATS.io for rpc invocation. Metrics are available at " + apiserv.MetricsRoutePath + "."))
			return
		}
		originalHandler.ServeHTTP(w, r)
	})

	natsCon, err := nats.Connect(serverURL, natsioOptions...)
	if err != nil {
		if shutdownErr := server.Shutdown(context.WithoutCancel(ctx)); shutdownErr != nil {
			err = errors.Join(err, fmt.Errorf("failed to shutdown metrics server: %w", shutdownErr))
		}
		return nil, err
	}

	return &worker{
		server,
		originalHandler,
		natsCon,
		routes,
	}, nil
}

// ListenAndServe implements Worker.
func (w *worker) ListenAndServe(ctx context.Context) error {
	var err error
	// Subscribe to NATS.io messages.
	for _, route := range w.routes {
		subscribePath := urlToSuscribeSubject(route.Path())
		msgHandler := w.natsMsgHandler(ctx, subscribePath, route.Path())
		_, routeErr := w.natsCon.Subscribe(subscribePath, msgHandler)
		slog.LogAttrs(ctx, slog.LevelInfo, "nats subscribe error", slog.String("subject", subscribePath))
		err = errors.Join(err, routeErr)
	}
	if err != nil {
		return err
	}

	return apiserv.ListenAndServe(ctx, w.server)
}

// Shutdown implements Worker.
func (w *worker) Shutdown(ctx context.Context) error {
	var err error
	if w.natsCon != nil && !w.natsCon.IsClosed() {
		err = errors.Join(err, w.natsCon.Drain())
		w.natsCon.Close()
	}
	if w.server != nil {
		err = errors.Join(err, w.server.Shutdown(context.WithoutCancel(ctx)))
	}
	return err
}

func (w *worker) natsMsgHandler(ctx context.Context, subscribePath, handlerPath string) nats.MsgHandler {
	return func(msg *nats.Msg) {
		req, err := NewRequestFromMessage(msg, subscribePath, handlerPath)
		if err != nil {
			gwErrror := errToSharedError(err)
			data, errRespond := protojson.Marshal(gwErrror)
			errRespond = errors.Join(errRespond, msg.Respond(data))
			slog.LogAttrs(ctx, slog.LevelError,
				"NATS.io to gPRC request error",
				slog.String("subject", msg.Subject),
				slog.String("error", errRespond.Error()),
				slog.String("parent_error", err.Error()),
			)
		}
		buffer := bytes.NewBufferString("")
		resp := NewStdResponseWriter(buffer)

		w.handler.ServeHTTP(resp, req)
		resp.Header().Set("X-Status-Code", strconv.Itoa(resp.statusCode))

		err = msg.RespondMsg(&nats.Msg{
			Subject: msg.Subject,
			Reply:   msg.Reply,
			Data:    buffer.Bytes(),
			Header:  nats.Header(resp.Header()),
		})
		if err != nil {
			slog.LogAttrs(ctx, slog.LevelError,
				"NATS.io to gPRC response error",
				slog.String("subject", msg.Subject),
				slog.String("error", err.Error()),
			)
		}
	}
}
