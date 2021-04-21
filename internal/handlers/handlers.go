package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/peakle/recycle/internal"
	"github.com/peakle/recycle/internal/storages"
	"github.com/valyala/fasthttp"
	"log"
)

type Handler struct {
	Manager *internal.SQLManager
	_       log.Logger
	_       Config
}

func (h *Handler) Subscribe(ctx *fasthttp.RequestCtx) {
	defer ctx.Response.Header.Set("Content-Type", "application/json")

	var success, err = storages.Subscribe(ctx, h.Manager, string(ctx.QueryArgs().Peek("order_id")))
	if err != nil {
		log.Printf("on Subscribe: %s", err)
		ctx.SetStatusCode(fasthttp.StatusInternalServerError) /// TODO delete
		return
	}

	_, _ = fmt.Fprintf(ctx, "{\"success\": \"%v\"}", success)
}

func (h *Handler) Create(ctx *fasthttp.RequestCtx) {
	defer ctx.Response.Header.Set("Content-Type", "application/json")

	var userId = string(ctx.QueryArgs().Peek("user_id")) // TODO validation
	var address = string(ctx.QueryArgs().Peek("address"))
	var maxSize = string(ctx.QueryArgs().Peek("max_size"))
	var eventAt = string(ctx.QueryArgs().Peek("event_at"))

	var orderId, err = storages.Create(ctx, h.Manager, userId, address, maxSize, eventAt)
	if err != nil {
		log.Printf("on Create: %s", err)
		ctx.SetStatusCode(fasthttp.StatusInternalServerError) /// TODO delete
		return
	}

	_, _ = fmt.Fprintf(ctx, "{\"order_id\": \"%v\"}", orderId)
}

func (h *Handler) List(ctx *fasthttp.RequestCtx) {
	defer ctx.Response.Header.Set("Content-Type", "application/json")

	var orders, err = storages.GetList(ctx, h.Manager)
	if err != nil {
		log.Printf("on List: %s", err)
		ctx.SetStatusCode(fasthttp.StatusInternalServerError) /// TODO delete
		return
	}

	resp, err := json.Marshal(orders)
	if err != nil {
		log.Printf("on List: on Marshal: %s", err)
		ctx.SetStatusCode(fasthttp.StatusInternalServerError) // TODO delete
		return
	}

	_, _ = fmt.Fprint(ctx, string(resp))
}

func (h *Handler) Info(ctx *fasthttp.RequestCtx) {
	defer ctx.Response.Header.Set("Content-Type", "application/json")

	var orders, err = storages.Info(ctx, h.Manager, string(ctx.QueryArgs().Peek("order_id")))
	if err != nil {
		log.Printf("on Info: %s", err)
		ctx.SetStatusCode(fasthttp.StatusInternalServerError) /// TODO delete
		return
	}

	resp, err := json.Marshal(orders)
	if err != nil {
		log.Printf("on Info: on Marshal: %s", err)
		ctx.SetStatusCode(fasthttp.StatusInternalServerError) // TODO delete
		return
	}

	_, _ = fmt.Fprint(ctx, string(resp))
}
