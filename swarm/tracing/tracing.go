
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 12:09:50</date>
//</624342684877721600>

package tracing

import (
	"io"
	"os"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/log"
	jaeger "github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	cli "gopkg.in/urfave/cli.v1"
)

var Enabled bool = false

//
const TracingEnabledFlag = "tracing"

var (
	Closer io.Closer
)

var (
	TracingFlag = cli.BoolFlag{
		Name:  TracingEnabledFlag,
		Usage: "Enable tracing",
	}
	TracingEndpointFlag = cli.StringFlag{
		Name:  "tracing.endpoint",
		Usage: "Tracing endpoint",
		Value: "0.0.0.0:6831",
	}
	TracingSvcFlag = cli.StringFlag{
		Name:  "tracing.svc",
		Usage: "Tracing service name",
		Value: "swarm",
	}
)

//
var Flags = []cli.Flag{
	TracingFlag,
	TracingEndpointFlag,
	TracingSvcFlag,
}

//
func init() {
	for _, arg := range os.Args {
		if flag := strings.TrimLeft(arg, "-"); flag == TracingEnabledFlag {
			Enabled = true
		}
	}
}

func Setup(ctx *cli.Context) {
	if Enabled {
		log.Info("Enabling opentracing")
		var (
			endpoint = ctx.GlobalString(TracingEndpointFlag.Name)
			svc      = ctx.GlobalString(TracingSvcFlag.Name)
		)

		Closer = initTracer(endpoint, svc)
	}
}

func initTracer(endpoint, svc string) (closer io.Closer) {
//
//
	cfg := jaegercfg.Configuration{
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:            true,
			BufferFlushInterval: 1 * time.Second,
			LocalAgentHostPort:  endpoint,
		},
	}

//
//
//
	jLogger := jaegerlog.StdLogger
//

//
	closer, err := cfg.InitGlobalTracer(
		svc,
		jaegercfg.Logger(jLogger),
//
//
	)
	if err != nil {
		log.Error("Could not initialize Jaeger tracer", "err", err)
	}

	return closer
}

