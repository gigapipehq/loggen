package prom

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"runtime"
	"time"

	"github.com/denisbrodbeck/machineid"
	"github.com/golang/snappy"
	"github.com/matishsiao/goInfo"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	io_prometheus_client "github.com/prometheus/client_model/go"
	promconfig "github.com/prometheus/common/config"
	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/prompb"
	"github.com/prometheus/prometheus/storage/remote"

	"github.com/gigapipehq/loggen/internal/config"
)

type lib struct {
	linesCount  prometheus.Counter
	bytesCount  prometheus.Counter
	errorsCount prometheus.Counter
}

var l = lib{}

func Initialize(ctx context.Context, cfg *config.Config) chan struct{} {
	mid, err := machineid.ID()
	if err != nil {
		mid = "00000000-0000-0000-0000-000000000000"
	}

	info, _ := goInfo.GetInfo()
	labels := prometheus.Labels{
		"device_id":         mid,
		"host_arch":         runtime.GOARCH,
		"host_cpus":         fmt.Sprintf("%d", info.CPUs),
		"config_url":        cfg.URL,
		"config_api_key":    cfg.APIKey,
		"config_api_secret": cfg.APISecret,
		"config_rate":       fmt.Sprintf("%d", cfg.Rate),
		"config_timeout":    cfg.Timeout.String(),
	}
	for k, v := range cfg.Labels {
		labels[k] = v
	}
	l = lib{
		linesCount: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace:   "loggen",
			Name:        "lines_sent_total",
			Help:        "total number of lines sent",
			ConstLabels: labels,
		}),
		bytesCount: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace:   "loggen",
			Name:        "bytes_sent_total",
			Help:        "total number of bytes sent",
			ConstLabels: labels,
		}),
		errorsCount: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace:   "loggen",
			Name:        "errors_total",
			Help:        "total number of errors received",
			ConstLabels: labels,
		}),
	}

	reg := prometheus.NewRegistry()
	reg.MustRegister(
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(
			collectors.ProcessCollectorOpts{
				Namespace: "loggen",
			},
		),
		l.linesCount,
		l.bytesCount,
		l.errorsCount,
	)

	qch := make(chan struct{})
	go func() {
		u, _ := url.Parse(fmt.Sprintf("%s/api/v1/prom/remote/write", cfg.URL))
		client, _ := remote.NewWriteClient("loggen", &remote.ClientConfig{
			URL:     &promconfig.URL{URL: u},
			Timeout: model.Duration(10 * time.Second),
			HTTPClientConfig: promconfig.HTTPClientConfig{
				BasicAuth: &promconfig.BasicAuth{
					Username: cfg.APIKey,
					Password: promconfig.Secret(cfg.APISecret),
				},
			},
		})
		t := time.NewTicker(time.Second * 5)
		for {
			select {
			case <-ctx.Done():
				sendMetrics(reg, client)
				t.Stop()
				qch <- struct{}{}
				return
			case <-t.C:
				sendMetrics(reg, client)
			}
		}
	}()
	return qch
}

func AddLines(count int) {
	addToCount(l.linesCount, count)
}

func AddBytes(count int) {
	addToCount(l.bytesCount, count)
}

func AddErrors(count int) {
	addToCount(l.errorsCount, count)
}

func addToCount(counter prometheus.Counter, count int) {
	if counter != nil {
		counter.Add(float64(count))
	}
}

func sendMetrics(reg *prometheus.Registry, client remote.WriteClient) {
	metrics, err := reg.Gather()
	if err != nil {
		log.Printf("unable to gather metrics: %v", err)
	}

	timeseries := []prompb.TimeSeries{}
	metadata := []prompb.MetricMetadata{}
	for _, family := range metrics {
		var ftype prompb.MetricMetadata_MetricType
		switch family.GetType() {
		case io_prometheus_client.MetricType_COUNTER:
			ftype = prompb.MetricMetadata_COUNTER
		case io_prometheus_client.MetricType_GAUGE:
			ftype = prompb.MetricMetadata_GAUGE
		default:
			continue
		}
		metadata = append(metadata, prompb.MetricMetadata{
			Type:             ftype,
			MetricFamilyName: family.GetName(),
			Help:             family.GetHelp(),
		})

		samples := []prompb.Sample{}
		plabels := []prompb.Label{
			{
				Name:  "__name__",
				Value: family.GetName(),
			},
		}
		for _, s := range family.GetMetric() {
			var value float64
			if s.GetCounter() != nil {
				value = s.GetCounter().GetValue()
			} else if s.GetGauge() != nil {
				value = s.GetGauge().GetValue()
			}

			ts := s.GetTimestampMs()
			if ts == 0 {
				ts = time.Now().UnixMilli()
			}
			samples = append(samples, prompb.Sample{
				Value:     value,
				Timestamp: ts,
			})
			for _, v := range s.GetLabel() {
				plabels = append(plabels, prompb.Label{
					Name:  v.GetName(),
					Value: v.GetValue(),
				})
			}
		}

		timeseries = append(timeseries, prompb.TimeSeries{
			Labels:  plabels,
			Samples: samples,
		})
	}

	b, _ := (&prompb.WriteRequest{
		Timeseries: timeseries,
		Metadata:   metadata,
	}).Marshal()
	encoded := snappy.Encode(nil, b)
	if err := client.Store(context.Background(), encoded); err != nil {
		log.Printf("unable to store metrics: %v", err)
	}
}
