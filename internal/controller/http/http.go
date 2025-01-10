package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/ogen-go/ogen/middleware"
	"github.com/pkg/errors"
	"github.com/upikoth/aireader-go/internal/config"
	"github.com/upikoth/aireader-go/internal/controller/http/handler"
	app "github.com/upikoth/aireader-go/internal/generated/app"
	"github.com/upikoth/aireader-go/internal/pkg/logger"
	"github.com/upikoth/aireader-go/internal/services"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type HTTP struct {
	logger    logger.Logger
	appServer *http.Server
}

func New(
	cfg *config.ControllerHTTP,
	loggerInstance logger.Logger,
	srvs *services.Services,
	tp trace.TracerProvider,
) (*HTTP, error) {
	handlerInstance := handler.New(
		loggerInstance,
		srvs,
		cfg,
	)

	srv, err := app.NewServer(
		handlerInstance,
		app.WithErrorHandler(getAppErrorHandler(handlerInstance)),
		app.WithTracerProvider(tp),
		app.WithMiddleware(
			httpTraceMiddleware,
			getHTTPLoggingMiddleware(loggerInstance),
		),
	)

	if err != nil {
		return nil, errors.WithStack(err)
	}

	tracedSrv := otelhttp.NewHandler(srv, "API")

	mux := http.NewServeMux()

	mux.Handle("/api/", corsMiddleware(tracedSrv))

	appInstance := http.FileServer(http.Dir("docs/app"))
	mux.Handle("/api/docs/app/", http.StripPrefix("/api/docs/app/", appInstance))

	oauthMailRuInstance := http.FileServer(http.Dir("docs/oauthmailru"))
	mux.Handle("/api/docs/oauthmailru/", http.StripPrefix("/api/docs/oauthmailru/", oauthMailRuInstance))

	oauthYandexInstance := http.FileServer(http.Dir("docs/oauthyandex"))
	mux.Handle("/api/docs/oauthyandex/", http.StripPrefix("/api/docs/oauthyandex/", oauthYandexInstance))

	swaggerUI := http.FileServer(http.Dir("docs/swagger-ui"))
	mux.Handle("/api/docs/swagger-ui/", http.StripPrefix("/api/docs/swagger-ui/", swaggerUI))

	appServer := &http.Server{
		Addr:              ":" + cfg.Port,
		ReadHeaderTimeout: time.Minute,
		Handler:           mux,
	}

	return &HTTP{
		logger:    loggerInstance,
		appServer: appServer,
	}, nil
}

func (h *HTTP) Start() error {
	return errors.WithStack(h.appServer.ListenAndServe())
}

func (h *HTTP) Stop(ctx context.Context) error {
	return errors.WithStack(h.appServer.Shutdown(ctx))
}

func getAppErrorHandler(h *handler.Handler) app.ErrorHandler {
	return func(context context.Context, w http.ResponseWriter, _ *http.Request, err error) {
		w.Header().Set("Content-Type", "application/json")

		res := h.NewError(context, err)
		w.WriteHeader(res.StatusCode)

		bytes, _ := res.Response.MarshalJSON()
		_, _ = w.Write(bytes)
	}
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Baggage, Sentry-Trace, Authorization-Token")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func httpTraceMiddleware(req middleware.Request, next middleware.Next) (middleware.Response, error) {
	span := trace.SpanFromContext(req.Context)

	bytes, _ := json.Marshal(req.Body)

	span.SetAttributes(
		attribute.String("http.request.body", string(bytes)),
	)

	for k, v := range req.Params {
		var val string

		switch vt := v.(type) {
		case app.OptInt:
			val = strconv.Itoa(vt.Value)
		default:
			val = fmt.Sprintf("%v", vt)
		}

		span.SetAttributes(
			attribute.String(
				fmt.Sprintf("http.request.params.%s", k.Name),
				val,
			),
		)
	}

	res, errorResponse := next(req)

	if errorResponse == nil {
		resBytes, _ := json.Marshal(res.Type)
		span.SetAttributes(
			attribute.String("http.response.body", string(resBytes)),
		)
	} else {
		span.SetAttributes(
			attribute.String("http.response.error", errorResponse.Error()),
		)
	}

	return res, errorResponse
}

func getHTTPLoggingMiddleware(
	logger logger.Logger,
) func(req middleware.Request, next middleware.Next) (middleware.Response, error) {
	return func(req middleware.Request, next middleware.Next) (middleware.Response, error) {
		type reqLogInfo struct {
			URL     string `json:"Url"`
			Body    string `json:"Body"`
			TraceID string `json:"TraceId"`
		}

		spanCtx := trace.SpanContextFromContext(req.Context)
		reqBody, _ := json.Marshal(req.Body)

		rli := reqLogInfo{
			URL:     req.Raw.URL.String(),
			Body:    string(reqBody),
			TraceID: spanCtx.TraceID().String(),
		}

		rliBytes, _ := json.Marshal(rli)

		logger.Debug(fmt.Sprintf("Request: %s", string(rliBytes)))

		res, errorResponse := next(req)

		type resLogInfo struct {
			URL     string `json:"Url"`
			Body    string `json:"Body"`
			Error   string `json:"Error"`
			TraceID string `json:"TraceId"`
		}

		resBody, _ := json.Marshal(res.Type)

		errStr := ""
		if errorResponse != nil {
			errStr = errorResponse.Error()
		}

		info := resLogInfo{
			URL:     req.Raw.URL.String(),
			Body:    string(resBody),
			Error:   errStr,
			TraceID: spanCtx.TraceID().String(),
		}

		infoBytes, _ := json.Marshal(info)

		logger.Debug(fmt.Sprintf("Response: %s", string(infoBytes)))

		return res, errorResponse
	}
}
