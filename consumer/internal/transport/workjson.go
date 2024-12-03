package transport

import (
	"bytes"
	"encoding/json"
	"time"

	myErrors "consumer/internal/errors"
	"consumer/internal/models"

	"github.com/valyala/fasthttp"
)

func ParseJsonOrder(ctx *fasthttp.RequestCtx) (models.Order, error) {
	var order models.Order
	err := json.NewDecoder(bytes.NewReader(ctx.Request.Body())).Decode(&order)
	if err != nil {
		return models.Order{}, myErrors.ErrParseJSON
	}

	_, err = time.Parse("2006/01/02", order.DateCreated)
	if err != nil {
		return order, nil
	}

	return models.Order{}, myErrors.ErrParseJSON
}
func WriteJson(ctx *fasthttp.RequestCtx, s string) error {
	ctx.SetContentType("application/json")
	ctx.Response.BodyWriter()
	err := json.NewEncoder((*ctx).Response.BodyWriter()).Encode(s)
	if err != nil {
		return err
	}
	return nil
}

func WriteJsonOrder(ctx *fasthttp.RequestCtx, order models.Order) error {
	ctx.SetContentType("application/json")
	ctx.Response.BodyWriter()
	err := json.NewEncoder((*ctx).Response.BodyWriter()).Encode(order)
	if err != nil {
		return err
	}
	return nil
}
