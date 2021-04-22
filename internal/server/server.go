package server

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/peakle/recycle/internal"
	"github.com/peakle/recycle/internal/handlers"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/pprofhandler"
)

func StartServer(ctx context.Context) error {
	var m = internal.InitManager()
	defer m.Close()

	var h = handlers.Handler{
		Manager: m,
	}

	var requestHandler = func(ctx *fasthttp.RequestCtx) {
		path := strings.ToLower(string(ctx.Path()))

		if strings.HasPrefix(path, "/v1/order/subscribe") && string(ctx.Request.Header.Method()) == fasthttp.MethodPost {
			h.Subscribe(ctx)
		} else if strings.HasPrefix(path, "/v1/order/create") && string(ctx.Request.Header.Method()) == fasthttp.MethodPost {
			h.Create(ctx)
		} else if strings.HasPrefix(path, "/v1/order/list") && string(ctx.Request.Header.Method()) == fasthttp.MethodGet {
			h.List(ctx)
		} else if strings.HasPrefix(path, "/v1/order/info") && string(ctx.Request.Header.Method()) == fasthttp.MethodGet {
			h.Info(ctx)
		} else if strings.HasPrefix(path, "/debug/pprof") {
			pprofhandler.PprofHandler(ctx)
		} else {
			ctx.SetConnectionClose()
			return
		}
	}

	var server = fasthttp.Server{
		Handler:         requestHandler,
		IdleTimeout:     1 * time.Minute,
		TCPKeepalive:    true,
		CloseOnShutdown: true,
	}

	go func() {
		<-ctx.Done()

		var err = server.Shutdown()
		if err != nil {
			log.Printf("on shutdown api server: %s", err)
		}
	}()

	return server.ListenAndServe(":80")
}
