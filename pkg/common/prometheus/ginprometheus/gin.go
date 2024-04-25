package ginprometheus

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	gin "github.com/gin-gonic/gin"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Metric struct {
	MetricCollector prometheus.Collector
	ID              string
	Name            string
	Description     string
	Type            string
	Args            []string
}

type RequestCounterURLLabelMappingFunc func(c *gin.Context) string
type Prometheus struct {
	PrometheusPushGateway
	reqCnt                    *prometheus.CounterVec
	reqDur                    *prometheus.HistogramVec
	reqSz, resSz              *prometheus.Summary
	router                    *gin.Engine
	listenAddress             string
	MetricsList               []*Metric
	MetricsPath               string
	ReqCntURILabelMappingFunc RequestCounterURLLabelMappingFunc
	URILabelFromContext       string
}

type PrometheusPushGateway struct {
	PushIntervalSeconds time.Duration
	PushGatewayURL      string
	MetricsURL          string
	Job                 string
}

func NewPrometheus(system string, customMetricsList ...[]*Metric) *Prometheus {
	p := &Prometheus{}

	return p
}

func (p *Prometheus) Use(e *gin.Engine) error {
	e.Use(p.HandlerFunc())
	return p.SetMetricsPath(e)
}

func (p *Prometheus) SetMetricsPath(e *gin.Engine) error {
	if p.listenAddress != "" {
		p.router.GET(p.MetricsPath, prometheusHandler())
		return p.start()
	}
	e.GET(p.MetricsPath, prometheusHandler())
	return nil
}

func (p *Prometheus) start() error {
	return p.router.Run(p.listenAddress)
}

func prometheusHandler() gin.HandlerFunc {
	h := promhttp.Handler()
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func (p *Prometheus) HandlerFunc() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.Path == p.MetricsPath {
			c.Next()
			return
		}
		c.Next()
	}
}

func (p *Prometheus) SetListenAddress(address string) {
	p.listenAddress = address
	if p.listenAddress != "" {
		p.router = gin.Default()
	}
}

func Exit() {
	progName := filepath.Base(os.Args[0])
	fmt.Fprintf(os.Stderr, "Warning %s receive process terminal SIGTERM exit 0\n", progName)
}
