package transport

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	myErrors "writer/internal/errors"
	myLog "writer/internal/logger"
	"writer/internal/models"

	"github.com/valyala/fasthttp"
)

func ParseJsonOrder(ctx *fasthttp.RequestCtx) (models.Order, error) {
	var order models.Order
	err := json.NewDecoder(bytes.NewReader(ctx.Request.Body())).Decode(&order)
	if err != nil {
		//myLog.Log.Errorf("error in parse json", err.Error())
		fmt.Println(err.Error())
		return models.Order{}, myErrors.ErrParseJSON
	}

	// _, err = time.Parse("2006/01/02", order.DateCreated)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// 	return order, nil
	// }

	return order, nil
}

func ParseGenerateJsonOrder(orderjson string) (models.Order, error) {
	var order models.Order
	err := json.NewDecoder(bytes.NewReader([]byte(orderjson))).Decode(&order)
	if err != nil {
		myLog.Log.Errorf("error in parse json", err.Error())
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
		return err //myErrors.ErrWriteJSON
	}
	return nil
}
