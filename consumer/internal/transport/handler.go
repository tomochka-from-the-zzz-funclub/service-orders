package transport

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"

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
	srv   service.InterfaceService
	rout  *router.Router
	templ *template.Template
}

func HandleCreate(cfg config.Config, s service.InterfaceService) {
	fmt.Println(os.Getwd())
	fmt.Println(os.ReadDir("./src"))
	t, err := os.Getwd()
	absolutePath, err := filepath.Abs(t)
	if err != nil {
		fmt.Println("Ошибка при получении абсолютного пути:", err)
		return
	}

	fmt.Println("Абсолютный путь к директории:", absolutePath)
	tmpl, err := template.ParseFiles("../app/fw.html")
	if err != nil {
		myLog.Log.Fatalf("GetHtml error during parsing of file: %v", err)
		return
	}

	hb := HandlersBuilder{
		srv:   s,
		rout:  router.New(),
		templ: tmpl,
	}

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(":8090", nil)
	}()

	hb.rout.GET("/WB/get", hb.Get())
	hb.rout.GET("/WB/get", hb.GetHtml())
	fmt.Println(fasthttp.ListenAndServe(":8080", hb.rout.Handler))
}

// func (hb *HandlersBuilder) Get() func(ctx *fasthttp.RequestCtx) {
// 	myLog.Log.Infof("start func Get")
// 	return metrics(func(ctx *fasthttp.RequestCtx) {
// 		if ctx.IsGet() {
// 			orderUUIDjson := string(ctx.QueryArgs().Peek("order_uid"))
// 			myLog.Log.Debugf("sucsess parse json in func Get with id %+v", orderUUIDjson)
// 			order, err := hb.srv.GetOrderSrv(orderUUIDjson)
// 			if err != nil {
// 				// err_ := WriteJson(ctx, err.Error())
// 				// if err_ != nil {
// 				myLog.Log.Warnf("there is no way to record an error")
// 				//}
// 				ctx.SetStatusCode(fasthttp.StatusNotFound)
// 			} else {
// 				myLog.Log.Debugf("sucsess get")
// 				err_ := WriteJsonOrder(ctx, order)
// 				if err_ != nil {
// 					myLog.Log.Warnf("there is no way to record an error")
// 				}
// 			}
// 		} else {
// 			err_ := WriteJson(ctx, my_errors.ErrMethodNotAllowed.Error())
// 			if err_ != nil {
// 				myLog.Log.Warnf("there is no way to record an error")
// 			}
// 			myLog.Log.Warnf("MethodNotAllowed")
// 		}
// 	}, "Get")
// }

func (hb *HandlersBuilder) GetHtml() func(ctx *fasthttp.RequestCtx) {
	return metrics(func(ctx *fasthttp.RequestCtx) {
		myLog.Log.Debugf("Start func GetHtml")
		if ctx.IsGet() {
			err := hb.templ.Execute(ctx.Response.BodyWriter(), nil)
			if err != nil {
				myLog.Log.Errorf("GetHtml error during executing of file: %v", err)
				ctx.Response.SetStatusCode(400)
				return
			}
			//ctx.Response.Header.Set("Content-Type", "text/plain")

			ctx.Response.Header.Add("content-type", "text/html")
		}
	}, "GetHtml")
}

func (hb *HandlersBuilder) Get() func(ctx *fasthttp.RequestCtx) {
	myLog.Log.Infof("Start func Get")
	return metrics(func(ctx *fasthttp.RequestCtx) {
		if ctx.IsGet() {
			orderUUID := string(ctx.QueryArgs().Peek("order_uid"))
			myLog.Log.Debugf("sucsess parse json in func Get with id %+v", orderUUID)

			//orderUUID, err_ := ParseJsonUUID(ctx)
			// if err_ != nil {
			// 	err_ := WriteJson(ctx, err_.Error())
			// 	if err_ != nil {
			// 	}
			// } else {
			//myLog.Log.Debugf("sucsess parse")
			order, err := hb.srv.GetOrderSrv(orderUUID)
			if err != nil {
				WriteJson(ctx, err.Error())
			} else {
				myLog.Log.Debugf("sucsess get")
				WriteJsonOrder(ctx, order)
			}
			//}
		} else {
			WriteJson(ctx, my_errors.ErrMethodNotAllowed.Error())
			ctx.SetStatusCode(fasthttp.StatusMethodNotAllowed)
			myLog.Log.Debugf("MethodNotAllowed")
		}
	}, "Get")
}
