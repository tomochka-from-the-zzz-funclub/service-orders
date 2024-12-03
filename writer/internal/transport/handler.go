package transport

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"strconv"

	"net/http"
	myErrors "writer/internal/errors"
	"writer/internal/generator"
	myLog "writer/internal/logger"
	"writer/internal/models"
	publisher "writer/internal/publisher"

	"github.com/fasthttp/router"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/valyala/fasthttp"
)

type HandlersBuilder struct {
	pub  publisher.InterfaceKafkaClient
	tmpl *template.Template
	rout *router.Router
}

func HandleCreate() {
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
		pub:  publisher.NewKafkaClient(),
		rout: router.New(),
		tmpl: tmpl,
	}
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(":8090", nil)
	}()

	hb.rout.POST("/api/v1/set", hb.Set())
	hb.rout.GET("/api/v1", hb.GetHtml())
	hb.rout.POST("/api/v1/set/generate", hb.SetGenerateJson())
	myLog.Log.Debugf("Service writer started")
	fmt.Println(fasthttp.ListenAndServe(":8081", hb.rout.Handler))
}

func (hb *HandlersBuilder) GetHtml() func(ctx *fasthttp.RequestCtx) {
	return metrics(func(ctx *fasthttp.RequestCtx) {
		myLog.Log.Debugf("Start func GetHtml")
		if ctx.IsGet() {
			err := hb.tmpl.Execute(ctx.Response.BodyWriter(), nil)
			if err != nil {
				myLog.Log.Errorf("GetHtml error during executing of file: %v", err)
				ctx.Response.SetStatusCode(400)
				return
			}

			ctx.Response.Header.Add("content-type", "text/html")
		}
	}, "GetHtml")
}

func (hb *HandlersBuilder) Set() func(ctx *fasthttp.RequestCtx) {
	return metrics(func(ctx *fasthttp.RequestCtx) {
		myLog.Log.Debugf("Start func Set: %+v", string(ctx.Request.Body()))
		if ctx.IsPost() {
			fmt.Println(string(ctx.Request.Header.Method()))
			order, err := ParseJsonOrder(ctx)
			if err != nil {
				myLog.Log.Errorf("err: %v", err.Error())
				WriteJson(ctx, "Error parse")
			}
			myLog.Log.Debugf("SEND")
			if hb.pub == nil {
				fmt.Println("клиент не инициализирован")
			}
			if models.CheckValidOrder(order) {
				err = hb.pub.SendOrderToKafka("records", order)

				if err != nil {
					myLog.Log.Errorf("SendOrderToKafka", err.Error())
				}
				myLog.Log.Debugf("SendOrderToKafka")
			} else {
				ctx.SetStatusCode(fasthttp.StatusBadRequest)
			}

		} else {
			ctx.SetStatusCode(fasthttp.StatusMethodNotAllowed)
			myLog.Log.Warnf("message from func Set %v", myErrors.ErrMethodNotAllowed.Error())
		}
	}, "Set")
}

func (hb *HandlersBuilder) SetGenerateJson() func(ctx *fasthttp.RequestCtx) {
	return metrics(func(ctx *fasthttp.RequestCtx) {

		myLog.Log.Debugf("Start func SetGenerateJson, Method: %+v", ctx.Request.Header.Method())
		if ctx.IsPost() {
			quantity_s := string(ctx.QueryArgs().Peek("quantity"))
			fmt.Println(quantity_s)
			if quantity_s == "" {
				fmt.Println("def 5")
				quantity_s = "5"
			} else {
				quantity, err := strconv.Atoi(quantity_s)
				if err != nil {
					ctx.SetStatusCode(fasthttp.StatusBadRequest)
				} else {
					for i := 0; i < quantity; i++ {
						orderjson, err := generator.GenerateOrder()
						order, err := ParseGenerateJsonOrder(string(orderjson))
						if err != nil {
							WriteJson(ctx, "Error parse")
						}
						myLog.Log.Debugf("SEND")
						if hb.pub == nil {
							fmt.Println("клиент не инициализирован")
						}
						if models.CheckValidOrder(order) {
							err = hb.pub.SendOrderToKafka("records", order)

							if err != nil {
								myLog.Log.Errorf("SendOrderToKafka", err.Error())
								ctx.SetStatusCode(fasthttp.StatusInternalServerError)
							}
							myLog.Log.Debugf("SendOrderToKafka")
						} else {
							ctx.SetStatusCode(fasthttp.StatusBadRequest)
						}
					}
				}

			}
		} else {
			ctx.SetStatusCode(fasthttp.StatusMethodNotAllowed)
			myLog.Log.Warnf("message from func Set %v", myErrors.ErrMethodNotAllowed.Error())
		}
	}, "SetGenerateJson")
}
