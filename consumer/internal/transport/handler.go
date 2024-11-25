package transport

import (
	"fmt"

	"consumer/internal/config"
	my_errors "consumer/internal/errors"
	myLog "consumer/internal/logger"
	"consumer/internal/service"

	"net/http"

	"github.com/fasthttp/router"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/valyala/fasthttp"
)

type HandlersBuilder struct {
	srv  *service.Srv
	rout *router.Router
}

func HandleCreate(cfg config.Config, s *service.Srv) {
	hb := HandlersBuilder{
		srv:  s,
		rout: router.New(),
	}

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(":8090", nil)
	}()

	hb.rout.GET("/WB/get", hb.Get())
	fmt.Println(fasthttp.ListenAndServe(":8080", hb.rout.Handler))
}

func (hb *HandlersBuilder) Get() func(ctx *fasthttp.RequestCtx) {
	myLog.Log.Infof("Start func Get")
	return metrics(func(ctx *fasthttp.RequestCtx) {
		if ctx.IsGet() {
			orderUUID, err_ := ParseJsonUUID(ctx)
			if err_ != nil {
				err_ := WriteJson(ctx, err_.Error())

				if err_ != nil {
				}

			} else {
				myLog.Log.Debugf("sucsess parse")
				order, err := hb.srv.GetOrderSrv(orderUUID)
				if err != nil {
					WriteJson(ctx, err.Error())
				} else {
					myLog.Log.Debugf("sucsess get")
					WriteJsonOrder(ctx, order)
				}
			}
		} else {
			WriteJson(ctx, my_errors.ErrMethodNotAllowed.Error())
			ctx.SetStatusCode(fasthttp.StatusMethodNotAllowed)
			myLog.Log.Debugf("MethodNotAllowed")
		}
	}, "Get")
}