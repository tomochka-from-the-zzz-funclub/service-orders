package transport

import (
	"fmt"
	"os"
	"time"

	my_errors "github.com/tomochka-from-the-zzz-funclub/go-L0-Kafka/internal/err"
	myLog "github.com/tomochka-from-the-zzz-funclub/go-L0-Kafka/internal/logger"
	publisher "github.com/tomochka-from-the-zzz-funclub/go-L0-Kafka/internal/publisher"

	"net/http"

	"github.com/fasthttp/router"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/valyala/fasthttp"

	"github.com/rs/zerolog"
	//"github.com/rs/zerolog/log"
)

type HandlersBuilder struct {
	pub  *publisher.KafkaClient
	lg   zerolog.Logger
	rout *router.Router
}

func HandleCreate() {
	hb := HandlersBuilder{
		pub:  publisher.NewKafkaClient(),
		lg:   zerolog.New(os.Stderr).With().Timestamp().Logger().Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.UnixDate}),
		rout: router.New(),
	}
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(":8090", nil)
	}()

	hb.rout.POST("/WB/set", hb.Set())
	fmt.Println(fasthttp.ListenAndServe(":8081", hb.rout.Handler))
}

func (hb *HandlersBuilder) Set() func(ctx *fasthttp.RequestCtx) {
	myLog.Log.Debugf("Start func Set")
	return metrics(func(ctx *fasthttp.RequestCtx) {
		if ctx.IsPost() {
			order, err := ParseJsonOrder(ctx)
			if err != nil {
				WriteJson(ctx, "Error parse")
			}
			myLog.Log.Debugf("SEND")
			//fmt.Println(hb.pub.Producer)
			if hb.pub == nil {
				fmt.Println("клиент не инициализирован")
			}
			err = hb.pub.SendOrderToKafka("records", order)

			//order, err := ParseJsonOrder(ctx)
			if err != nil {
				myLog.Log.Errorf("SendOrderToKafka", err.Error())
				WriteJson(ctx, "Error SendMessageOrder")
			}
			hb.lg.Info().Msg("SendOrderToKafka")
		} else {
			WriteJson(ctx, my_errors.ErrMethodNotAllowed.Error())
			ctx.SetStatusCode(fasthttp.StatusMethodNotAllowed)
			hb.lg.Warn().
				Msgf("message from func Set %v", my_errors.ErrMethodNotAllowed.Error())
		}
	}, "Set")
}
